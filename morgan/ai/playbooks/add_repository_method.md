# Playbook: Add Repository Method

## Purpose
Add a new method to a repository interface and its PostgreSQL implementation.

## Triggers
- New data access requirement identified.
- `add_usecase` playbook requires a missing repository method.

## Inputs
- **Method Signature**: Name, arguments, return types.
- **SQL Query**: The raw SQL or query logic required.
- **Target Repository**: The domain repository interface to modify.

## Preconditions
- Domain entity exists.
- Repository interface exists in `module/<name>/domain/repository.go`.
- Implementation exists in `module/<name>/repository/postgresql/`.

## Rules to Load
- [ai/rules/00_principles.md](../rules/00_principles.md)
- [ai/rules/10_repository.md](../rules/10_repository.md)
- [ai/rules/50_helpers.md](../rules/50_helpers.md)

## Templates to Use
- [ai/templates/repository.md](../templates/repository.md)

## Execution Phases

### Phase 0 — Environment & Context Check
1. Identify the target module and domain.
2. Locate `domain/repository.go` and `repository/postgresql/repository.go`.
3. Verify the entity struct definition.

### Phase 1 — Analysis & Decomposition
1. Define the method signature in the interface.
   - *Constraint*: Must accept `context.Context`.
   - *Constraint*: Must return `error`.
2. Draft the SQL query.
   - *Constraint*: Use `pgx` placeholders (`$1`, `$2`).

### Phase 2 — Implementation
1. **Update Interface**: Add method to `domain/repository.go`.
2. **Create Implementation File**:
   - If the method is complex, create a new file `repository/postgresql/<method_name>.go`.
   - If simple, append to `repository/postgresql/repository.go` (only if file is small).
   - *Preference*: Create new file for distinct operations (e.g., `find_by_email.go`).
3. **Implement Logic**:
   - Use `r.db.Query` or `r.db.QueryRow`.
   - Map SQL rows to Domain Entity.
   - Handle `pgx.ErrNoRows` -> return `nil, nil` or specific domain error based on requirements.
   - Wrap errors with context.

### Phase 3 — Validation
1. **Compile**: Run `go build ./...` to ensure interface satisfaction.
2. **Test**: Add unit/integration test for the new method (if testing playbook is active).
3. **Checklist**: Execute `ai/checklists/definition_of_done.md`.

## Outputs
- Updated `domain/repository.go`
- New/Updated `repository/postgresql/<method>.go`

## Failure Modes
- **Interface Mismatch**: Forgetting to update the struct implementation after changing the interface.
- **SQL Injection**: Concatenating strings instead of using placeholders.
- **Layer Violation**: Returning DTOs instead of Domain Entities.
