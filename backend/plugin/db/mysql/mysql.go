// Package mysql is the plugin for MySQL driver.
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/blang/semver/v4"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/crypto/ssh"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/bytebase/bytebase/backend/common"
	"github.com/bytebase/bytebase/backend/common/log"
	"github.com/bytebase/bytebase/backend/plugin/db"
	"github.com/bytebase/bytebase/backend/plugin/db/util"
	"github.com/bytebase/bytebase/backend/plugin/parser/base"
	mysqlparser "github.com/bytebase/bytebase/backend/plugin/parser/mysql"
	storepb "github.com/bytebase/bytebase/proto/generated-go/store"
	v1pb "github.com/bytebase/bytebase/proto/generated-go/v1"
)

var (
	baseTableType = "BASE TABLE"
	viewTableType = "VIEW"

	_ db.Driver = (*Driver)(nil)
)

func init() {
	db.Register(storepb.Engine_MYSQL, newDriver)
	db.Register(storepb.Engine_MARIADB, newDriver)
	db.Register(storepb.Engine_OCEANBASE, newDriver)
}

// Driver is the MySQL driver.
type Driver struct {
	connectionCtx db.ConnectionContext
	connCfg       db.ConnectionConfig
	dbType        storepb.Engine
	dbBinDir      string
	db            *sql.DB
	databaseName  string
	sshClient     *ssh.Client
}

func newDriver(dc db.DriverConfig) db.Driver {
	return &Driver{
		dbBinDir: dc.DbBinDir,
	}
}

// Open opens a MySQL driver.
func (driver *Driver) Open(ctx context.Context, dbType storepb.Engine, connCfg db.ConnectionConfig) (db.Driver, error) {
	var dsn string
	if connCfg.AuthenticationType == storepb.DataSourceOptions_GOOGLE_CLOUD_SQL_IAM {
		connStr, err := getCloudSQLConnection(ctx, connCfg)
		if err != nil {
			return nil, err
		}
		dsn = connStr
	} else {
		connStr, err := driver.getMySQLConnection(connCfg)
		if err != nil {
			return nil, err
		}
		dsn = connStr
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	driver.dbType = dbType
	driver.db = db
	// TODO(d): remove the work-around once we have clean-up the migration connection hack.
	db.SetConnMaxLifetime(2 * time.Hour)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(15)
	driver.connectionCtx = connCfg.ConnectionContext
	driver.connCfg = connCfg
	driver.databaseName = connCfg.Database

	return driver, nil
}

func (driver *Driver) getMySQLConnection(connCfg db.ConnectionConfig) (string, error) {
	protocol := "tcp"
	if strings.HasPrefix(connCfg.Host, "/") {
		protocol = "unix"
	}

	params := []string{"multiStatements=true", "maxAllowedPacket=0"}
	if connCfg.SSHConfig.Host != "" {
		sshClient, err := util.GetSSHClient(connCfg.SSHConfig)
		if err != nil {
			return "", err
		}
		driver.sshClient = sshClient
		// Now we register the dialer with the ssh connection as a parameter.
		mysql.RegisterDialContext("mysql+tcp", func(_ context.Context, addr string) (net.Conn, error) {
			return sshClient.Dial("tcp", addr)
		})
		protocol = "mysql+tcp"
	}

	// TODO(zp): mysql and mysqlbinlog doesn't support SSL yet. We need to write certs to temp files and load them as CLI flags.
	tlsConfig, err := connCfg.TLSConfig.GetSslConfig()
	if err != nil {
		return "", errors.Wrap(err, "sql: tls config error")
	}
	tlsKey := "db.mysql.tls"
	if tlsConfig != nil {
		if err := mysql.RegisterTLSConfig(tlsKey, tlsConfig); err != nil {
			return "", errors.Wrap(err, "sql: failed to register tls config")
		}
		// TLS config is only used during sql.Open, so should be safe to deregister afterwards.
		defer mysql.DeregisterTLSConfig(tlsKey)
		params = append(params, fmt.Sprintf("tls=%s", tlsKey))
	}

	return fmt.Sprintf("%s:%s@%s(%s:%s)/%s?%s", connCfg.Username, connCfg.Password, protocol, connCfg.Host, connCfg.Port, connCfg.Database, strings.Join(params, "&")), nil
}

func getCloudSQLConnection(ctx context.Context, connCfg db.ConnectionConfig) (string, error) {
	d, err := cloudsqlconn.NewDialer(ctx, cloudsqlconn.WithIAMAuthN())
	if err != nil {
		return "", err
	}
	mysql.RegisterDialContext("cloudsqlconn",
		func(ctx context.Context, _ string) (net.Conn, error) {
			return d.Dial(ctx, connCfg.Host)
		})

	return fmt.Sprintf("%s:empty@cloudsqlconn(localhost:3306)/%s?parseTime=true",
		connCfg.Username, connCfg.Database), nil
}

// Close closes the driver.
func (driver *Driver) Close(context.Context) error {
	var err error
	err = multierr.Append(err, driver.db.Close())
	if driver.sshClient != nil {
		err = multierr.Append(err, driver.sshClient.Close())
	}
	return err
}

// Ping pings the database.
func (driver *Driver) Ping(ctx context.Context) error {
	return driver.db.PingContext(ctx)
}

// GetType returns the database type.
func (driver *Driver) GetType() storepb.Engine {
	return driver.dbType
}

// GetDB gets the database.
func (driver *Driver) GetDB() *sql.DB {
	return driver.db
}

// getVersion gets the version.
func (driver *Driver) getVersion(ctx context.Context) (string, string, error) {
	query := "SELECT VERSION()"
	var version string
	if err := driver.db.QueryRowContext(ctx, query).Scan(&version); err != nil {
		if err == sql.ErrNoRows {
			return "", "", common.FormatDBErrorEmptyRowWithQuery(query)
		}
		return "", "", util.FormatErrorWithQuery(err, query)
	}

	return parseVersion(version)
}

func (driver *Driver) getReadOnly() bool {
	if driver.dbType == storepb.Engine_OCEANBASE {
		return false
	}
	// MariaDB 5.5 doesn't support READ ONLY transactions.
	// Error 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MariaDB server version for the right syntax to use near 'READ ONLY' at line 1
	v, err := semver.Make(driver.connectionCtx.EngineVersion)
	if err != nil {
		slog.Debug("invalid version", slog.String("version", driver.connectionCtx.EngineVersion))
		return true
	}
	if v.GT(semver.Version{Major: 5, Minor: 5}) {
		return true
	}
	return false
}

func parseVersion(version string) (string, string, error) {
	if loc := regexp.MustCompile(`^\d+.\d+.\d+`).FindStringIndex(version); loc != nil {
		return version[loc[0]:loc[1]], version[loc[1]:], nil
	}
	return "", "", errors.Errorf("failed to parse version %q", version)
}

// Execute executes a SQL statement.
func (driver *Driver) Execute(ctx context.Context, statement string, opts db.ExecuteOptions) (int64, error) {
	statement, err := mysqlparser.DealWithDelimiter(statement)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to deal with delimiter")
	}

	conn, err := driver.db.Conn(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	connectionID, err := getConnectionID(ctx, conn)
	if err != nil {
		return 0, err
	}
	slog.Debug("connectionID", slog.String("connectionID", connectionID))

	if opts.BeginFunc != nil {
		if err := opts.BeginFunc(ctx, conn); err != nil {
			return 0, err
		}
	}

	var totalCommands int
	var chunks [][]base.SingleSQL
	if opts.ChunkedSubmission && len(statement) <= common.MaxSheetCheckSize {
		singleSQLs, err := mysqlparser.SplitSQL(statement)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to split sql")
		}
		singleSQLs = base.FilterEmptySQL(singleSQLs)
		if len(singleSQLs) == 0 {
			return 0, nil
		}
		totalCommands = len(singleSQLs)
		ret, err := util.ChunkedSQLScript(singleSQLs, common.MaxSheetChunksCount)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to chunk sql")
		}
		chunks = ret
	} else {
		chunks = [][]base.SingleSQL{
			{
				base.SingleSQL{
					Text: statement,
				},
			},
		}
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to begin execute transaction")
	}
	defer tx.Rollback()

	currentIndex := 0
	var totalRowsAffected int64
	for _, chunk := range chunks {
		if len(chunk) == 0 {
			continue
		}
		// Start the current chunk.

		// Set the progress information for the current chunk.
		if opts.UpdateExecutionStatus != nil {
			opts.UpdateExecutionStatus(&v1pb.TaskRun_ExecutionDetail{
				CommandsTotal:     int32(totalCommands),
				CommandsCompleted: int32(currentIndex),
				CommandStartPosition: &v1pb.TaskRun_ExecutionDetail_Position{
					Line:   int32(chunk[0].FirstStatementLine),
					Column: int32(chunk[0].FirstStatementColumn),
				},
				CommandEndPosition: &v1pb.TaskRun_ExecutionDetail_Position{
					Line:   int32(chunk[len(chunk)-1].LastLine),
					Column: int32(chunk[len(chunk)-1].LastColumn),
				},
			})
		}

		chunkText, err := util.ConcatChunk(chunk)
		if err != nil {
			return 0, err
		}

		sqlResult, err := tx.ExecContext(ctx, chunkText)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				slog.Info("cancel connection", slog.String("connectionID", connectionID))
				if err := driver.StopConnectionByID(connectionID); err != nil {
					slog.Error("failed to cancel connection", slog.String("connectionID", connectionID), log.BBError(err))
				}
			}

			return 0, &db.ErrorWithPosition{
				Err: errors.Wrapf(err, "failed to execute context in a transaction"),
				Start: &storepb.TaskRunResult_Position{
					Line:   int32(chunk[0].FirstStatementLine),
					Column: int32(chunk[0].FirstStatementColumn),
				},
				End: &storepb.TaskRunResult_Position{
					Line:   int32(chunk[len(chunk)-1].LastLine),
					Column: int32(chunk[len(chunk)-1].LastColumn),
				},
			}
		}
		rowsAffected, err := sqlResult.RowsAffected()
		if err != nil {
			// Since we cannot differentiate DDL and DML yet, we have to ignore the error.
			slog.Debug("rowsAffected returns error", log.BBError(err))
		}
		totalRowsAffected += rowsAffected
		currentIndex += len(chunk)
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.Wrapf(err, "failed to commit execute transaction")
	}

	return totalRowsAffected, nil
}

// QueryConn queries a SQL statement in a given connection.
func (driver *Driver) QueryConn(ctx context.Context, conn *sql.Conn, statement string, queryContext *db.QueryContext) ([]*v1pb.QueryResult, error) {
	if queryContext.ReadOnly {
		queryContext.ReadOnly = driver.getReadOnly()
	}

	singleSQLs, err := base.SplitMultiSQL(storepb.Engine_MYSQL, statement)
	if err != nil {
		return nil, err
	}
	singleSQLs = base.FilterEmptySQL(singleSQLs)
	if len(singleSQLs) == 0 {
		return nil, nil
	}

	connectionID, err := getConnectionID(ctx, conn)
	if err != nil {
		return nil, err
	}
	slog.Debug("connectionID", slog.String("connectionID", connectionID))
	var results []*v1pb.QueryResult
	for _, singleSQL := range singleSQLs {
		result, err := driver.querySingleSQL(ctx, conn, singleSQL, queryContext)
		if err != nil {
			results = append(results, &v1pb.QueryResult{
				Error: err.Error(),
			})
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				slog.Info("cancel connection", slog.String("connectionID", connectionID))
				if err := driver.StopConnectionByID(connectionID); err != nil {
					slog.Error("failed to cancel connection", slog.String("connectionID", connectionID), log.BBError(err))
				}
				break
			}
		} else {
			results = append(results, result)
		}
	}

	return results, nil
}

func (driver *Driver) StopConnectionByID(id string) error {
	// We cannot use placeholder parameter because TiDB doesn't accept it.
	_, err := driver.db.Exec(fmt.Sprintf("KILL QUERY %s", id))
	return err
}

func getConnectionID(ctx context.Context, conn *sql.Conn) (string, error) {
	var id string
	if err := conn.QueryRowContext(ctx, `SELECT CONNECTION_ID();`).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (driver *Driver) querySingleSQL(ctx context.Context, conn *sql.Conn, singleSQL base.SingleSQL, queryContext *db.QueryContext) (*v1pb.QueryResult, error) {
	statement := strings.TrimLeft(strings.TrimRight(singleSQL.Text, " \n\t;"), " \n\t")
	isExplain := strings.HasPrefix(statement, "EXPLAIN")
	isSet, _ := regexp.MatchString(`(?i)^SET\s+?`, statement)

	stmt := statement
	if !isExplain && !isSet && queryContext.Limit > 0 {
		stmt = getStatementWithResultLimit(stmt, queryContext.Limit)
	}

	if queryContext.SensitiveSchemaInfo != nil {
		for _, database := range queryContext.SensitiveSchemaInfo.DatabaseList {
			if len(database.SchemaList) == 0 {
				continue
			}
			if len(database.SchemaList) > 1 {
				return nil, errors.Errorf("MySQL schema info should only have one schema per database, but got %d, %v", len(database.SchemaList), database.SchemaList)
			}
			if database.SchemaList[0].Name != "" {
				return nil, errors.Errorf("MySQL schema info should have empty schema name, but got %s", database.SchemaList[0].Name)
			}
		}
	}

	startTime := time.Now()
	result, err := util.Query(ctx, driver.dbType, conn, stmt, queryContext)
	if err != nil {
		return nil, err
	}
	result.Latency = durationpb.New(time.Since(startTime))
	result.Statement = statement
	return result, nil
}

// RunStatement runs a SQL statement in a given connection.
func (*Driver) RunStatement(ctx context.Context, conn *sql.Conn, statement string) ([]*v1pb.QueryResult, error) {
	return util.RunStatement(ctx, storepb.Engine_MYSQL, conn, statement)
}
