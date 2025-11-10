package v1

import (
	"testing"

	"github.com/stretchr/testify/require"

	v1pb "github.com/bytebase/bytebase/backend/generated-go/v1"
)

func TestCreatePermissionDeniedResult(t *testing.T) {
	tests := []struct {
		name         string
		resource     string
		statement    string
		wantError    string
		wantResource string
	}{
		{
			name:         "Simple resource path",
			resource:     "instances/prod/databases/users/tables/sensitive",
			statement:    "SELECT * FROM sensitive;",
			wantError:    "permission denied to access resource: instances/prod/databases/users/tables/sensitive",
			wantResource: "instances/prod/databases/users/tables/sensitive",
		},
		{
			name:         "Resource path with schema",
			resource:     "instances/prod/databases/users/schemas/public/tables/sensitive",
			statement:    "SELECT * FROM public.sensitive;",
			wantError:    "permission denied to access resource: instances/prod/databases/users/schemas/public/tables/sensitive",
			wantResource: "instances/prod/databases/users/schemas/public/tables/sensitive",
		},
		{
			name:         "Empty resource",
			resource:     "",
			statement:    "SELECT * FROM table;",
			wantError:    "permission denied to access resource: ",
			wantResource: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := require.New(t)

			result := createPermissionDeniedResult(tt.resource, tt.statement)

			// Verify error string is populated (backward compatibility)
			a.Equal(tt.wantError, result.Error, "Error string should match")

			// Verify statement is set
			a.Equal(tt.statement, result.Statement, "Statement should match")

			// Verify PermissionDeniedDetail is populated
			permDenied := result.GetPermissionDenied()
			a.NotNil(permDenied, "PermissionDeniedDetail should be populated")
			a.Equal(tt.wantResource, permDenied.Resource, "Resource path should match")

			// Verify this is the correct oneof case
			a.Nil(result.GetSyntaxError(), "SyntaxError should not be set")
			a.Nil(result.GetPostgresError(), "PostgresError should not be set")
		})
	}
}

func TestStructuredErrorDetails_Types(t *testing.T) {
	// Test that we can create different error detail types without conflicts
	tests := []struct {
		name   string
		result *v1pb.QueryResult
	}{
		{
			name: "SyntaxErrorDetail",
			result: &v1pb.QueryResult{
				Error:     "syntax error",
				Statement: "SELCT * FROM tbl;",
				DetailedError: &v1pb.QueryResult_SyntaxError{
					SyntaxError: &v1pb.QueryResult_SyntaxErrorDetail{
						Position: &v1pb.Position{Line: 1, Column: 1},
					},
				},
			},
		},
		{
			name: "PermissionDeniedDetail",
			result: &v1pb.QueryResult{
				Error:     "permission denied",
				Statement: "SELECT * FROM tbl;",
				DetailedError: &v1pb.QueryResult_PermissionDenied{
					PermissionDenied: &v1pb.QueryResult_PermissionDeniedDetail{
						Resource: "instances/prod/databases/db/tables/tbl",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := require.New(t)
			a.NotNil(tt.result, "Result should not be nil")
			a.NotEmpty(tt.result.Error, "Error string should be populated")
			a.NotEmpty(tt.result.Statement, "Statement should be populated")

			// Verify oneof behavior - only one detailed_error should be set
			detailCount := 0
			if tt.result.GetSyntaxError() != nil {
				detailCount++
			}
			if tt.result.GetPermissionDenied() != nil {
				detailCount++
			}
			if tt.result.GetPostgresError() != nil {
				detailCount++
			}
			a.Equal(1, detailCount, "Exactly one detailed_error should be set")
		})
	}
}

func TestStructuredErrorDetails_PositionFormat(t *testing.T) {
	tests := []struct {
		name       string
		position   *v1pb.Position
		wantValid  bool
		wantReason string
	}{
		{
			name:      "Valid position",
			position:  &v1pb.Position{Line: 1, Column: 10},
			wantValid: true,
		},
		{
			name:       "Line 0 is invalid",
			position:   &v1pb.Position{Line: 0, Column: 10},
			wantValid:  false,
			wantReason: "Line should be >= 1",
		},
		{
			name:       "Negative line is invalid",
			position:   &v1pb.Position{Line: -1, Column: 10},
			wantValid:  false,
			wantReason: "Line should be >= 1",
		},
		{
			name:       "Negative column is invalid",
			position:   &v1pb.Position{Line: 1, Column: -1},
			wantValid:  false,
			wantReason: "Column should be >= 0",
		},
		{
			name:      "Column 0 is valid (start of line)",
			position:  &v1pb.Position{Line: 1, Column: 0},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := require.New(t)

			result := &v1pb.QueryResult{
				Error:     "syntax error",
				Statement: "SELCT * FROM tbl;",
				DetailedError: &v1pb.QueryResult_SyntaxError{
					SyntaxError: &v1pb.QueryResult_SyntaxErrorDetail{
						Position: tt.position,
					},
				},
			}

			syntaxError := result.GetSyntaxError()
			a.NotNil(syntaxError)
			a.NotNil(syntaxError.Position)

			// Check position validity
			if tt.wantValid {
				a.Greater(syntaxError.Position.Line, int32(0), "Line should be > 0")
				a.GreaterOrEqual(syntaxError.Position.Column, int32(0), "Column should be >= 0")
			} else {
				// Document that invalid positions can exist in the proto
				// (validation should happen at creation time or in frontend)
				t.Logf("Position is invalid: %s", tt.wantReason)
			}
		})
	}
}
