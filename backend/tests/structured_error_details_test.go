package tests

import (
	"context"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"

	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
	v1pb "github.com/bytebase/bytebase/backend/generated-go/v1"
)

func TestStructuredErrorDetails_SyntaxError(t *testing.T) {
	tests := []struct {
		name              string
		databaseName      string
		dbType            storepb.Engine
		prepareStatements string
		query             string
		wantSyntaxError   bool
		wantPosition      bool
	}{
		{
			name:              "MySQL - Syntax error populates SyntaxErrorDetail",
			databaseName:      "TestMySQLSyntaxError",
			dbType:            storepb.Engine_MYSQL,
			prepareStatements: "CREATE TABLE tbl(id INT PRIMARY KEY, name VARCHAR(64));",
			query:             "SELCT * FROM tbl;", // Intentional typo: SELCT instead of SELECT
			wantSyntaxError:   true,
			wantPosition:      true,
		},
		{
			name:              "PostgreSQL - Syntax error populates SyntaxErrorDetail",
			databaseName:      "TestPostgresSyntaxError",
			dbType:            storepb.Engine_POSTGRES,
			prepareStatements: "CREATE TABLE tbl(id INT PRIMARY KEY, name VARCHAR(64));",
			query:             "SELCT * FROM tbl;", // Intentional typo
			wantSyntaxError:   true,
			wantPosition:      true,
		},
		{
			name:              "MySQL - Multiple statements with syntax error in second",
			databaseName:      "TestMySQLMultipleSyntaxError",
			dbType:            storepb.Engine_MYSQL,
			prepareStatements: "CREATE TABLE tbl(id INT PRIMARY KEY);",
			query:             "INSERT INTO tbl VALUES(1); SELCT * FROM tbl;",
			wantSyntaxError:   true,
			wantPosition:      true,
		},
		{
			name:              "PostgreSQL - Invalid column reference",
			databaseName:      "TestPostgresInvalidColumn",
			dbType:            storepb.Engine_POSTGRES,
			prepareStatements: "CREATE TABLE tbl(id INT PRIMARY KEY);",
			query:             "SELECT nonexistent FROM tbl;",
			wantSyntaxError:   false, // This is a semantic error, not syntax
			wantPosition:      false,
		},
	}

	t.Parallel()
	a := require.New(t)
	ctx := context.Background()
	ctl := &controller{}
	ctx, err := ctl.StartServerWithExternalPg(ctx)
	a.NoError(err)
	defer ctl.Close(ctx)

	mysqlContainer, err := getMySQLContainer(ctx)
	a.NoError(err)
	t.Cleanup(func() {
		mysqlContainer.Close(ctx)
	})

	pgContainer, err := getPgContainer(ctx)
	a.NoError(err)
	t.Cleanup(func() {
		pgContainer.Close(ctx)
	})

	mysqlInstanceResp, err := ctl.instanceServiceClient.CreateInstance(ctx, connect.NewRequest(&v1pb.CreateInstanceRequest{
		InstanceId: generateRandomString("instance"),
		Instance: &v1pb.Instance{
			Title:       "mysqlInstance",
			Engine:      v1pb.Engine_MYSQL,
			Environment: stringPtr("environments/prod"),
			Activation:  true,
			DataSources: []*v1pb.DataSource{{Type: v1pb.DataSourceType_ADMIN, Host: mysqlContainer.host, Port: mysqlContainer.port, Username: "root", Password: "root-password", Id: "admin"}},
		},
	}))
	a.NoError(err)
	mysqlInstance := mysqlInstanceResp.Msg

	pgInstanceResp, err := ctl.instanceServiceClient.CreateInstance(ctx, connect.NewRequest(&v1pb.CreateInstanceRequest{
		InstanceId: generateRandomString("instance"),
		Instance: &v1pb.Instance{
			Title:       "pgInstance",
			Engine:      v1pb.Engine_POSTGRES,
			Environment: stringPtr("environments/prod"),
			Activation:  true,
			DataSources: []*v1pb.DataSource{{Type: v1pb.DataSourceType_ADMIN, Host: pgContainer.host, Port: pgContainer.port, Username: "postgres", Password: "root-password", Id: "admin"}},
		},
	}))
	a.NoError(err)
	pgInstance := pgInstanceResp.Msg

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := require.New(t)

			var instance *v1pb.Instance
			databaseOwner := ""
			switch tt.dbType {
			case storepb.Engine_MYSQL:
				instance = mysqlInstance
			case storepb.Engine_POSTGRES:
				instance = pgInstance
				databaseOwner = "postgres"
			default:
				a.FailNow("unsupported db type")
			}

			err = ctl.createDatabaseV2(ctx, ctl.project, instance, nil, tt.databaseName, databaseOwner)
			a.NoError(err)

			databaseResp, err := ctl.databaseServiceClient.GetDatabase(ctx, connect.NewRequest(&v1pb.GetDatabaseRequest{
				Name: fmt.Sprintf("%s/databases/%s", instance.Name, tt.databaseName),
			}))
			a.NoError(err)
			database := databaseResp.Msg

			sheetResp, err := ctl.sheetServiceClient.CreateSheet(ctx, connect.NewRequest(&v1pb.CreateSheetRequest{
				Parent: ctl.project.Name,
				Sheet: &v1pb.Sheet{
					Title:   "prepareStatements",
					Content: []byte(tt.prepareStatements),
				},
			}))
			a.NoError(err)
			sheet := sheetResp.Msg

			a.NotNil(database.InstanceResource)
			a.Equal(1, len(database.InstanceResource.DataSources))

			err = ctl.changeDatabase(ctx, ctl.project, database, sheet, v1pb.MigrationType_DDL)
			a.NoError(err)

			// Execute query that should produce a syntax error
			queryResp, err := ctl.sqlServiceClient.Query(ctx, connect.NewRequest(&v1pb.QueryRequest{
				Name:         database.Name,
				Statement:    tt.query,
				DataSourceId: "admin",
			}))

			if tt.wantSyntaxError {
				// Query should return results with error details, not a connect error
				a.NoError(err, "Query should return results with error details")
				a.NotNil(queryResp)
				a.NotEmpty(queryResp.Msg.Results, "Should have at least one result with error")

				// Find the result with syntax error
				var foundSyntaxError bool
				for _, result := range queryResp.Msg.Results {
					if result.GetSyntaxError() != nil {
						foundSyntaxError = true

						// Verify error string is populated (backward compatibility)
						a.NotEmpty(result.Error, "Error string should be populated for backward compatibility")

						// Verify SyntaxErrorDetail is populated
						syntaxError := result.GetSyntaxError()
						a.NotNil(syntaxError, "SyntaxErrorDetail should be populated")

						if tt.wantPosition {
							// Verify position information is present
							a.NotNil(syntaxError.Position, "Position should be present")
							a.Greater(syntaxError.Position.Line, int32(0), "Line should be greater than 0")
							a.GreaterOrEqual(syntaxError.Position.Column, int32(0), "Column should be >= 0")
						}

						// Verify statement is set
						a.NotEmpty(result.Statement, "Statement should be set in result")
						break
					}
				}
				a.True(foundSyntaxError, "Should have found a result with SyntaxErrorDetail")
			} else if err == nil {
				// For non-syntax errors, we might still get an error but not SyntaxErrorDetail
				a.NotNil(queryResp)
				for _, result := range queryResp.Msg.Results {
					a.Nil(result.GetSyntaxError(), "Should not have SyntaxErrorDetail for non-syntax errors")
				}
			}
		})
	}
}

func TestStructuredErrorDetails_BackwardCompatibility(t *testing.T) {
	tests := []struct {
		name              string
		databaseName      string
		dbType            storepb.Engine
		prepareStatements string
		query             string
		wantErrorString   bool
	}{
		{
			name:              "MySQL - Error string always populated with SyntaxErrorDetail",
			databaseName:      "TestBackwardCompatMySQL",
			dbType:            storepb.Engine_MYSQL,
			prepareStatements: "CREATE TABLE tbl(id INT PRIMARY KEY);",
			query:             "SELCT * FROM tbl;",
			wantErrorString:   true,
		},
		{
			name:              "PostgreSQL - Error string always populated with SyntaxErrorDetail",
			databaseName:      "TestBackwardCompatPostgres",
			dbType:            storepb.Engine_POSTGRES,
			prepareStatements: "CREATE TABLE tbl(id INT PRIMARY KEY);",
			query:             "SELCT * FROM tbl;",
			wantErrorString:   true,
		},
	}

	t.Parallel()
	a := require.New(t)
	ctx := context.Background()
	ctl := &controller{}
	ctx, err := ctl.StartServerWithExternalPg(ctx)
	a.NoError(err)
	defer ctl.Close(ctx)

	mysqlContainer, err := getMySQLContainer(ctx)
	a.NoError(err)
	t.Cleanup(func() {
		mysqlContainer.Close(ctx)
	})

	pgContainer, err := getPgContainer(ctx)
	a.NoError(err)
	t.Cleanup(func() {
		pgContainer.Close(ctx)
	})

	mysqlInstanceResp, err := ctl.instanceServiceClient.CreateInstance(ctx, connect.NewRequest(&v1pb.CreateInstanceRequest{
		InstanceId: generateRandomString("instance"),
		Instance: &v1pb.Instance{
			Title:       "mysqlInstance",
			Engine:      v1pb.Engine_MYSQL,
			Environment: stringPtr("environments/prod"),
			Activation:  true,
			DataSources: []*v1pb.DataSource{{Type: v1pb.DataSourceType_ADMIN, Host: mysqlContainer.host, Port: mysqlContainer.port, Username: "root", Password: "root-password", Id: "admin"}},
		},
	}))
	a.NoError(err)
	mysqlInstance := mysqlInstanceResp.Msg

	pgInstanceResp, err := ctl.instanceServiceClient.CreateInstance(ctx, connect.NewRequest(&v1pb.CreateInstanceRequest{
		InstanceId: generateRandomString("instance"),
		Instance: &v1pb.Instance{
			Title:       "pgInstance",
			Engine:      v1pb.Engine_POSTGRES,
			Environment: stringPtr("environments/prod"),
			Activation:  true,
			DataSources: []*v1pb.DataSource{{Type: v1pb.DataSourceType_ADMIN, Host: pgContainer.host, Port: pgContainer.port, Username: "postgres", Password: "root-password", Id: "admin"}},
		},
	}))
	a.NoError(err)
	pgInstance := pgInstanceResp.Msg

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := require.New(t)

			var instance *v1pb.Instance
			databaseOwner := ""
			switch tt.dbType {
			case storepb.Engine_MYSQL:
				instance = mysqlInstance
			case storepb.Engine_POSTGRES:
				instance = pgInstance
				databaseOwner = "postgres"
			default:
				a.FailNow("unsupported db type")
			}

			err = ctl.createDatabaseV2(ctx, ctl.project, instance, nil, tt.databaseName, databaseOwner)
			a.NoError(err)

			databaseResp, err := ctl.databaseServiceClient.GetDatabase(ctx, connect.NewRequest(&v1pb.GetDatabaseRequest{
				Name: fmt.Sprintf("%s/databases/%s", instance.Name, tt.databaseName),
			}))
			a.NoError(err)
			database := databaseResp.Msg

			sheetResp, err := ctl.sheetServiceClient.CreateSheet(ctx, connect.NewRequest(&v1pb.CreateSheetRequest{
				Parent: ctl.project.Name,
				Sheet: &v1pb.Sheet{
					Title:   "prepareStatements",
					Content: []byte(tt.prepareStatements),
				},
			}))
			a.NoError(err)
			sheet := sheetResp.Msg

			err = ctl.changeDatabase(ctx, ctl.project, database, sheet, v1pb.MigrationType_DDL)
			a.NoError(err)

			// Execute query
			queryResp, err := ctl.sqlServiceClient.Query(ctx, connect.NewRequest(&v1pb.QueryRequest{
				Name:         database.Name,
				Statement:    tt.query,
				DataSourceId: "admin",
			}))

			a.NoError(err)
			a.NotNil(queryResp)
			a.NotEmpty(queryResp.Msg.Results)

			// Check backward compatibility: error string should always be populated
			for _, result := range queryResp.Msg.Results {
				if result.GetSyntaxError() != nil {
					if tt.wantErrorString {
						a.NotEmpty(result.Error, "Error string must be populated for backward compatibility")
						// Verify both structured and string errors have consistent information
						a.Contains(result.Error, "syntax", "Error string should mention syntax")
					}

					// Verify structured error is also present
					a.NotNil(result.GetSyntaxError(), "Structured error should be present")
				}
			}
		})
	}
}

func TestStructuredErrorDetails_StopOnError(t *testing.T) {
	tests := []struct {
		name              string
		databaseName      string
		dbType            storepb.Engine
		prepareStatements string
		query             string
		wantResults       int // Number of successful results before syntax error
	}{
		{
			name:              "MySQL - Stop on syntax error in second statement",
			databaseName:      "TestStopOnSyntaxMySQL",
			dbType:            storepb.Engine_MYSQL,
			prepareStatements: "CREATE TABLE tbl(id INT PRIMARY KEY);",
			query:             "INSERT INTO tbl VALUES(1); SELCT * FROM tbl; INSERT INTO tbl VALUES(2);",
			wantResults:       1, // First insert succeeds, syntax error stops execution
		},
		{
			name:              "PostgreSQL - Stop on syntax error in second statement",
			databaseName:      "TestStopOnSyntaxPostgres",
			dbType:            storepb.Engine_POSTGRES,
			prepareStatements: "CREATE TABLE tbl(id INT PRIMARY KEY);",
			query:             "INSERT INTO tbl VALUES(1); SELCT * FROM tbl; INSERT INTO tbl VALUES(2);",
			wantResults:       1, // First insert succeeds, syntax error stops execution
		},
	}

	t.Parallel()
	a := require.New(t)
	ctx := context.Background()
	ctl := &controller{}
	ctx, err := ctl.StartServerWithExternalPg(ctx)
	a.NoError(err)
	defer ctl.Close(ctx)

	mysqlContainer, err := getMySQLContainer(ctx)
	a.NoError(err)
	t.Cleanup(func() {
		mysqlContainer.Close(ctx)
	})

	pgContainer, err := getPgContainer(ctx)
	a.NoError(err)
	t.Cleanup(func() {
		pgContainer.Close(ctx)
	})

	mysqlInstanceResp, err := ctl.instanceServiceClient.CreateInstance(ctx, connect.NewRequest(&v1pb.CreateInstanceRequest{
		InstanceId: generateRandomString("instance"),
		Instance: &v1pb.Instance{
			Title:       "mysqlInstance",
			Engine:      v1pb.Engine_MYSQL,
			Environment: stringPtr("environments/prod"),
			Activation:  true,
			DataSources: []*v1pb.DataSource{{Type: v1pb.DataSourceType_ADMIN, Host: mysqlContainer.host, Port: mysqlContainer.port, Username: "root", Password: "root-password", Id: "admin"}},
		},
	}))
	a.NoError(err)
	mysqlInstance := mysqlInstanceResp.Msg

	pgInstanceResp, err := ctl.instanceServiceClient.CreateInstance(ctx, connect.NewRequest(&v1pb.CreateInstanceRequest{
		InstanceId: generateRandomString("instance"),
		Instance: &v1pb.Instance{
			Title:       "pgInstance",
			Engine:      v1pb.Engine_POSTGRES,
			Environment: stringPtr("environments/prod"),
			Activation:  true,
			DataSources: []*v1pb.DataSource{{Type: v1pb.DataSourceType_ADMIN, Host: pgContainer.host, Port: pgContainer.port, Username: "postgres", Password: "root-password", Id: "admin"}},
		},
	}))
	a.NoError(err)
	pgInstance := pgInstanceResp.Msg

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := require.New(t)

			var instance *v1pb.Instance
			databaseOwner := ""
			switch tt.dbType {
			case storepb.Engine_MYSQL:
				instance = mysqlInstance
			case storepb.Engine_POSTGRES:
				instance = pgInstance
				databaseOwner = "postgres"
			default:
				a.FailNow("unsupported db type")
			}

			err = ctl.createDatabaseV2(ctx, ctl.project, instance, nil, tt.databaseName, databaseOwner)
			a.NoError(err)

			databaseResp, err := ctl.databaseServiceClient.GetDatabase(ctx, connect.NewRequest(&v1pb.GetDatabaseRequest{
				Name: fmt.Sprintf("%s/databases/%s", instance.Name, tt.databaseName),
			}))
			a.NoError(err)
			database := databaseResp.Msg

			sheetResp, err := ctl.sheetServiceClient.CreateSheet(ctx, connect.NewRequest(&v1pb.CreateSheetRequest{
				Parent: ctl.project.Name,
				Sheet: &v1pb.Sheet{
					Title:   "prepareStatements",
					Content: []byte(tt.prepareStatements),
				},
			}))
			a.NoError(err)
			sheet := sheetResp.Msg

			err = ctl.changeDatabase(ctx, ctl.project, database, sheet, v1pb.MigrationType_DDL)
			a.NoError(err)

			// Execute query that should stop on syntax error
			queryResp, err := ctl.sqlServiceClient.Query(ctx, connect.NewRequest(&v1pb.QueryRequest{
				Name:         database.Name,
				Statement:    tt.query,
				DataSourceId: "admin",
			}))

			a.NoError(err)
			a.NotNil(queryResp)

			// Count successful results before syntax error
			successfulResults := 0
			var foundSyntaxError bool
			for _, result := range queryResp.Msg.Results {
				if result.Error == "" {
					successfulResults++
				} else if result.GetSyntaxError() != nil {
					foundSyntaxError = true
					// Verify structured error details
					a.NotNil(result.GetSyntaxError().Position)
					break
				}
			}

			a.Equal(tt.wantResults, successfulResults, "Expected %d successful results before syntax error", tt.wantResults)
			a.True(foundSyntaxError, "Should have encountered a syntax error")
			// Third statement should not execute (stop-on-error behavior)
			a.LessOrEqual(len(queryResp.Msg.Results), tt.wantResults+1, "Execution should stop after syntax error")
		})
	}
}
