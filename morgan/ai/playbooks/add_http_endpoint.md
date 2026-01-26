# Playbook: Add HTTP Endpoint

## Purpose
Expose a usecase via REST HTTP API (Fiber).

## Triggers
- New API endpoint required.
- `build_feature_from_openapi` execution.

## Inputs
- **Method**: GET, POST, PUT, DELETE, etc.
- **Path**: URL path (e.g., `/v1/users`).
- **Request/Response Schema**: JSON structure.
- **Target Usecase**: The business logic to invoke.

## Preconditions
- Usecase exists (or `add_usecase` is scheduled).
- Module `delivery/http` folder exists.

## Rules to Load
- [ai/rules/00_principles.md](../rules/00_principles.md)
- [ai/rules/30_handler_http.md](../rules/30_handler_http.md)
- [ai/rules/50_helpers.md](../rules/50_helpers.md)

## Templates to Use
- [ai/templates/http_handler.md](../templates/http_handler.md)
- [ai/templates/VALIDATION.md](../templates/VALIDATION.md)

## Execution Phases

### Phase 0 — Environment & Context Check
1. Locate `module/<name>/delivery/http/handler.go`.
2. Check `Register` method for existing routes.

### Phase 1 — Analysis & Decomposition
1. Define the DTO (Data Transfer Object) for the request.
2. Define validation rules (struct tags).
3. Determine the HTTP status codes (200, 201, 400, 404, 500).

### Phase 2 — Implementation
1. **Create Handler File**:
   - Create `module/<name>/delivery/http/<action>_<resource>.go` (e.g., `create_user.go`).
   - Use `http_handler.md` template.
2. **Implement Handler**:
   - Parse Body/Params.
   - Validate DTO.
   - Call Usecase.
   - Map Usecase error to HTTP error (`response.Error`).
   - Return JSON response (`response.Success`).
3. **Register Route**:
   - Edit `module/<name>/delivery/http/handler.go`.
   - Add `group.Method("/path", h.MethodName)`.

### Phase 3 — Validation
1. **Compile**: `go build`.
2. **Run**: Start server and test with `curl` or Postman.
3. **Checklist**: Execute `ai/checklists/definition_of_done.md`.

## Outputs
- New `delivery/http/<action>.go`
- Updated `delivery/http/handler.go`

## Failure Modes
- **Business Logic in Handler**: Performing DB calls or complex logic in the handler.
- **Raw Errors**: Returning raw Go errors to the client instead of standardized JSON.
- **Missing Validation**: Trusting user input without validation.
