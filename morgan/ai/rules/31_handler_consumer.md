# Async Consumers

## Purpose
To standardize asynchronous message processing, ensuring reliable consumption and error handling.

## Observed Evidence
- **Signature**: Returns `consumer.Handler` (function closure).
- **Tracing**: Extracts trace ID from `msg.MessageId` and enriches context.
- **Ack/Nack**: Explicitly handles Ack/Nack.
- **Error Strategy**: Distinguishes between System errors (Requeue) and other errors (Discard).

## Rules (MUST)
- **MUST** return a `consumer.Handler` function.
- **MUST** extract the trace ID from the message (e.g., `msg.MessageId`) and add it to the context.
- **MUST** manually Acknowledge (`Ack`) or Negative Acknowledge (`Nack`) messages.
- **MUST** Requeue (`Nack(false, true)`) only for transient System errors.
- **MUST** Discard (`Nack(false, false)`) for malformed data or business logic errors to prevent infinite loops.

## Conventions (SHOULD)
- **SHOULD** log the incoming event and processing result with the trace ID.
- **SHOULD** unmarshal the message body into a typed struct before processing.

## Open Questions
- None.
