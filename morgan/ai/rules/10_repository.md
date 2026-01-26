# Repository Layer

## Purpose
To standardize the persistence layer implementation, ensuring consistent data access and separation of concerns.

## Observed Evidence
- **Implementation**: `module/user/repository/postgresql/repository.go` uses `pgx/v5`.
- **Interface**: `module/user/domain/repository.go` defines the contract.
- **Pattern**: Repository methods return domain objects, not database entities (implied by `domain.Repository` signature).

## Rules (MUST)
- **MUST** implement the repository interface defined in the `domain` package.
- **MUST** use `github.com/jackc/pgx/v5` for PostgreSQL interactions.
- **MUST** accept `context.Context` as the first argument.
- **MUST NOT** leak database-specific types (like `pgx.Rows` or SQL null types) out of the repository layer.
- **MUST** map database entities to domain objects before returning.

## Conventions (SHOULD)
- **SHOULD** use named parameters (e.g., `@id`, `@name`) in SQL queries for clarity and safety.
- **SHOULD** separate complex queries into distinct files (e.g., `find_by_id.go`, `store.go`) if the repository grows large.

## Open Questions
- Handling of transactions is not explicitly observed in the sample code. Need to clarify the transaction management strategy (e.g., passing `pgx.Tx` in context or method arguments).
