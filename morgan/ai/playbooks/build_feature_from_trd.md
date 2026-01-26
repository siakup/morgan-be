# Playbook: Build Feature from TRD

## Purpose
End-to-end implementation of a feature based on a Technical Requirement Document (TRD).

## Triggers
- A TRD is approved and ready for implementation.

## Inputs
- **TRD**: Document describing the feature, data models, APIs, and events.

## Preconditions
- TRD is complete and unambiguous.

## Rules to Load
- [ai/rules/00_principles.md](../rules/00_principles.md)
- [ai/rules/01_repo_structure.md](../rules/01_repo_structure.md)

## Templates to Use
- All templates as needed by sub-playbooks.

## Execution Phases

### Phase 0 — Environment & Context Check
1. Read the TRD fully.
2. Identify which modules are affected.
3. Create a plan listing the components to build (Entities, Repos, Usecases, Handlers).

### Phase 1 — Analysis & Decomposition
1. **Data Modeling**:
   - Identify new entities or changes to existing ones.
   - *Action*: Define structs in `domain/`.
2. **API Design**:
   - Identify new endpoints.
   - *Action*: Define DTOs.
3. **Dependency Graph**:
   - Map out: Handler -> Usecase -> Repository/Publisher.

### Phase 2 — Implementation (Iterative)
*Execute the following sub-playbooks in order:*

1. **Database Layer**:
   - For each new query/method:
   - Execute [add_repository_method.md](add_repository_method.md).

2. **Business Logic Layer**:
   - For each unit of logic:
   - Execute [add_usecase.md](add_usecase.md).

3. **Interface Layer (Choose applicable)**:
   - For HTTP APIs: Execute [add_http_endpoint.md](add_http_endpoint.md).
   - For Consumers: Execute [add_rabbitmq_consumer.md](add_rabbitmq_consumer.md).
   - For Cron Jobs: Execute [add_cron_job.md](add_cron_job.md).

4. **Wiring**:
   - Ensure all new components are registered in `module.go` (Fx options).

### Phase 3 — Validation
1. **Integration Test**: Verify the flow from Handler to DB.
2. **Checklist**: Execute `ai/checklists/definition_of_done.md`.

## Outputs
- Fully implemented feature across all layers.

## Failure Modes
- **Missing Dependencies**: Implementing Handler before Usecase. *Follow the order: Domain -> Repo -> Usecase -> Handler.*
- **Scope Creep**: Implementing features not in the TRD.
