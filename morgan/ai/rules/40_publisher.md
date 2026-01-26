# Publishers / Producers

## Purpose
To standardize the publication of domain events to message brokers, ensuring consistent event structures and tracing propagation.

## Observed Evidence
- **Library**: Uses `yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/publisher`.
- **Usage**: `module/user/usecase/create.go` publishes events after successful persistence.
- **Event Structure**: `module/audittrail/domain/event.go` (read in previous turns) defines event structures with `QueueExchange`, `QueueTopic`, `MessageIds`, etc.
- **Pattern**: Fire-and-forget (errors are logged but do not fail the parent operation).

## Rules (MUST)
- **MUST** use `yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/publisher` for publishing events.
- **MUST** implement the `publisher.Message` interface for event structs (methods: `Exchange`, `Topic`, `MessageId`, `ContentType`, `Body`).
- **MUST** publish events **AFTER** the primary business transaction (e.g., database commit) is successful.
- **MUST** propagate the Trace ID from the context to the event message (`MessageIds`).
- **MUST** use `zerolog` to log publishing failures.

## Conventions (SHOULD)
- **SHOULD** define event structures in the `domain` package.
- **SHOULD** NOT fail the main usecase operation if event publishing fails (unless distributed transaction guarantees are required).
- **SHOULD** use constants for Event Types and Service Names.

## Open Questions
- **Reliability**: Current implementation seems to be "at most once" or "at least once" depending on the library. Need to clarify if outbox pattern is needed for critical events.
