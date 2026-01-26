# Playbook: Add Usecase

## Purpose
Implement a business logic unit (usecase) that coordinates data and domain rules.

## Triggers
- New feature requirement.
- `add_http_endpoint` requires a handler to call business logic.

## Inputs
- **Usecase Name**: Verb-Noun (e.g., `CreateUser`, `ProcessOrder`).
- **Input/Output**: Request struct and Response struct/entity.
- **Business Rules**: Validation and logic steps.

## Preconditions
- Domain entities exist.
- Repository methods exist (or are planned).

## Rules to Load
- [ai/rules/00_principles.md](../rules/00_principles.md)
- [ai/rules/20_usecase_service.md](../rules/20_usecase_service.md)
- [ai/rules/50_helpers.md](../rules/50_helpers.md)

## Templates to Use
- [ai/templates/usecase.md](../templates/usecase.md)

## Execution Phases

### Phase 0 — Environment & Context Check
1. Locate `module/<name>/domain/usecase.go` (Interface).
2. Locate `module/<name>/usecase/` (Implementation).

### Phase 1 — Analysis & Decomposition
1. Define the `Usecase` interface method.
2. Define Input/Output structs (if not using raw entities).
3. Identify required dependencies (Repositories, other Usecases, Event Publishers).

### Phase 2 — Implementation
1. **Update Interface**: Add method to `domain/usecase.go`.
2. **Create Implementation**:
   - Create `module/<name>/usecase/<action>.go` (e.g., `create.go`).
   - Use the `usecase.md` template.
3. **Inject Dependencies**:
   - Ensure the `Usecase` struct has necessary fields.
   - Update constructor in `module.go` or `usecase/usecase.go` if needed.
4. **Implement Logic**:
   - Validate input.
   - Call repository methods.
   - Apply domain logic.
   - Publish events (if applicable).
   - Return result.

### Phase 3 — Validation
1. **Compile**: Ensure `Fx` dependency injection is valid (if applicable) and interfaces match.
2. **Test**: Create unit test mocking the repository.
3. **Checklist**: Execute `ai/checklists/definition_of_done.md`.

## Outputs
- Updated `domain/usecase.go`
- New `usecase/<action>.go`

## Failure Modes
- **God Object**: Adding too many methods to one file. Keep 1 file per major method.
- **Logic Leakage**: Putting HTTP logic (status codes) in the usecase.
- **Transaction Missing**: Not using transactions for multi-step DB operations (if supported).
