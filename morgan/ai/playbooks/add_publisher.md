# Playbook: Add Publisher

## Purpose
Publish domain events to RabbitMQ.

## Triggers
- Usecase needs to notify other systems/modules.

## Inputs
- **Exchange Name**: Target exchange.
- **Routing Key**: Topic or key.
- **Event Schema**: Data to send.

## Preconditions
- RabbitMQ connection available.

## Rules to Load
- [ai/rules/00_principles.md](../rules/00_principles.md)
- [ai/rules/40_publisher.md](../rules/40_publisher.md)

## Templates to Use
- [ai/templates/publisher.md](../templates/publisher.md)

## Execution Phases

### Phase 0 — Environment & Context Check
1. Identify where the event originates (Usecase).

### Phase 1 — Analysis & Decomposition
1. Define the Event struct (if not shared).
2. Determine if publishing is transactional (Outbox pattern) or fire-and-forget. *Default to fire-and-forget unless Outbox specified.*

### Phase 2 — Implementation
1. **Create Publisher Interface** (Optional but recommended for testing):
   - Define in `domain/repository.go` or `domain/publisher.go`.
2. **Implement Publisher**:
   - Create `module/<name>/delivery/messaging/publisher.go` (or `infrastructure/messaging`).
   - Use `publisher.md` template.
   - Use `amqp` library to publish.
3. **Integrate into Usecase**:
   - Inject publisher interface into Usecase.
   - Call `Publish` method after successful domain logic.

### Phase 3 — Validation
1. **Compile**: `go build`.
2. **Test**: Mock publisher in Usecase tests.
3. **Checklist**: Execute `ai/checklists/definition_of_done.md`.

## Outputs
- New Publisher implementation.
- Updated Usecase.

## Failure Modes
- **Publishing before Commit**: If DB transaction rolls back, event is still sent (Phantom event). *Mitigation: Publish after commit or use Outbox.*
- **Connection Leaks**: Not reusing channels properly.
