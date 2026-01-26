# HTTP Handlers

## Purpose
To standardize HTTP API implementation using Fiber, ensuring consistent routing, middleware usage, and response formatting.

## Observed Evidence
- **Framework**: Uses `github.com/gofiber/fiber/v2`.
- **Routing**: Implements `RegisterRoutes(app *gofiber.App)`.
- **Middleware**: Extracts `requestid`, sets up trace context, and logger.
- **Response**: Uses `responses.Success` and `responses.Fail` helpers.
- **Error Handling**: Uses `handleError` helper to map `AppError` to HTTP status codes.

## Rules (MUST)
- **MUST** implement `RegisterRoutes` to define endpoints.
- **MUST** use `c.UserContext()` to propagate context to the UseCase.
- **MUST** extract and propagate the Request ID (Trace ID).
- **MUST** use `responses.Success` and `responses.Fail` for JSON responses.
- **MUST** map `errors.AppError` to appropriate HTTP status codes (e.g., `NotFound` -> 404).
- **MUST NOT** contain business logic; delegate to the UseCase.

## Conventions (SHOULD)
- **SHOULD** define Request and Response structs within the handler package/file.
- **SHOULD** use `c.BodyParser` for parsing JSON bodies.

## Open Questions
- None.
