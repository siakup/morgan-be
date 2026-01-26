# Helpers & Utilities

## Purpose
To document the usage of shared utilities for configuration, logging, and error handling.

## Observed Evidence
- **Config**: `config/config.go` uses `config:",squash"` and helper functions to expose sub-configs.
- **Logging**: `zerolog` is the standard logger.
- **Errors**: `yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors` provides typed errors.

## Rules (MUST)
- **MUST** define application configuration in `config/config.go`.
- **MUST** use `zerolog` for all logging.
- **MUST** use `yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors` for creating and checking errors.

## Conventions (SHOULD)
- **SHOULD** expose configuration sections via helper functions (e.g., `func Postgres(app *ApplicationConfig)`).

## Open Questions
- None.
