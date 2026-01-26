# Playbook: Build Feature from OpenAPI and Database Migrations

## Purpose
Implement a feature when the contract (OpenAPI) and storage (SQL Migrations) are already defined.
This is a "meet-in-the-middle" approach.

## Triggers
- OpenAPI spec provided.
- SQL migration files provided.

## Inputs
- **OpenAPI Spec**: YAML/JSON defining endpoints.
- **SQL Migrations**: `.sql` files defining tables.

## Preconditions
- Files are present in workspace.

## Rules to Load
- [ai/rules/00_principles.md](../rules/00_principles.md)
- [ai/rules/01_repo_structure.md](../rules/01_repo_structure.md)

## Templates to Use
- [ai/templates/repository.md](../templates/repository.md)
- [ai/templates/usecase.md](../templates/usecase.md)
- [ai/templates/http_handler.md](../templates/http_handler.md)

## Execution Phases

### Phase 0 — Environment & Context Check
1. Analyze SQL to understand the Data Model.
2. Analyze OpenAPI to understand the API Contract.

### Phase 1 — Analysis & Decomposition
1. **Map SQL to Domain**:
   - Create Domain Entities matching SQL tables.
2. **Map OpenAPI to DTOs**:
   - Create Request/Response structs matching API spec.
3. **Gap Analysis**:
   - Identify what logic is needed to transform SQL Data <-> API Data.

### Phase 2 — Implementation
1. **Repository Layer**:
   - Execute [add_repository_method.md](add_repository_method.md) for all CRUD operations implied by the API.
2. **Usecase Layer**:
   - Execute [add_usecase.md](add_usecase.md) to bridge the gap.
   - Implement validation logic defined in OpenAPI (e.g., required fields, max length).
3. **HTTP Layer**:
   - Execute [add_http_endpoint.md](add_http_endpoint.md).
   - Ensure route paths and methods match OpenAPI exactly.
   - Ensure JSON tags match OpenAPI property names.

### Phase 3 — Validation
1. **Contract Testing**: Verify endpoints match OpenAPI spec (if tools available).
2. **Checklist**: Execute `ai/checklists/definition_of_done.md`.

## Outputs
- Complete implementation adhering to both contracts.

## Failure Modes
- **Field Mismatch**: JSON tags not matching OpenAPI.
- **Type Mismatch**: SQL types (e.g., `TIMESTAMP`) not mapping correctly to Go types (`time.Time`).
