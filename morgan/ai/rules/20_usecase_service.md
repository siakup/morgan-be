# UseCase / Service Layer

## Purpose
To define the business logic layer, ensuring proper orchestration, error handling, and observability.

## Observed Evidence
- **Tracing**: `module/user/usecase/create.go` starts a trace span: `ctx, span := u.tracer.Start(ctx, "Create")`.
- **Logging**: Uses `zerolog.Ctx(ctx)` for contextual logging.
- **Error Handling**: Returns `errors.InternalServerError` or similar from `yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors`.
- **Events**: Publishes events using `eventPublisher`.

## Rules (MUST)
- **MUST** implement the UseCase interface defined in the `domain` package.
- **MUST** start a new OpenTelemetry trace span for each public method.
- **MUST** use `zerolog.Ctx(ctx)` to obtain a logger with context fields (trace ID, etc.).
- **MUST** return errors using `yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors` types (e.g., `AppError`).
- **MUST NOT** return raw database errors; wrap or transform them into domain errors.
- **MUST** publish domain events for significant state changes (e.g., "UserCreated").

## Conventions (SHOULD)
- **SHOULD** validate input data before processing or delegating to the repository.
- **SHOULD** keep business logic pure and independent of delivery mechanisms (HTTP, AMQP).

## Open Questions
- None.
