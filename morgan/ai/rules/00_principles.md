# Global Principles

## Purpose
To establish the fundamental architectural principles and design patterns that govern the codebase, ensuring consistency and maintainability.

## Observed Evidence
- **Dependency Injection**: `cmd/serve.go` and `module/user/module.go` extensively use `go.uber.org/fx`.
- **Clean Architecture**: Clear separation of `delivery`, `usecase`, `repository`, and `domain` packages in `module/user`.
- **Framework Usage**: Heavy reliance on `yuhuu.universitaspertamina.ac.id/siak/siakup/backend/framework` for core capabilities (fiber, postgres, redis, logger).

## Rules (MUST)
- **MUST** follow Clean Architecture principles:
  - Dependencies point inward: Delivery -> UseCase -> Repository -> Domain.
  - Domain layer **MUST NOT** depend on outer layers.
- **MUST** use `go.uber.org/fx` for dependency injection.
  - Each module **MUST** define a `Module` variable of type `fx.Option`.
  - Constructors **MUST** be provided via `fx.Provide` or `fx.Annotate`.
- **MUST** use `context.Context` as the first argument for all blocking or long-running methods (Repositories, UseCases).
- **MUST** use `yuhuu.universitaspertamina.ac.id/siak/siakup/backend/framework` modules for infrastructure setup (Logger, Config, Fiber, Postgres, etc.).

## Conventions (SHOULD)
- **SHOULD** prefer interface-based dependency injection to allow for mocking and loose coupling.
- **SHOULD** keep `cmd/` lightweight, delegating initialization logic to Fx modules.

## Open Questions
- None.
