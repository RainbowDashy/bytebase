# Structured Error Details Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add structured error details (SyntaxErrorDetail, PermissionDeniedDetail) to QueryResult proto to enable better error handling in the frontend without string parsing.

**Architecture:** Extend the existing `detailed_error` oneof in QueryResult proto with two new minimal message types. Update backend to populate these fields when specific errors occur. Update frontend to consume structured errors instead of parsing error strings.

**Tech Stack:** Protocol Buffers, Go, TypeScript/Vue, buf generate

---

## Task 1: Update Proto Definition

**Files:**
- Modify: `proto/v1/v1/sql_service.proto:180-227` (QueryResult message)

**Step 1: Add new error detail messages to proto**

In `proto/v1/v1/sql_service.proto`, after the existing `PostgresError` message (around line 227), add:

```protobuf
  // Syntax error with position information for editor highlighting
  message SyntaxErrorDetail {
    // Position information for highlighting in editor
    Position position = 1;
  }

  // Permission denied with resource information for "Request Access" button
  message PermissionDeniedDetail {
    // Resource path: instances/{instance}/databases/{database}/schemas/{schema}/tables/{table}
    string resource = 1;
  }
```

**Step 2: Extend the detailed_error oneof**

In `proto/v1/v1/sql_service.proto`, modify the `oneof detailed_error` (around line 202-204):

```protobuf
  oneof detailed_error {
    PostgresError postgres_error = 9;
    SyntaxErrorDetail syntax_error = 13;
    PermissionDeniedDetail permission_denied = 14;
  }
```

**Step 3: Format the proto file**

Run: `buf format -w proto`
Expected: File formatted according to buf style

**Step 4: Lint the proto file**

Run: `buf lint proto`
Expected: No linting errors

**Step 5: Generate code from proto**

Run: `cd proto && buf generate`
Expected: Generated files updated in:
- `backend/generated-go/v1/sql_service.pb.go`
- `frontend/src/types/proto-es/v1/sql_service_pb.d.ts`
- `frontend/src/types/proto-es/v1/sql_service_pb.js`

**Step 6: Commit proto changes**

```bash
git add proto/v1/v1/sql_service.proto
git add backend/generated-go/v1/
git add frontend/src/types/proto-es/v1/
git commit -m "$(cat <<'EOF'
feat: add structured error details to QueryResult

- Add SyntaxErrorDetail with position for editor highlighting
- Add PermissionDeniedDetail with resource path for access requests
- Extend detailed_error oneof to include new error types

This enables frontend to handle errors without string parsing.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
EOF
)"
```

---

## Task 2: Backend - Populate SyntaxErrorDetail

**Files:**
- Modify: `backend/api/v1/sql_service.go:235-258` (Query method error handling)

**Step 1: Update syntax error handling to populate detailed_error**

In `backend/api/v1/sql_service.go`, replace the existing syntax error handling (lines 240-256):

```go
} else if syntaxErr, ok := queryErr.(*parserbase.SyntaxError); ok {
	err := connect.NewError(connect.CodeInvalidArgument, syntaxErr)
	if detail, detailErr := connect.NewErrorDetail(&v1pb.PlanCheckRun_Result{
		Code:    int32(advisor.StatementSyntaxError),
		Content: syntaxErr.Message,
		Title:   "Syntax error",
		Status:  v1pb.Advice_ERROR,
		Report: &v1pb.PlanCheckRun_Result_SqlReviewReport_{
			SqlReviewReport: &v1pb.PlanCheckRun_Result_SqlReviewReport{
				StartPosition: convertToPosition(syntaxErr.Position),
			},
		},
	}); detailErr == nil {
		err.AddDetail(detail)
	}
	return nil, err
}
```

With:

```go
} else if syntaxErr, ok := queryErr.(*parserbase.SyntaxError); ok {
	// Return error with structured details in QueryResult
	// Note: We return error in results array, not as connect.Error,
	// to support stop-on-error behavior where we need partial results
	results = []*v1pb.QueryResult{
		{
			Error:     syntaxErr.Message,
			Statement: statement,
			DetailedError: &v1pb.QueryResult_SyntaxError{
				SyntaxError: &v1pb.QueryResult_SyntaxErrorDetail{
					Position: convertToPosition(syntaxErr.Position),
				},
			},
		},
	}
	return connect.NewResponse(&v1pb.QueryResponse{Results: results}), nil
}
```

**Step 2: Format the modified file**

Run: `gofmt -w backend/api/v1/sql_service.go`
Expected: File formatted

**Step 3: Verify build**

Run: `go build -ldflags "-w -s" -p=16 -o ./bytebase-build/bytebase ./backend/bin/server/main.go`
Expected: Build succeeds without errors

**Step 4: Commit**

```bash
git add backend/api/v1/sql_service.go
git commit -m "$(cat <<'EOF'
feat: populate SyntaxErrorDetail in QueryResult

Return syntax errors in QueryResult with structured position info
instead of connect.Error to support stop-on-error with partial results.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: Backend - Populate PermissionDeniedDetail

**Files:**
- Modify: `backend/api/v1/sql_service.go:1427-1430` (accessCheck function)

**Step 1: Find where permission denied errors are created**

Search for the permission denied error creation in accessCheck:

Run: `grep -n "permission denied to access resource" backend/api/v1/sql_service.go`
Expected: Find the line where the error message is constructed

**Step 2: Create helper function to populate PermissionDeniedDetail**

Add this function after the `queryRetryStopOnError` function (around line 728):

```go
func createPermissionDeniedResult(resource string, statement string) *v1pb.QueryResult {
	return &v1pb.QueryResult{
		Error:     fmt.Sprintf("permission denied to access resource: %s", resource),
		Statement: statement,
		DetailedError: &v1pb.QueryResult_PermissionDenied{
			PermissionDenied: &v1pb.QueryResult_PermissionDeniedDetail{
				Resource: resource,
			},
		},
	}
}
```

**Step 3: Update accessCheck to return structured error**

In `backend/api/v1/sql_service.go`, find the permission denied error (around line 1427-1430) and note that it returns a connect.Error. We need to handle this in queryRetry instead.

**Step 4: Update queryRetry to detect and structure permission denied errors**

In `backend/api/v1/sql_service.go`, in the `queryRetry` function around line 464, after the accessCheck call:

```go
if optionalAccessCheck != nil {
	// Check query access
	if err := optionalAccessCheck(ctx, instance, database, user, spans, queryContext.Explain); err != nil {
		// Check if it's a permission denied error and extract resource
		if connectErr, ok := err.(*connect.Error); ok && connectErr.Code() == connect.CodePermissionDenied {
			errMsg := connectErr.Message()
			if strings.Contains(errMsg, "permission denied to access resource:") {
				resource := strings.TrimPrefix(errMsg, "permission denied to access resource: ")
				// Return as QueryResult with structured error instead of error return
				return []*v1pb.QueryResult{createPermissionDeniedResult(resource, statement)}, nil, time.Duration(0), nil
			}
		}
		return nil, nil, time.Duration(0), err
	}
	slog.Debug("optional access check", slog.String("instance", instance.ResourceID), slog.String("database", database.DatabaseName))
}
```

**Step 5: Format the modified file**

Run: `gofmt -w backend/api/v1/sql_service.go`
Expected: File formatted

**Step 6: Verify build**

Run: `go build -ldflags "-w -s" -p=16 -o ./bytebase-build/bytebase ./backend/bin/server/main.go`
Expected: Build succeeds without errors

**Step 7: Commit**

```bash
git add backend/api/v1/sql_service.go
git commit -m "$(cat <<'EOF'
feat: populate PermissionDeniedDetail in QueryResult

Extract resource from permission denied errors and populate
structured PermissionDeniedDetail for frontend consumption.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: Frontend - Update Error Handling for PermissionDeniedDetail

**Files:**
- Modify: `frontend/src/views/sql-editor/EditorCommon/ResultView/ResultViewV1.vue:256-269`

**Step 1: Update missingResource computed to use structured error**

In `frontend/src/views/sql-editor/EditorCommon/ResultView/ResultViewV1.vue`, replace the `missingResource` computed (lines 256-269):

```typescript
const missingResource = computed((): DatabaseResource | undefined => {
  if (props.resultSet?.status !== Code.PermissionDenied) {
    return;
  }

  // Check for structured error first
  const firstResult = props.resultSet.results[0];
  if (firstResult?.detailedError?.case === 'permissionDenied') {
    const resource = firstResult.detailedError.value.resource;
    if (resource) {
      return parseStringToResource(resource);
    }
  }

  // Fallback to parsing error string for backward compatibility
  const prefix = "permission denied to access resource: ";
  if (!props.resultSet.error.includes(prefix)) {
    return;
  }
  const resource = props.resultSet.error.split(prefix).pop();
  if (!resource) {
    return;
  }
  return parseStringToResource(resource);
});
```

**Step 2: Format the file**

Run: `cd frontend && pnpm prettier --write src/views/sql-editor/EditorCommon/ResultView/ResultViewV1.vue`
Expected: File formatted

**Step 3: Type check**

Run: `pnpm --dir frontend type-check`
Expected: No type errors

**Step 4: Commit**

```bash
git add frontend/src/views/sql-editor/EditorCommon/ResultView/ResultViewV1.vue
git commit -m "$(cat <<'EOF'
feat: use structured PermissionDeniedDetail in ResultView

Check detailedError for structured permission info before
falling back to string parsing for backward compatibility.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: Frontend - Add SyntaxError Highlighting Support

**Files:**
- Create: `frontend/src/views/sql-editor/EditorCommon/ResultView/ErrorView/SyntaxError.vue`
- Modify: `frontend/src/views/sql-editor/EditorCommon/ResultView/ErrorView/ErrorView.vue:11`

**Step 1: Create SyntaxError component**

Create `frontend/src/views/sql-editor/EditorCommon/ResultView/ErrorView/SyntaxError.vue`:

```vue
<template>
  <div
    v-for="(error, i) in syntaxErrors"
    :key="i"
    class="text-sm grid gap-1 pl-8"
    style="grid-template-columns: auto 1fr"
  >
    <template v-if="error.position">
      <div>POSITION:</div>
      <div>Line {{ error.position.line }}, Column {{ error.position.column }}</div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import type { SQLResultSetV1 } from "@/types";
import type { QueryResult_SyntaxErrorDetail } from "@/types/proto-es/v1/sql_service_pb";

const props = defineProps<{
  resultSet: SQLResultSetV1;
}>();

const syntaxErrors = computed(() => {
  const errors: QueryResult_SyntaxErrorDetail[] = [];
  props.resultSet.results.forEach((result) => {
    if (result.detailedError?.case === "syntaxError") {
      errors.push(result.detailedError.value);
    }
  });
  return errors;
});
</script>
```

**Step 2: Import and use SyntaxError component in ErrorView**

In `frontend/src/views/sql-editor/EditorCommon/ResultView/ErrorView/ErrorView.vue`, add import:

```typescript
import SyntaxError from "./SyntaxError.vue";
```

And add component after PostgresError (line 11):

```vue
<PostgresError v-if="resultSet" :result-set="resultSet" />
<SyntaxError v-if="resultSet" :result-set="resultSet" />
```

**Step 3: Format the files**

Run: `cd frontend && pnpm prettier --write "src/views/sql-editor/EditorCommon/ResultView/ErrorView/*.vue"`
Expected: Files formatted

**Step 4: Lint**

Run: `pnpm --dir frontend lint --fix`
Expected: No lint errors

**Step 5: Type check**

Run: `pnpm --dir frontend type-check`
Expected: No type errors

**Step 6: Commit**

```bash
git add frontend/src/views/sql-editor/EditorCommon/ResultView/ErrorView/
git commit -m "$(cat <<'EOF'
feat: add SyntaxError component for structured error display

Display syntax error position information from SyntaxErrorDetail
in a consistent format with PostgresError.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: Testing

**Files:**
- Modify: `backend/tests/sql_stop_on_error_test.go` (add error detail tests)

**Step 1: Add test for SyntaxErrorDetail**

Add test case to `backend/tests/sql_stop_on_error_test.go`:

```go
{
	name:              "MySQL - Syntax error with structured detail",
	databaseName:      "TestSyntaxError",
	dbType:            storepb.Engine_MYSQL,
	prepareStatements: "CREATE TABLE tbl(id INT PRIMARY KEY);",
	query:             "SELCT * FROM tbl;", // Intentional typo
	wantResults:       0,
	wantError:         true,
},
```

**Step 2: Add test for PermissionDeniedDetail**

Add test case to `backend/tests/sql_stop_on_error_test.go`:

```go
{
	name:              "MySQL - Permission denied with structured detail",
	databaseName:      "TestPermissionDenied",
	dbType:            storepb.Engine_MYSQL,
	prepareStatements: "CREATE TABLE tbl(id INT PRIMARY KEY);",
	query:             "SELECT * FROM tbl;", // Will need permission setup
	wantResults:       0,
	wantError:         true,
},
```

**Step 3: Format test file**

Run: `gofmt -w backend/tests/sql_stop_on_error_test.go`
Expected: File formatted

**Step 4: Run tests (if Docker available)**

Run: `cd backend/tests && go test -v -count=1 -run TestSQLQueryStopOnError`
Expected: Tests pass (or skip if no Docker)

**Step 5: Commit**

```bash
git add backend/tests/sql_stop_on_error_test.go
git commit -m "$(cat <<'EOF'
test: add structured error detail test cases

Add test cases for SyntaxErrorDetail and PermissionDeniedDetail
in stop-on-error scenarios.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
EOF
)"
```

---

## Task 7: Documentation

**Files:**
- Create: `docs/design/structured-error-details.md`

**Step 1: Write design documentation**

Create `docs/design/structured-error-details.md`:

```markdown
# Structured Error Details

## Overview

QueryResult now supports structured error details via the `detailed_error` oneof field, enabling frontends to handle errors without string parsing.

## Error Types

### SyntaxErrorDetail
- **When**: SQL syntax errors
- **Contains**: Position (line, column) for editor highlighting
- **Frontend use**: Highlight error location in SQL editor

### PermissionDeniedDetail
- **When**: User lacks permission to access resource
- **Contains**: Resource path (instances/{}/databases/{}/tables/{})
- **Frontend use**: Show "Request Access" button with specific resource

### PostgresError (existing)
- **When**: PostgreSQL-specific errors
- **Contains**: Detail, hint, where fields
- **Frontend use**: Display detailed Postgres error information

## Usage

### Backend (Populating)

```go
// Syntax error
result := &v1pb.QueryResult{
	Error: "syntax error near 'SELCT'",
	DetailedError: &v1pb.QueryResult_SyntaxError{
		SyntaxError: &v1pb.QueryResult_SyntaxErrorDetail{
			Position: &v1pb.Position{Line: 1, Column: 1},
		},
	},
}

// Permission denied
result := &v1pb.QueryResult{
	Error: "permission denied to access resource: instances/prod/databases/users/tables/sensitive",
	DetailedError: &v1pb.QueryResult_PermissionDenied{
		PermissionDenied: &v1pb.QueryResult_PermissionDeniedDetail{
			Resource: "instances/prod/databases/users/tables/sensitive",
		},
	},
}
```

### Frontend (Consuming)

```typescript
// Check error type
if (result.detailedError?.case === 'syntaxError') {
  const position = result.detailedError.value.position;
  highlightError(position);
}

if (result.detailedError?.case === 'permissionDenied') {
  const resource = result.detailedError.value.resource;
  showRequestAccessButton(resource);
}
```

## Design Principles

1. **Minimal**: Only structured data that can't be in a string
2. **Backward compatible**: `error` string field still populated
3. **Type-safe**: Frontend uses discriminated union (oneof)
4. **Extensible**: Easy to add new error types in the future
```

**Step 2: Commit documentation**

```bash
git add docs/design/structured-error-details.md
git commit -m "$(cat <<'EOF'
docs: add structured error details design doc

Document the new SyntaxErrorDetail and PermissionDeniedDetail
error types and their usage patterns.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
EOF
)"
```

---

## Verification Checklist

After completing all tasks, verify:

- [ ] Proto builds without errors (`buf lint proto` passes)
- [ ] Backend builds without errors (`go build` succeeds)
- [ ] Frontend type-checks without errors (`pnpm type-check` passes)
- [ ] Frontend lints without errors (`pnpm lint` passes)
- [ ] All commits follow conventional commit format
- [ ] Documentation is complete and accurate

## Notes

- The `error` string field is always populated for backward compatibility
- Frontend checks `detailedError` first, falls back to string parsing
- Stop-on-error behavior means errors can be in individual QueryResults, not just response-level
- New error types can be added by extending the `detailed_error` oneof
