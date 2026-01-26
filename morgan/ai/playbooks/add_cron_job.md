# Playbook: Add Cron Job

## Purpose
Schedule recurring background tasks.

## Triggers
- Periodic cleanup, reporting, or sync requirements.

## Inputs
- **Schedule**: Cron expression (e.g., `0 0 * * *`).
- **Task**: Usecase to execute.

## Preconditions
- Usecase exists.
- Cron library/scheduler is set up in `cmd/` or `module/`.

## Rules to Load
- [ai/rules/00_principles.md](../rules/00_principles.md)

## Templates to Use
- [ai/templates/cron_handler.md](../templates/cron_handler.md)

## Execution Phases

### Phase 0 — Environment & Context Check
1. Check if a Cron scheduler is already running (e.g., `robfig/cron`).

### Phase 1 — Analysis & Decomposition
1. Determine if the job needs to be a singleton (distributed lock) or can run on all instances.
2. *Constraint*: If running in Kubernetes with multiple replicas, ensure idempotency or use a locking mechanism.

### Phase 2 — Implementation
1. **Create Handler**:
   - Create `module/<name>/delivery/cron/handler.go`.
   - Use `cron_handler.md` template.
2. **Implement Logic**:
   - Call Usecase.
   - Log start/end and errors.
3. **Register Job**:
   - In `module.go` or `cmd/root.go`, register the cron job with the scheduler.

### Phase 3 — Validation
1. **Compile**: `go build`.
2. **Test**: Manually trigger the job or wait for schedule.
3. **Checklist**: Execute `ai/checklists/definition_of_done.md`.

## Outputs
- New `delivery/cron/handler.go`
- Registration code.

## Failure Modes
- **Overlapping Jobs**: Job takes longer than interval. *Mitigation: Skip if running.*
- **Silent Failures**: Errors not logged or alerted.
