# Structured Error Details Design Document

## Overview

This document describes the structured error details feature added to the QueryResult proto message. This feature enables the frontend to handle specific error types (syntax errors and permission denied errors) without parsing error strings, improving type safety, maintainability, and enabling better user experiences.

## Motivation

### Problems with String-Based Error Handling

Prior to this feature, the frontend had to parse error strings to extract meaningful information:

1. **Fragile**: String parsing breaks when error message format changes
2. **Type-unsafe**: No compile-time guarantees about error structure
3. **Incomplete**: Position information for syntax errors required parsing PlanCheckRun_Result details
4. **Localization issues**: Error string parsing doesn't work with internationalized messages

### Use Cases Enabled

1. **Syntax Error Highlighting**: Display error position (line/column) in SQL editor for immediate visual feedback
2. **Request Access Button**: Show "Request Access" button with specific resource path when user lacks permissions
3. **Better Error UX**: Provide structured, actionable error information to users

## Design Decisions

### 1. Extend Existing `detailed_error` Oneof

We extended the existing `detailed_error` oneof field in QueryResult rather than creating a new field:

```protobuf
oneof detailed_error {
  PostgresError postgres_error = 9;           // Existing
  SyntaxErrorDetail syntax_error = 13;        // New
  PermissionDeniedDetail permission_denied = 14; // New
}
```

**Rationale**:
- Consistent with existing PostgresError pattern
- Type-safe discriminated union (only one error type at a time)
- Clean API surface without field proliferation
- Easy to extend with new error types in the future

### 2. Minimal Message Definitions

We keep the error detail messages minimal, containing only structured data that cannot be represented as strings:

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

**Rationale**:
- Position (line/column) is the only syntax error data needed for highlighting
- Resource path is the only permission data needed for access requests
- Error message stays in the `error` string field for backward compatibility
- Minimal fields reduce maintenance burden and API surface

### 3. Backward Compatibility Guarantees

The `error` string field is **always populated** alongside structured details:

**Backend**:
```go
result := &v1pb.QueryResult{
    Error:     syntaxErr.Message,  // Always set
    Statement: statement,
    DetailedError: &v1pb.QueryResult_SyntaxError{
        SyntaxError: &v1pb.QueryResult_SyntaxErrorDetail{
            Position: convertToPosition(syntaxErr.Position),
        },
    },
}
```

**Frontend**:
```typescript
// Check for structured error first
if (firstResult?.detailedError?.case === 'permissionDenied') {
  const resource = firstResult.detailedError.value.resource;
  return parseStringToResource(resource);
}
// Fallback to parsing error string
const resource = props.resultSet.error.split(prefix).pop();
```

**Rationale**:
- Old clients continue to work without changes
- New clients prefer structured data but gracefully fall back
- Enables gradual migration and safe rollbacks
- Tests verify both fields are populated

### 4. Error in QueryResult, Not Response-Level

Errors are returned in individual QueryResult objects within the results array, not as connect.Error at the response level:

```go
// Return error in results array for stop-on-error behavior
results = []*v1pb.QueryResult{
    {
        Error:     syntaxErr.Message,
        Statement: statement,
        DetailedError: &v1pb.QueryResult_SyntaxError{...},
    },
}
return connect.NewResponse(&v1pb.QueryResponse{Results: results}), nil
```

**Rationale**:
- Supports stop-on-error behavior: return partial results before error
- Consistent with multi-statement execution model
- Each statement's result can have its own error
- RPC call succeeds (HTTP 200) even when queries fail

## API Changes

### Proto Definition

File: `proto/v1/v1/sql_service.proto`

```protobuf
message QueryResult {
  // ... existing fields ...

  // The error message if the query failed.
  string error = 6;

  oneof detailed_error {
    PostgresError postgres_error = 9;
    SyntaxErrorDetail syntax_error = 13;
    PermissionDeniedDetail permission_denied = 14;
  }

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

  // ... existing PostgresError message ...
}
```

### Generated Code

After running `cd proto && buf generate`, the following files are updated:

1. **Go**: `backend/generated-go/v1/sql_service.pb.go`
   - `QueryResult_SyntaxError` struct wrapping `SyntaxErrorDetail`
   - `QueryResult_PermissionDenied` struct wrapping `PermissionDeniedDetail`
   - `GetSyntaxError()` and `GetPermissionDenied()` accessor methods

2. **TypeScript**: `frontend/src/types/proto-es/v1/sql_service_pb.ts`
   - Type-safe discriminated union for `detailedError`
   - Case checking via `detailedError?.case === 'syntaxError'`

## Backend Implementation

### 1. Syntax Error Handling

**Location**: `backend/api/v1/sql_service.go:456-467, 555-566`

Syntax errors are detected during query parsing and converted to structured QueryResult:

```go
spans, err := parseAndCheckWithObjectCatalog(
    ctx, stores, licenseService, instance, database,
    statement, queryContext.Schema, !store.IsObjectCaseSensitive(instance),
)
if err != nil {
    // Handle syntax errors with structured details
    if syntaxErr, ok := err.(*parserbase.SyntaxError); ok {
        result := &v1pb.QueryResult{
            Error:     syntaxErr.Message,
            Statement: statement,
            DetailedError: &v1pb.QueryResult_SyntaxError{
                SyntaxError: &v1pb.QueryResult_SyntaxErrorDetail{
                    Position: convertToPosition(syntaxErr.Position),
                },
            },
        }
        return []*v1pb.QueryResult{result}, nil, time.Duration(0), nil
    }
    return nil, nil, time.Duration(0), err
}
```

**Implementation in Two Places**:
1. `queryRetry` function (line 456): First-statement syntax errors
2. `queryRetryStopOnError` function (line 555): Multi-statement syntax errors

**convertToPosition Helper**:
Converts parser position to proto Position with line and column numbers (1-indexed for user display).

### 2. Permission Denied Handling

**Location**: `backend/api/v1/sql_service.go:477-486, 766-776`

Permission errors are detected during access checks and converted to structured QueryResult:

```go
if optionalAccessCheck != nil {
    if err := optionalAccessCheck(ctx, instance, database, user, spans, queryContext.Explain); err != nil {
        // Check if it's a permission denied error and extract resource
        if connectErr, ok := err.(*connect.Error); ok && connectErr.Code() == connect.CodePermissionDenied {
            errMsg := connectErr.Message()
            if strings.Contains(errMsg, "permission denied to access resource:") {
                resource := strings.TrimSpace(strings.TrimPrefix(errMsg, "permission denied to access resource:"))
                // Return as QueryResult with structured error
                return []*v1pb.QueryResult{createPermissionDeniedResult(resource, statement)}, nil, time.Duration(0), nil
            }
        }
        return nil, nil, time.Duration(0), err
    }
}
```

**Helper Function** (line 766):
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

**Resource Path Format**:
```
instances/{instance}/databases/{database}/schemas/{schema}/tables/{table}
```

## Frontend Implementation

### 1. Permission Denied Handling

**Location**: `frontend/src/views/sql-editor/EditorCommon/ResultView/ResultViewV1.vue:256-280`

The frontend checks for structured error details first, then falls back to string parsing:

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

**Key Points**:
- Checks `detailedError?.case === 'permissionDenied'` for type-safe access
- Extracts resource path from structured data
- Falls back to string parsing if structured data unavailable
- Used to show "Request Access" button with specific resource

### 2. Syntax Error Display

**Location**: `frontend/src/views/sql-editor/EditorCommon/ResultView/ErrorView/SyntaxError.vue`

New component displays syntax error position information:

```vue
<template>
  <div
    v-for="(error, i) in syntaxErrors"
    :key="i"
    class="text-sm grid gap-1 pl-8"
    style="grid-template-columns: auto 1fr"
  >
    <template v-if="error.position">
      <div>{{ $t("common.position") }}:</div>
      <div>
        {{
          $t("sql-editor.syntax-error-position", {
            line: error.position.line,
            column: error.position.column,
          })
        }}
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
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

**Integration**: Added to `ErrorView.vue` alongside existing PostgresError component:

```vue
<PostgresError v-if="resultSet" :result-set="resultSet" />
<SyntaxError v-if="resultSet" :result-set="resultSet" />
```

**Key Points**:
- Type-safe access via discriminated union
- Displays line and column numbers with i18n support
- Consistent styling with PostgresError component
- Prepares for future editor highlighting integration

## Testing Approach

### Test Structure

File: `backend/tests/structured_error_details_test.go`

Three comprehensive test suites verify the implementation:

#### 1. TestStructuredErrorDetails_SyntaxError

Tests that syntax errors populate SyntaxErrorDetail with position information:

```go
{
    name:              "MySQL - Syntax error populates SyntaxErrorDetail",
    databaseName:      "TestMySQLSyntaxError",
    dbType:            storepb.Engine_MYSQL,
    prepareStatements: "CREATE TABLE tbl(id INT PRIMARY KEY, name VARCHAR(64));",
    query:             "SELCT * FROM tbl;", // Intentional typo
    wantSyntaxError:   true,
    wantPosition:      true,
},
```

**Verifications**:
- SyntaxErrorDetail is populated
- Position information (line/column) is present and valid
- Error string is populated for backward compatibility
- Statement is preserved in result
- Tests both MySQL and PostgreSQL

#### 2. TestStructuredErrorDetails_BackwardCompatibility

Tests that error strings are always populated alongside structured details:

```go
// Check backward compatibility: error string should always be populated
for _, result := range queryResp.Msg.Results {
    if result.GetSyntaxError() != nil {
        a.NotEmpty(result.Error, "Error string must be populated for backward compatibility")
        a.Contains(result.Error, "syntax", "Error string should mention syntax")
        a.NotNil(result.GetSyntaxError(), "Structured error should be present")
    }
}
```

**Key Points**:
- Verifies both structured and string errors coexist
- Ensures error strings contain meaningful content
- Prevents regression where only structured details are populated

#### 3. TestStructuredErrorDetails_StopOnError

Tests that execution stops on syntax error in multi-statement queries:

```go
{
    name:              "MySQL - Stop on syntax error in second statement",
    databaseName:      "TestStopOnSyntaxMySQL",
    dbType:            storepb.Engine_MYSQL,
    prepareStatements: "CREATE TABLE tbl(id INT PRIMARY KEY);",
    query:             "INSERT INTO tbl VALUES(1); SELCT * FROM tbl; INSERT INTO tbl VALUES(2);",
    wantResults:       1, // First insert succeeds, syntax error stops execution
},
```

**Verifications**:
- First statement executes successfully
- Syntax error in second statement stops execution
- Third statement does not execute
- SyntaxErrorDetail is present with position
- Result array length is correct (successful + error)

### Test Coverage

- **Databases**: MySQL and PostgreSQL
- **Scenarios**: Single statement errors, multi-statement errors, semantic vs syntax errors
- **Compatibility**: Backward compatibility with string-only clients
- **Behavior**: Stop-on-error execution model

## Backward Compatibility

### For Existing Clients

1. **Error String Always Populated**: Old clients relying on `result.error` continue to work
2. **Ignore Unknown Fields**: Proto3 allows clients to ignore `detailed_error` field
3. **No Breaking Changes**: No existing fields removed or changed

### For New Clients

1. **Prefer Structured Data**: Check `detailedError` first for type-safe access
2. **Graceful Fallback**: Fall back to string parsing if structured data unavailable
3. **Progressive Enhancement**: Better UX when structured data available

### Migration Strategy

**Phase 1: Backend Deployment**
- Deploy backend with structured error population
- Both error formats available immediately
- No client changes required

**Phase 2: Frontend Migration**
- Update frontend to prefer structured errors
- Keep fallback to string parsing
- Gradual rollout with feature flags if desired

**Phase 3: Deprecation (Future)**
- Monitor usage of string-based error handling
- Eventually deprecate string parsing in favor of structured details
- Long deprecation period (6+ months) before removal

## Usage Examples

### Backend: Populating Syntax Error

```go
// When syntax error is detected during parsing
if syntaxErr, ok := err.(*parserbase.SyntaxError); ok {
    result := &v1pb.QueryResult{
        Error:     syntaxErr.Message,  // "syntax error near 'SELCT'"
        Statement: statement,           // "SELCT * FROM users;"
        DetailedError: &v1pb.QueryResult_SyntaxError{
            SyntaxError: &v1pb.QueryResult_SyntaxErrorDetail{
                Position: &v1pb.Position{Line: 1, Column: 1},
            },
        },
    }
    return []*v1pb.QueryResult{result}, nil, time.Duration(0), nil
}
```

### Backend: Populating Permission Denied

```go
// When access check fails
func createPermissionDeniedResult(resource string, statement string) *v1pb.QueryResult {
    return &v1pb.QueryResult{
        Error:     fmt.Sprintf("permission denied to access resource: %s", resource),
        Statement: statement,
        DetailedError: &v1pb.QueryResult_PermissionDenied{
            PermissionDenied: &v1pb.QueryResult_PermissionDeniedDetail{
                Resource: "instances/prod/databases/users/tables/sensitive_data",
            },
        },
    }
}
```

### Frontend: Consuming Syntax Error

```typescript
// Check for syntax error with position
function getSyntaxErrorPosition(result: QueryResult): Position | undefined {
  if (result.detailedError?.case === 'syntaxError') {
    return result.detailedError.value.position;
  }
  return undefined;
}

// Use position to highlight error in editor
const position = getSyntaxErrorPosition(result);
if (position) {
  highlightErrorInEditor(position.line, position.column);
}
```

### Frontend: Consuming Permission Denied

```typescript
// Extract resource for "Request Access" button
function getRestrictedResource(result: QueryResult): string | undefined {
  if (result.detailedError?.case === 'permissionDenied') {
    return result.detailedError.value.resource;
  }

  // Fallback to string parsing for backward compatibility
  const prefix = "permission denied to access resource: ";
  if (result.error?.includes(prefix)) {
    return result.error.split(prefix).pop();
  }

  return undefined;
}

// Show request access button
const resource = getRestrictedResource(result);
if (resource) {
  showRequestAccessButton(resource);
}
```

## Limitations and Edge Cases

### Current Limitations

1. **No Multi-Position Errors**: SyntaxErrorDetail has single position, can't represent errors spanning multiple locations
2. **No Suggestion Field**: Syntax errors don't include fix suggestions (could be added in future)
3. **Resource Format Rigidity**: Permission denied resource path format is fixed, can't represent arbitrary permission structures
4. **No Error Codes**: Syntax errors don't have numeric error codes for programmatic handling

### Edge Cases Handled

1. **Semantic vs Syntax Errors**: Only syntax errors populate SyntaxErrorDetail; semantic errors (invalid column) use generic error
2. **Multiple Statements**: Each statement's result can have its own error detail
3. **Missing Position**: Position may be nil if parser doesn't provide it (gracefully handled)
4. **Empty Resource**: Permission denied may have empty resource string (fallback to error string)

### Future Considerations

1. **Additional Error Types**: Easy to add new error detail types to `detailed_error` oneof
2. **Richer Error Data**: Could add fields to existing messages (e.g., suggestions, related errors)
3. **Error Codes**: Could add error code enums for programmatic handling
4. **Localization**: Error messages in `error` string should use i18n (structured details are locale-neutral)

## Future Enhancements

### Short Term (Next Iteration)

1. **Editor Highlighting**: Use position information to highlight syntax errors in Monaco editor
2. **Quick Fixes**: Add suggestion field to SyntaxErrorDetail for "Did you mean?" fixes
3. **More Error Types**: Add ValidationErrorDetail, TimeoutErrorDetail, etc.

### Medium Term

1. **Error Codes**: Add numeric error codes for programmatic error handling
2. **Rich Formatting**: Support markdown in error messages for better formatting
3. **Error Context**: Add surrounding code snippet to show error in context
4. **Batch Permissions**: Support multiple resources in PermissionDeniedDetail

### Long Term

1. **AI-Powered Suggestions**: Use LLM to generate syntax error fix suggestions
2. **Error Analytics**: Track common errors to improve developer experience
3. **Custom Error Types**: Allow plugins to define custom error detail types
4. **Error Recovery**: Suggest recoverable errors vs fatal errors

## References

### Related Files

- **Proto Definition**: `proto/v1/v1/sql_service.proto:180-241`
- **Backend Implementation**: `backend/api/v1/sql_service.go:456-467, 555-566, 766-776`
- **Frontend Permission Handling**: `frontend/src/views/sql-editor/EditorCommon/ResultView/ResultViewV1.vue:256-280`
- **Frontend Syntax Display**: `frontend/src/views/sql-editor/EditorCommon/ResultView/ErrorView/SyntaxError.vue`
- **Tests**: `backend/tests/structured_error_details_test.go`

### Related Commits

- `553a4ea28d`: feat: add structured error details to QueryResult (proto changes)
- `e4e8a73bcc`: feat: populate SyntaxErrorDetail for syntax errors (backend)
- `5eb4c8c845`: feat: use structured PermissionDeniedDetail in ResultView (frontend permission)
- `bd491bfb57`: feat: add syntax error highlighting with position (frontend syntax)
- `fa5a378866`: test: add tests for structured error details (tests)

### Design Principles Applied

1. **Minimal API Surface**: Only essential structured data, keep messages small
2. **Backward Compatible**: Always populate error strings alongside structured details
3. **Type Safe**: Use discriminated unions (oneof) for compile-time safety
4. **Extensible**: Easy to add new error types without breaking changes
5. **Fail Safe**: Graceful fallbacks when structured data unavailable
6. **User Focused**: Enable better UX (highlighting, request access) over perfection

## Questions and Answers

**Q: Why not put all error details in a single message?**
A: Different error types need different fields. A discriminated union (oneof) is more type-safe and prevents field bloat.

**Q: Why keep the error string if we have structured details?**
A: Backward compatibility with existing clients and as a human-readable fallback.

**Q: Can we add more fields to SyntaxErrorDetail later?**
A: Yes! Proto3 allows adding optional fields without breaking existing clients.

**Q: Why not use gRPC error details instead of QueryResult fields?**
A: We need errors in individual QueryResults for stop-on-error behavior, not at the RPC response level.

**Q: What about other database engines besides MySQL and PostgreSQL?**
A: The design is engine-agnostic. Any engine that provides position information can populate SyntaxErrorDetail.

**Q: How do we test this in development?**
A: Use intentional typos like "SELCT" in SQL queries to trigger syntax errors, or query restricted tables for permission errors.
