package taskrun

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/common"
	"github.com/bytebase/bytebase/backend/common/log"
	"github.com/bytebase/bytebase/backend/component/config"
	"github.com/bytebase/bytebase/backend/component/dbfactory"
	"github.com/bytebase/bytebase/backend/component/state"
	api "github.com/bytebase/bytebase/backend/legacyapi"
	"github.com/bytebase/bytebase/backend/plugin/db"
	"github.com/bytebase/bytebase/backend/store"
	"github.com/bytebase/bytebase/backend/store/model"
	"github.com/bytebase/bytebase/backend/utils"
	storepb "github.com/bytebase/bytebase/proto/generated-go/store"
)

// Executor is the task executor.
type Executor interface {
	// RunOnce will be called periodically by the scheduler until terminated is true.
	//
	// NOTE
	//
	// 1. It's possible that err could be non-nil while terminated is false, which
	// usually indicates a transient error and will make scheduler retry later.
	// 2. If err is non-nil, then the detail field will be ignored since info is provided in the err.
	// driverCtx is used by the database driver so that we can cancel the query
	// while have the ability to cleanup migration history etc.
	RunOnce(ctx context.Context, driverCtx context.Context, task *store.TaskMessage, taskRunUID int) (terminated bool, result *storepb.TaskRunResult, err error)
}

// RunExecutorOnce wraps a TaskExecutor.RunOnce call with panic recovery.
func RunExecutorOnce(ctx context.Context, driverCtx context.Context, exec Executor, task *store.TaskMessage, taskRunUID int) (terminated bool, result *storepb.TaskRunResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			panicErr, ok := r.(error)
			if !ok {
				panicErr = errors.Errorf("%v", r)
			}
			slog.Error("TaskExecutor PANIC RECOVER", log.BBError(panicErr), log.BBStack("panic-stack"))
			terminated = true
			result = nil
			err = errors.Errorf("TaskExecutor PANIC RECOVER, err: %v", panicErr)
		}
	}()

	return exec.RunOnce(ctx, driverCtx, task, taskRunUID)
}

func getMigrationInfo(ctx context.Context, stores *store.Store, profile *config.Profile, task *store.TaskMessage, migrationType db.MigrationType, statement string, schemaVersion model.Version, sheetID *int) (*db.MigrationInfo, error) {
	if schemaVersion.Version == "" {
		return nil, errors.Errorf("empty schema version")
	}
	instance, err := stores.GetInstanceV2(ctx, &store.FindInstanceMessage{UID: &task.InstanceID})
	if err != nil {
		return nil, err
	}
	database, err := stores.GetDatabaseV2(ctx, &store.FindDatabaseMessage{UID: task.DatabaseID})
	if err != nil {
		return nil, err
	}
	if database == nil {
		return nil, errors.Errorf("database not found")
	}
	environment, err := stores.GetEnvironmentV2(ctx, &store.FindEnvironmentMessage{ResourceID: &database.EffectiveEnvironmentID})
	if err != nil {
		return nil, err
	}

	mi := &db.MigrationInfo{
		InstanceID:     &instance.UID,
		DatabaseID:     &database.UID,
		CreatorID:      task.CreatorID,
		ReleaseVersion: profile.Version,
		Type:           migrationType,
		Version:        schemaVersion,
		Description:    task.Name,
		Environment:    environment.ResourceID,
		Database:       database.DatabaseName,
		Namespace:      database.DatabaseName,
		Payload:        &storepb.InstanceChangeHistoryPayload{},
	}

	plans, err := stores.ListPlans(ctx, &store.FindPlanMessage{PipelineID: &task.PipelineID})
	if err != nil {
		return nil, err
	}
	if len(plans) == 1 {
		planTypes := []store.PlanCheckRunType{store.PlanCheckDatabaseStatementSummaryReport}
		status := []store.PlanCheckRunStatus{store.PlanCheckRunStatusDone}
		runs, err := stores.ListPlanCheckRuns(ctx, &store.FindPlanCheckRunMessage{
			PlanUID: &plans[0].UID,
			Type:    &planTypes,
			Status:  &status,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to list plan check runs")
		}
		sort.Slice(runs, func(i, j int) bool {
			return runs[i].UID > runs[j].UID
		})
		foundChangedResources := false
		for _, run := range runs {
			if foundChangedResources {
				break
			}
			if run.Config.InstanceUid != int32(task.InstanceID) {
				continue
			}
			if run.Config.DatabaseName != database.DatabaseName {
				continue
			}
			if sheetID != nil && run.Config.SheetUid != int32(*sheetID) {
				continue
			}
			if run.Result == nil {
				continue
			}
			for _, result := range run.Result.Results {
				if result.Status != storepb.PlanCheckRunResult_Result_SUCCESS {
					continue
				}
				if report := result.GetSqlSummaryReport(); report != nil {
					mi.Payload.ChangedResources = report.ChangedResources
					foundChangedResources = true
					break
				}
			}
		}
	}

	issue, err := stores.GetIssueV2(ctx, &store.FindIssueMessage{PipelineID: &task.PipelineID})
	if err != nil {
		slog.Error("failed to find containing issue", log.BBError(err))
	}
	if issue != nil {
		// Concat issue title and task name as the migration description so that user can see
		// more context of the migration.
		mi.Description = fmt.Sprintf("%s - %s", issue.Title, task.Name)
		mi.ProjectUID = &issue.Project.UID
		mi.IssueUID = &issue.UID
	}

	mi.Source = db.UI
	creator, err := stores.GetUserByID(ctx, task.CreatorID)
	if err != nil {
		// If somehow we unable to find the principal, we just emit the error since it's not
		// critical enough to fail the entire operation.
		slog.Error("Failed to fetch creator for composing the migration info",
			slog.Int("task_id", task.ID),
			log.BBError(err),
		)
	} else {
		mi.Creator = creator.Name
		mi.CreatorID = creator.ID
	}

	statement = strings.TrimSpace(statement)
	// Only baseline and SDL migration can have empty sql statement, which indicates empty database.
	if mi.Type != db.Baseline && mi.Type != db.MigrateSDL && statement == "" {
		return nil, errors.Errorf("empty statement")
	}
	return mi, nil
}

func getCreateTaskRunLog(ctx context.Context, taskRunUID int, s *store.Store, profile *config.Profile) func(t time.Time, e *storepb.TaskRunLog) error {
	return func(t time.Time, e *storepb.TaskRunLog) error {
		return s.CreateTaskRunLog(ctx, taskRunUID, t.UTC(), profile.DeployID, e)
	}
}

func doMigration(
	ctx context.Context,
	driverCtx context.Context,
	stores *store.Store,
	dbFactory *dbfactory.DBFactory,
	stateCfg *state.State,
	profile *config.Profile,
	task *store.TaskMessage,
	taskRunUID int,
	statement string,
	sheetID *int,
	mi *db.MigrationInfo,
) (string, string, error) {
	instance, err := stores.GetInstanceV2(ctx, &store.FindInstanceMessage{UID: &task.InstanceID})
	if err != nil {
		return "", "", err
	}
	database, err := stores.GetDatabaseV2(ctx, &store.FindDatabaseMessage{UID: task.DatabaseID})
	if err != nil {
		return "", "", err
	}

	driver, err := dbFactory.GetAdminDatabaseDriver(ctx, instance, database, db.ConnectionContext{})
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to get driver connection for instance %q", instance.ResourceID)
	}
	defer driver.Close(ctx)

	statementRecord, _ := common.TruncateString(statement, common.MaxSheetSize)
	slog.Debug("Start migration...",
		slog.String("instance", instance.ResourceID),
		slog.String("database", database.DatabaseName),
		slog.String("source", string(mi.Source)),
		slog.String("type", string(mi.Type)),
		slog.String("statement", statementRecord),
	)

	var migrationID string
	opts := db.ExecuteOptions{}

	opts.SetConnectionID = func(id string) {
		stateCfg.TaskRunConnectionID.Store(taskRunUID, id)
	}
	opts.DeleteConnectionID = func() {
		stateCfg.TaskRunConnectionID.Delete(taskRunUID)
	}

	if stateCfg != nil {
		switch task.Type {
		case api.TaskDatabaseSchemaUpdate, api.TaskDatabaseDataUpdate:
			switch instance.Engine {
			case storepb.Engine_MYSQL, storepb.Engine_TIDB, storepb.Engine_OCEANBASE,
				storepb.Engine_STARROCKS, storepb.Engine_DORIS, storepb.Engine_POSTGRES,
				storepb.Engine_REDSHIFT, storepb.Engine_RISINGWAVE, storepb.Engine_ORACLE,
				storepb.Engine_DM, storepb.Engine_OCEANBASE_ORACLE, storepb.Engine_MSSQL,
				storepb.Engine_DYNAMODB:
				opts.CreateTaskRunLog = getCreateTaskRunLog(ctx, taskRunUID, stores, profile)
			default:
				// do nothing
			}
		}
	}

	migrationID, schema, err := executeMigrationDefault(ctx, driverCtx, stores, stateCfg, taskRunUID, driver, mi, statement, sheetID, opts)
	if err != nil {
		return "", "", err
	}

	return migrationID, schema, nil
}

func postMigration(ctx context.Context, stores *store.Store, task *store.TaskMessage, mi *db.MigrationInfo, migrationID string, sheetID *int) (bool, *storepb.TaskRunResult, error) {
	instance, err := stores.GetInstanceV2(ctx, &store.FindInstanceMessage{UID: &task.InstanceID})
	if err != nil {
		return true, nil, err
	}
	database, err := stores.GetDatabaseV2(ctx, &store.FindDatabaseMessage{UID: task.DatabaseID})
	if err != nil {
		return true, nil, err
	}

	if mi.Type == db.Migrate || mi.Type == db.MigrateSDL {
		if _, err := stores.UpdateDatabase(ctx, &store.UpdateDatabaseMessage{
			InstanceID:    instance.ResourceID,
			DatabaseName:  database.DatabaseName,
			SchemaVersion: &mi.Version,
		}, api.SystemBotID); err != nil {
			return true, nil, errors.Errorf("failed to update database %q for instance %q", database.DatabaseName, instance.ResourceID)
		}
	}

	slog.Debug("Post migration...",
		slog.String("instance", instance.ResourceID),
		slog.String("database", database.DatabaseName),
	)

	// Set schema config.
	if sheetID != nil && task.DatabaseID != nil {
		sheet, err := stores.GetSheet(ctx, &store.FindSheetMessage{
			UID: sheetID,
		})
		if err != nil {
			slog.Error("Failed to get sheet from store", slog.Int("sheetID", *sheetID), log.BBError(err))
		} else if sheet.Payload != nil && (sheet.Payload.DatabaseConfig != nil || sheet.Payload.BaselineDatabaseConfig != nil) {
			databaseSchema, err := stores.GetDBSchema(ctx, *task.DatabaseID)
			if err != nil {
				slog.Error("Failed to get database config from store", slog.Int("sheetID", *sheetID), slog.Int("databaseUID", *task.DatabaseID), log.BBError(err))
			} else {
				updatedDatabaseConfig := utils.MergeDatabaseConfig(sheet.Payload.DatabaseConfig, sheet.Payload.BaselineDatabaseConfig, databaseSchema.GetConfig())
				err = stores.UpdateDBSchema(ctx, *task.DatabaseID, &store.UpdateDBSchemaMessage{
					Config: updatedDatabaseConfig,
				}, api.SystemBotID)
				if err != nil {
					slog.Error("Failed to update database config", slog.Int("sheetID", *sheetID), slog.Int("databaseUID", *task.DatabaseID), log.BBError(err))
				}
			}
		}
	}

	// Remove schema drift anomalies.
	if err := stores.ArchiveAnomalyV2(ctx, &store.ArchiveAnomalyMessage{
		DatabaseUID: task.DatabaseID,
		Type:        api.AnomalyDatabaseSchemaDrift,
	}); err != nil && common.ErrorCode(err) != common.NotFound {
		slog.Error("Failed to archive anomaly",
			slog.String("instance", instance.ResourceID),
			slog.String("database", database.DatabaseName),
			slog.String("type", string(api.AnomalyDatabaseSchemaDrift)),
			log.BBError(err))
	}

	detail := fmt.Sprintf("Applied migration version %s to database %q.", mi.Version.Version, database.DatabaseName)
	if mi.Type == db.Baseline {
		detail = fmt.Sprintf("Established baseline version %s for database %q.", mi.Version.Version, database.DatabaseName)
	}

	storedVersion, err := mi.Version.Marshal()
	if err != nil {
		slog.Error("failed to convert database schema version",
			slog.String("version", mi.Version.Version),
			log.BBError(err),
		)
	}
	return true, &storepb.TaskRunResult{
		Detail:        detail,
		ChangeHistory: fmt.Sprintf("instances/%s/databases/%s/changeHistories/%s", instance.ResourceID, database.DatabaseName, migrationID),
		Version:       storedVersion,
	}, nil
}

func runMigration(ctx context.Context, driverCtx context.Context, store *store.Store, dbFactory *dbfactory.DBFactory, stateCfg *state.State, profile *config.Profile, task *store.TaskMessage, taskRunUID int, migrationType db.MigrationType, statement string, schemaVersion model.Version, sheetID *int) (terminated bool, result *storepb.TaskRunResult, err error) {
	mi, err := getMigrationInfo(ctx, store, profile, task, migrationType, statement, schemaVersion, sheetID)
	if err != nil {
		return true, nil, err
	}

	migrationID, _, err := doMigration(ctx, driverCtx, store, dbFactory, stateCfg, profile, task, taskRunUID, statement, sheetID, mi)
	if err != nil {
		return true, nil, err
	}
	return postMigration(ctx, store, task, mi, migrationID, sheetID)
}

// executeMigrationDefault executes migration.
func executeMigrationDefault(ctx context.Context, driverCtx context.Context, store *store.Store, _ *state.State, taskRunUID int, driver db.Driver, mi *db.MigrationInfo, statement string, sheetID *int, opts db.ExecuteOptions) (migrationHistoryID string, updatedSchema string, resErr error) {
	execFunc := func(ctx context.Context, execStatement string) error {
		if _, err := driver.Execute(ctx, execStatement, opts); err != nil {
			return err
		}
		return nil
	}
	return executeMigrationWithFunc(ctx, driverCtx, store, taskRunUID, driver, mi, statement, sheetID, execFunc, opts)
}

// executeMigrationWithFunc executes the migration with custom migration function.
func executeMigrationWithFunc(ctx context.Context, driverCtx context.Context, s *store.Store, _ int, driver db.Driver, m *db.MigrationInfo, statement string, sheetID *int, execFunc func(ctx context.Context, execStatement string) error, opts db.ExecuteOptions) (migrationHistoryID string, updatedSchema string, resErr error) {
	var prevSchemaBuf bytes.Buffer
	if m.Type.NeedDump() {
		opts.LogSchemaDumpStart()
		// Don't record schema if the database hasn't existed yet or is schemaless, e.g. MongoDB.
		// For baseline migration, we also record the live schema to detect the schema drift.
		// See https://bytebase.com/blog/what-is-database-schema-drift
		if err := driver.Dump(ctx, &prevSchemaBuf); err != nil {
			opts.LogSchemaDumpEnd(err.Error())
			return "", "", err
		}
		opts.LogSchemaDumpEnd("")
	}

	insertedID, err := beginMigration(ctx, s, m, prevSchemaBuf.String(), statement, sheetID)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to begin migration")
	}

	startedNs := time.Now().UnixNano()

	defer func() {
		if err := endMigration(ctx, s, startedNs, insertedID, updatedSchema, prevSchemaBuf.String(), sheetID, resErr == nil /* isDone */); err != nil {
			slog.Error("Failed to update migration history record",
				log.BBError(err),
				slog.String("migration_id", migrationHistoryID),
			)
		}
	}()

	// Phase 3 - Executing migration
	// Branch migration type always has empty sql.
	// Baseline migration type could has non-empty sql but will not execute.
	// https://github.com/bytebase/bytebase/issues/394
	doMigrate := true
	if statement == "" || m.Type == db.Baseline {
		doMigrate = false
	}
	if doMigrate {
		renderedStatement := statement
		// The m.DatabaseID is nil means the migration is a instance level migration
		if m.DatabaseID != nil {
			database, err := s.GetDatabaseV2(ctx, &store.FindDatabaseMessage{
				UID: m.DatabaseID,
			})
			if err != nil {
				return "", "", err
			}
			if database == nil {
				return "", "", errors.Errorf("database %d not found", *m.DatabaseID)
			}
			materials := utils.GetSecretMapFromDatabaseMessage(database)
			// To avoid leak the rendered statement, the error message should use the original statement and not the rendered statement.
			renderedStatement = utils.RenderStatement(statement, materials)
		}

		if err := execFunc(driverCtx, renderedStatement); err != nil {
			return "", "", err
		}
	}

	// Phase 4 - Dump the schema after migration
	var afterSchemaBuf bytes.Buffer
	if m.Type.NeedDump() {
		opts.LogSchemaDumpStart()
		if err := driver.Dump(ctx, &afterSchemaBuf); err != nil {
			// We will ignore the dump error if the database is dropped.
			if strings.Contains(err.Error(), "not found") {
				return insertedID, "", nil
			}
			opts.LogSchemaDumpEnd(err.Error())
			return "", "", err
		}
		opts.LogSchemaDumpEnd("")
	}

	return insertedID, afterSchemaBuf.String(), nil
}

// beginMigration checks before executing migration and inserts a migration history record with pending status.
func beginMigration(ctx context.Context, stores *store.Store, m *db.MigrationInfo, prevSchema, statement string, sheetID *int) (string, error) {
	// Phase 1 - Pre-check before executing migration
	// Check if the same migration version has already been applied.
	if list, err := stores.ListInstanceChangeHistory(ctx, &store.FindInstanceChangeHistoryMessage{
		InstanceID: m.InstanceID,
		DatabaseID: m.DatabaseID,
		Version:    &m.Version,
	}); err != nil {
		return "", errors.Wrap(err, "failed to check duplicate version")
	} else if len(list) > 0 {
		migrationHistory := list[0]
		switch migrationHistory.Status {
		case db.Done:
			return "", common.Errorf(common.MigrationAlreadyApplied, "database %q has already applied version %s, hint: the version might be duplicate, please check the version", m.Database, m.Version.Version)
		case db.Pending:
			err := errors.Errorf("database %q version %s migration is already in progress", m.Database, m.Version.Version)
			slog.Debug(err.Error())
			// For force migration, we will ignore the existing migration history and continue to migration.
			return migrationHistory.UID, nil
		case db.Failed:
			err := errors.Errorf("database %q version %s migration has failed, please check your database to make sure things are fine and then start a new migration using a new version", m.Database, m.Version.Version)
			slog.Debug(err.Error())
			// For force migration, we will ignore the existing migration history and continue to migration.
			return migrationHistory.UID, nil
		}
	}

	// Phase 2 - Record migration history as PENDING.
	statementRecord, _ := common.TruncateString(statement, common.MaxSheetSize)
	insertedID, err := stores.CreatePendingInstanceChangeHistory(ctx, prevSchema, m, statementRecord, sheetID)
	if err != nil {
		return "", err
	}

	return insertedID, nil
}

// endMigration updates the migration history record to DONE or FAILED depending on migration is done or not.
func endMigration(ctx context.Context, storeInstance *store.Store, startedNs int64, insertedID string, updatedSchema, schemaPrev string, sheetID *int, isDone bool) error {
	migrationDurationNs := time.Now().UnixNano() - startedNs
	update := &store.UpdateInstanceChangeHistoryMessage{
		ID:                  insertedID,
		ExecutionDurationNs: &migrationDurationNs,
		// Update the sheet ID just in case it has been updated.
		Sheet: sheetID,
		// Update schemaPrev because we might be re-using a previous change history entry.
		SchemaPrev: &schemaPrev,
	}
	if isDone {
		// Upon success, update the migration history as 'DONE', execution_duration_ns, updated schema.
		status := db.Done
		update.Status = &status
		update.Schema = &updatedSchema
	} else {
		// Otherwise, update the migration history as 'FAILED', execution_duration.
		status := db.Failed
		update.Status = &status
	}
	return storeInstance.UpdateInstanceChangeHistory(ctx, update)
}
