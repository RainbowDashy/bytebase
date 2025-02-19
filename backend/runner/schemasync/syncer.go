// Package schemasync is a runner that synchronize database schemas.
package schemasync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bytebase/bytebase/backend/common"
	"github.com/bytebase/bytebase/backend/common/log"
	"github.com/bytebase/bytebase/backend/component/config"
	"github.com/bytebase/bytebase/backend/component/dbfactory"
	"github.com/bytebase/bytebase/backend/component/state"
	enterpriseAPI "github.com/bytebase/bytebase/backend/enterprise/api"
	api "github.com/bytebase/bytebase/backend/legacyapi"
	"github.com/bytebase/bytebase/backend/plugin/db"
	"github.com/bytebase/bytebase/backend/store"
	"github.com/bytebase/bytebase/backend/utils"
	storepb "github.com/bytebase/bytebase/proto/generated-go/store"
)

const (
	schemaSyncInterval  = 1 * time.Minute
	defaultSyncInterval = 24 * time.Hour
)

// NewSyncer creates a schema syncer.
func NewSyncer(store *store.Store, dbFactory *dbfactory.DBFactory, stateCfg *state.State, profile config.Profile, licenseService enterpriseAPI.LicenseService) *Syncer {
	return &Syncer{
		store:          store,
		dbFactory:      dbFactory,
		stateCfg:       stateCfg,
		profile:        profile,
		licenseService: licenseService,
	}
}

// Syncer is the schema syncer.
type Syncer struct {
	store          *store.Store
	dbFactory      *dbfactory.DBFactory
	stateCfg       *state.State
	profile        config.Profile
	licenseService enterpriseAPI.LicenseService
}

// Run will run the schema syncer once.
func (s *Syncer) Run(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(schemaSyncInterval)
	defer ticker.Stop()
	defer wg.Done()
	log.Debug(fmt.Sprintf("Schema syncer started and will run every %v", schemaSyncInterval))
	for {
		select {
		case <-ticker.C:
			s.trySyncAll(ctx)
		case instance := <-s.stateCfg.InstanceDatabaseSyncChan:
			// Sync all databases for instance.
			s.syncAllDatabases(ctx, instance)
		case <-ctx.Done(): // if cancel() execute
			return
		}
	}
}

func (s *Syncer) trySyncAll(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = errors.Errorf("%v", r)
			}
			log.Error("Instance syncer PANIC RECOVER", zap.Error(err), zap.Stack("panic-stack"))
		}
	}()
	instances, err := s.store.ListInstancesV2(ctx, &store.FindInstanceMessage{})
	if err != nil {
		log.Error("Failed to retrieve instances", zap.Error(err))
		return
	}

	now := time.Now()
	for _, instance := range instances {
		interval := getOrDefaultSyncInterval(instance)
		lastSyncTime := getOrDefaultLastSyncTime(instance.Metadata.LastSyncTime)
		// lastSyncTime + syncInterval > now
		// Next round not started yet.
		nextSyncTime := lastSyncTime.Add(interval)
		if now.Before(nextSyncTime) {
			continue
		}

		log.Debug("Sync instance schema", zap.String("instance", instance.ResourceID))
		if err := s.SyncInstance(ctx, instance); err != nil {
			log.Debug("Failed to sync instance",
				zap.String("instance", instance.ResourceID),
				zap.String("error", err.Error()))
		}
	}

	instancesMap := map[string]*store.InstanceMessage{}
	for _, instance := range instances {
		instancesMap[instance.ResourceID] = instance
	}

	databases, err := s.store.ListDatabases(ctx, &store.FindDatabaseMessage{})
	if err != nil {
		log.Error("Failed to retrieve databases", zap.Error(err))
		return
	}
	for _, database := range databases {
		instance, ok := instancesMap[database.InstanceID]
		if !ok {
			continue
		}
		// The database inherits the sync interval from the instance.
		interval := getOrDefaultSyncInterval(instance)
		lastSyncTime := getOrDefaultLastSyncTime(database.Metadata.LastSyncTime)
		// lastSyncTime + syncInterval > now
		// Next round not started yet.
		nextSyncTime := lastSyncTime.Add(interval)
		if now.Before(nextSyncTime) {
			continue
		}
		if err := s.SyncDatabaseSchema(ctx, database, false /* force */); err != nil {
			log.Debug("Failed to sync database schema",
				zap.String("instance", instance.ResourceID),
				zap.String("databaseName", database.DatabaseName),
				zap.Error(err))
		}
	}
}

func (s *Syncer) syncAllDatabases(ctx context.Context, instance *store.InstanceMessage) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = errors.Errorf("%v", r)
			}
			log.Error("Database syncer PANIC RECOVER", zap.Error(err), zap.Stack("panic-stack"))
		}
	}()

	find := &store.FindDatabaseMessage{}
	if instance != nil {
		find.InstanceID = &instance.ResourceID
	}
	databases, err := s.store.ListDatabases(ctx, find)
	if err != nil {
		log.Debug("Failed to find databases to sync",
			zap.String("error", err.Error()))
		return
	}

	instanceMap := make(map[string][]*store.DatabaseMessage)
	for _, database := range databases {
		// Skip deleted databases.
		if database.SyncState != api.OK {
			continue
		}
		instanceMap[database.InstanceID] = append(instanceMap[database.InstanceID], database)
	}

	var instanceWG sync.WaitGroup
	for _, databaseList := range instanceMap {
		instanceWG.Add(1)
		go func(databaseList []*store.DatabaseMessage) {
			defer instanceWG.Done()

			if len(databaseList) == 0 {
				return
			}
			for _, database := range databaseList {
				instanceID := database.InstanceID
				log.Debug("Sync database schema",
					zap.String("instance", instanceID),
					zap.String("database", database.DatabaseName),
					zap.Int64("lastSuccessfulSyncTs", database.SuccessfulSyncTimeTs),
				)
				// If we fail to sync a particular database due to permission issue, we will continue to sync the rest of the databases.
				// We don't force dump database schema because it's rarely changed till the metadata is changed.
				if err := s.SyncDatabaseSchema(ctx, database, false /* force */); err != nil {
					log.Debug("Failed to sync database schema",
						zap.String("instance", instanceID),
						zap.String("databaseName", database.DatabaseName),
						zap.Error(err))
				}
			}
		}(databaseList)
	}
	instanceWG.Wait()
}

// SyncInstance syncs the schema for all databases in an instance.
func (s *Syncer) SyncInstance(ctx context.Context, instance *store.InstanceMessage) error {
	driver, err := s.dbFactory.GetAdminDatabaseDriver(ctx, instance, nil /* database */)
	if err != nil {
		return err
	}
	defer driver.Close(ctx)

	instanceMeta, err := driver.SyncInstance(ctx)
	if err != nil {
		return err
	}

	updateInstance := &store.UpdateInstanceMessage{
		UpdaterID:     api.SystemBotID,
		EnvironmentID: instance.EnvironmentID,
		ResourceID:    instance.ResourceID,
		Metadata: &storepb.InstanceMetadata{
			LastSyncTime: timestamppb.Now(),
		},
	}
	if instanceMeta.Version != instance.EngineVersion {
		updateInstance.EngineVersion = &instanceMeta.Version
	}
	if !equalInstanceMetadata(instanceMeta.Metadata, instance.Metadata) {
		updateInstance.Metadata.MysqlLowerCaseTableNames = instanceMeta.Metadata.GetMysqlLowerCaseTableNames()
	}
	if _, err := s.store.UpdateInstanceV2(ctx, updateInstance, -1); err != nil {
		return err
	}

	var instanceUsers []*store.InstanceUserMessage
	for _, instanceUser := range instanceMeta.InstanceRoles {
		instanceUsers = append(instanceUsers, &store.InstanceUserMessage{
			Name:  instanceUser.Name,
			Grant: instanceUser.Grant,
		})
	}
	if err := s.store.UpsertInstanceUsers(ctx, instance.UID, instanceUsers); err != nil {
		return err
	}

	databases, err := s.store.ListDatabases(ctx, &store.FindDatabaseMessage{InstanceID: &instance.ResourceID})
	if err != nil {
		return errors.Wrapf(err, "failed to sync database for instance: %s. Failed to find database list", instance.ResourceID)
	}
	for _, databaseMetadata := range instanceMeta.Databases {
		exist := false
		for _, database := range databases {
			if database.DatabaseName == databaseMetadata.Name {
				exist = true
				break
			}
		}
		if !exist {
			// Create the database in the default project.
			if err := s.store.CreateDatabaseDefault(ctx, &store.DatabaseMessage{
				InstanceID:   instance.ResourceID,
				DatabaseName: databaseMetadata.Name,
				DataShare:    databaseMetadata.Datashare,
				ServiceName:  databaseMetadata.ServiceName,
				ProjectID:    api.DefaultProjectID,
			}); err != nil {
				return errors.Wrapf(err, "failed to create instance %q database %q in sync runner", instance.ResourceID, databaseMetadata.Name)
			}
		}
	}

	for _, database := range databases {
		exist := false
		for _, databaseMetadata := range instanceMeta.Databases {
			if database.DatabaseName == databaseMetadata.Name {
				exist = true
				break
			}
		}
		if !exist {
			syncStatus := api.NotFound
			if _, err := s.store.UpdateDatabase(ctx, &store.UpdateDatabaseMessage{
				InstanceID:   instance.ResourceID,
				DatabaseName: database.DatabaseName,
				SyncState:    &syncStatus,
			}, api.SystemBotID); err != nil {
				return errors.Errorf("failed to update database %q for instance %q", database.DatabaseName, instance.ResourceID)
			}
		}
	}

	return nil
}

// SyncDatabaseSchema will sync the schema for a database.
func (s *Syncer) SyncDatabaseSchema(ctx context.Context, database *store.DatabaseMessage, force bool) error {
	instance, err := s.store.GetInstanceV2(ctx, &store.FindInstanceMessage{ResourceID: &database.InstanceID})
	if err != nil {
		return err
	}
	if instance == nil {
		return errors.Errorf("instance %q not found", database.InstanceID)
	}
	driver, err := s.dbFactory.GetAdminDatabaseDriver(ctx, instance, database)
	if err != nil {
		return err
	}
	defer driver.Close(ctx)
	// Sync database schema
	databaseMetadata, err := driver.SyncDBSchema(ctx)
	if err != nil {
		return err
	}
	setClassificationAndUserCommentFromComment(databaseMetadata)

	var patchSchemaVersion *string
	if force {
		// When there are too many databases, this might have performance issue and will
		// cause frontend timeout since we set a 30s limit (INSTANCE_OPERATION_TIMEOUT).
		schemaVersion, err := utils.GetLatestSchemaVersion(ctx, s.store, instance.UID, database.UID, databaseMetadata.Name)
		if err != nil {
			return err
		}
		patchSchemaVersion = &schemaVersion
	}

	syncStatus := api.OK
	ts := time.Now().Unix()
	if _, err := s.store.UpdateDatabase(ctx, &store.UpdateDatabaseMessage{
		InstanceID:           database.InstanceID,
		DatabaseName:         database.DatabaseName,
		SyncState:            &syncStatus,
		SuccessfulSyncTimeTs: &ts,
		SchemaVersion:        patchSchemaVersion,
		MetadataUpsert: &storepb.DatabaseMetadata{
			LastSyncTime: timestamppb.New(time.Unix(ts, 0)),
		},
	}, api.SystemBotID); err != nil {
		return errors.Wrapf(err, "failed to update database %q for instance %q", database.DatabaseName, database.InstanceID)
	}

	dbSchema, err := s.store.GetDBSchema(ctx, database.UID)
	if err != nil {
		return err
	}
	var oldDatabaseMetadata *storepb.DatabaseSchemaMetadata
	var rawDump []byte
	if dbSchema != nil {
		oldDatabaseMetadata = dbSchema.Metadata
		rawDump = dbSchema.Schema
	}

	if !cmp.Equal(oldDatabaseMetadata, databaseMetadata, protocmp.Transform()) {
		// Avoid updating dump everytime by dumping the schema only when the database metadata is changed.
		// if oldDatabaseMetadata is nil and databaseMetadata is not, they are not equal resulting a sync.
		if force || !equalDatabaseMetadata(oldDatabaseMetadata, databaseMetadata) {
			var schemaBuf bytes.Buffer
			if _, err := driver.Dump(ctx, &schemaBuf, true /* schemaOnly */); err != nil {
				return err
			}
			rawDump = schemaBuf.Bytes()
		}

		if err := s.store.UpsertDBSchema(ctx, database.UID, &store.DBSchema{
			Metadata: databaseMetadata,
			Schema:   rawDump,
		}, api.SystemBotID); err != nil {
			return err
		}
	}

	// Check schema drift
	if s.licenseService.IsFeatureEnabledForInstance(api.FeatureSchemaDrift, instance) == nil {
		// Redis and MongoDB are schemaless.
		if disableSchemaDriftAnomalyCheck(instance.Engine) {
			return nil
		}
		limit := 1
		list, err := s.store.FindInstanceChangeHistoryList(ctx, &db.MigrationHistoryFind{
			InstanceID: &instance.UID,
			DatabaseID: &database.UID,
			Database:   &database.DatabaseName,
			Limit:      &limit,
		})
		if err != nil {
			log.Error("Failed to check anomaly",
				zap.String("instance", instance.ResourceID),
				zap.String("database", database.DatabaseName),
				zap.String("type", string(api.AnomalyDatabaseSchemaDrift)),
				zap.Error(err))
			return nil
		}
		latestSchema := string(rawDump)
		if len(list) > 0 {
			if list[0].Schema != latestSchema {
				anomalyPayload := api.AnomalyDatabaseSchemaDriftPayload{
					Version: list[0].Version,
					Expect:  list[0].Schema,
					Actual:  latestSchema,
				}
				payload, err := json.Marshal(anomalyPayload)
				if err != nil {
					log.Error("Failed to marshal anomaly payload",
						zap.String("instance", instance.ResourceID),
						zap.String("database", database.DatabaseName),
						zap.String("type", string(api.AnomalyDatabaseSchemaDrift)),
						zap.Error(err))
				} else {
					if _, err = s.store.UpsertActiveAnomalyV2(ctx, api.SystemBotID, &store.AnomalyMessage{
						InstanceID:  instance.ResourceID,
						DatabaseUID: &database.UID,
						Type:        api.AnomalyDatabaseSchemaDrift,
						Payload:     string(payload),
					}); err != nil {
						log.Error("Failed to create anomaly",
							zap.String("instance", instance.ResourceID),
							zap.String("database", database.DatabaseName),
							zap.String("type", string(api.AnomalyDatabaseSchemaDrift)),
							zap.Error(err))
					}
				}
			} else {
				err := s.store.ArchiveAnomalyV2(ctx, &store.ArchiveAnomalyMessage{
					DatabaseUID: &database.UID,
					Type:        api.AnomalyDatabaseSchemaDrift,
				})
				if err != nil && common.ErrorCode(err) != common.NotFound {
					log.Error("Failed to close anomaly",
						zap.String("instance", instance.ResourceID),
						zap.String("database", database.DatabaseName),
						zap.String("type", string(api.AnomalyDatabaseSchemaDrift)),
						zap.Error(err))
				}
			}
		}
	}
	return nil
}

func equalInstanceMetadata(x, y *storepb.InstanceMetadata) bool {
	return cmp.Equal(x, y, protocmp.Transform(), protocmp.IgnoreFields(&storepb.InstanceMetadata{}, "last_sync_time"))
}

func equalDatabaseMetadata(x, y *storepb.DatabaseSchemaMetadata) bool {
	return cmp.Equal(x, y, protocmp.Transform(),
		protocmp.IgnoreFields(&storepb.TableMetadata{}, "row_count", "data_size", "index_size", "data_free"),
	)
}

func setClassificationAndUserCommentFromComment(dbSchema *storepb.DatabaseSchemaMetadata) {
	for _, schema := range dbSchema.Schemas {
		for _, table := range schema.Tables {
			table.Classification, table.UserComment = common.GetClassificationAndUserComment(table.Comment)
			for _, col := range table.Columns {
				col.Classification, col.UserComment = common.GetClassificationAndUserComment(col.Comment)
			}
		}
	}
}

func getOrDefaultSyncInterval(instance *store.InstanceMessage) time.Duration {
	if instance.Activation && instance.Options.SyncInterval.IsValid() {
		return instance.Options.SyncInterval.AsDuration()
	}
	return defaultSyncInterval
}

func getOrDefaultLastSyncTime(t *timestamppb.Timestamp) time.Time {
	if t.IsValid() {
		return t.AsTime()
	}
	return time.Unix(0, 0)
}

func disableSchemaDriftAnomalyCheck(dbTp db.Type) bool {
	m := map[db.Type]struct{}{
		db.MongoDB:  {},
		db.Redis:    {},
		db.Oracle:   {},
		db.MSSQL:    {},
		db.Redshift: {},
	}
	_, ok := m[dbTp]
	return ok
}
