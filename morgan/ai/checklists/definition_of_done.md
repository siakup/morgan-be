# Definition of Done

## Purpose
To ensure every code change meets the quality standards of the repository before being merged.

## Checklist
- [ ] **Compiles**: Code builds without errors (`make build`).
- [ ] **Linted**: Code passes linter checks (`make lint` or `golangci-lint run`).
- [ ] **Tested**: Unit tests added/updated and passing (`go test ./...`).
- [ ] **Clean Architecture**: No layer violations (e.g., Domain depending on Delivery).
- [ ] **Error Handling**: Errors are typed (`AppError`) and not swallowed.
- [ ] **Observability**: Tracing started for new operations; logs include context.
- [ ] **Configuration**: No hardcoded secrets or config values.
- [ ] **Formatting**: Code formatted with `gofmt`.
