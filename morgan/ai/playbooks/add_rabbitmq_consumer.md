# Playbook: Add RabbitMQ Consumer

## Purpose
Process asynchronous messages from a RabbitMQ queue.

## Triggers
- New background processing requirement.
- Event-driven architecture implementation.

## Inputs
- **Queue Name**: The queue to consume from.
- **Exchange/Routing Key**: (If binding is required).
- **Message Schema**: JSON structure of the event.
- **Target Usecase**: Logic to execute upon message receipt.

## Preconditions
- Usecase exists.
- RabbitMQ infrastructure is configured.

## Rules to Load
- [ai/rules/00_principles.md](../rules/00_principles.md)
- [ai/rules/31_handler_consumer.md](../rules/31_handler_consumer.md)

## Templates to Use
- [ai/templates/consumer.md](../templates/consumer.md)

## Execution Phases

### Phase 0 — Environment & Context Check
1. Locate `module/<name>/delivery/messaging/` (or create if missing).
2. Check `module.go` for consumer registration.

### Phase 1 — Analysis & Decomposition
1. Define the Event DTO.
2. Determine idempotency requirements.
3. Determine error handling strategy (Retry, DLQ, Ack/Nack).

### Phase 2 — Implementation
1. **Create Consumer File**:
   - Create `module/<name>/delivery/messaging/consumer.go` (or specific file if multiple).
   - Use `consumer.md` template.
2. **Implement Logic**:
   - Unmarshal JSON body.
   - Call Usecase.
   - Handle Errors:
     - If temporary -> Nack (requeue).
     - If permanent -> Ack (discard) or DLQ.
     - If success -> Ack.
3. **Register Consumer**:
   - Ensure the consumer is started in `module.go` or `cmd/serve.go` (depending on app structure).
   - Usually registered via Fx lifecycle.

### Phase 3 — Validation
1. **Compile**: `go build`.
2. **Test**: Publish a test message to the queue and verify processing.
3. **Checklist**: Execute `ai/checklists/definition_of_done.md`.

## Outputs
- New/Updated `delivery/messaging/consumer.go`

## Failure Modes
- **Blocking Main Thread**: Consumer must run in a goroutine (handled by library usually, but verify).
- **Lost Messages**: Auto-acking before processing is complete.
- **Infinite Loops**: Requeuing permanent errors forever.
