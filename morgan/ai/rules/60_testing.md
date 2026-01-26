# Testing Strategy

## Purpose
To establish a standard testing strategy, ensuring code reliability and preventing regressions.

## Observed Evidence
- **Current State**: No test files (`_test.go`) were observed in the scanned directories (`module/user`, `module/audittrail`).
- **Framework**: `go.mod` does not list `testify` or `gomock`, implying reliance on the standard library or that tests are missing.

## Rules (MUST)
- **MUST** place unit tests in the same package as the code being tested (e.g., `usecase_test.go` in `usecase` package).
- **MUST** use the standard Go `testing` package.
- **MUST** name test functions `Test<Function>_<Scenario>`.
- **MUST** mock external dependencies (Repositories, Publishers) when testing UseCases.
- **MUST** mock UseCases when testing Handlers.

## Conventions (SHOULD)
- **SHOULD** use Table-Driven Tests for covering multiple scenarios.
- **SHOULD** aim for high coverage in `usecase` and `domain` logic.
- **SHOULD** use `testcontainers` or similar for integration testing Repositories (if introduced).

## Open Questions
- **Mocking Library**: Should we introduce `github.com/stretchr/testify` and `go.uber.org/mock`? (Recommended).
- **Integration Tests**: Where should end-to-end tests live? (Suggested: `tests/integrations/`).
