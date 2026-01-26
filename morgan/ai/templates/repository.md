# Template: Repository Layer

## Location
```
module/<module>/repository/postgresql/
├── repository.go           # Repository struct and constructor
├── find_by_<field>.go     # Query methods (one per file)
├── find_all.go            # List/filter query
├── store.go               # Create operation
├── update.go              # Update operation
└── delete.go              # Delete operation (soft or hard)
```

## Purpose
Implements the data access layer for a specific entity. The repository:
- Implements the domain.Repository interface
- Handles PostgreSQL-specific operations using pgx/v5
- Performs entity-to-domain object mapping
- Isolates database concerns from business logic
- Uses named parameters for SQL queries

## Rules Applied
- Layer boundary enforcement: Repository must NOT expose pgx types or database entities outside its package
- Use domain objects (domain.<Entity>) in method signatures
- Use database entities (domain.<Entity>Entity) internally only
- Each operation in a separate file for maintainability
- Use named parameters (@param) for SQL queries
- Handle errors appropriately (return db errors as-is for usecase to interpret)

---

## Code Skeleton

### File: `module/<module>/domain/repository.go`

```go
package domain

import (
	"context"
	"time"
)

// <Entity>Entity represents the schema in the database.
// Tags: db for database column mapping, map for object transformation
type <Entity>Entity struct {
	Id        string     `db:"id" map:"Id"`
	// Add fields matching your database schema
	<Field1>  <Type1>    `db:"<field1>" map:"<Field1>"`
	<Field2>  <Type2>    `db:"<field2>" map:"<Field2>"`
	CreatedAt time.Time  `db:"created_at" map:"CreatedAt"`
	UpdatedAt time.Time  `db:"updated_at" map:"UpdatedAt"`
	DeletedAt *time.Time `db:"deleted_at" map:"DeletedAt"`
}

// Repository defines the methods for interacting with the <entity> storage.
// It uses Domain Objects (*<Entity>) in its signature to decouple the caller from DB Entities.
type Repository interface {
	FindAll(ctx context.Context, filter <Entity>Filter) ([]*<Entity>, error)
	FindByID(ctx context.Context, id string) (*<Entity>, error)
	// Add custom finders as needed:
	// FindBy<Field>(ctx context.Context, <field> <Type>) (*<Entity>, error)
	Store(ctx context.Context, <entity> *<Entity>) error
	Update(ctx context.Context, <entity> *<Entity>) error
	Delete(ctx context.Context, id string) error
}
```

### File: `module/<module>/repository/postgresql/repository.go`

```go
package postgresql

import (
	"github.com/<org>/<project>/module/<module>/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ domain.Repository = (*Repository)(nil)

// Repository implements the domain.Repository interface for PostgreSQL.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new instance of the PostgreSQL Repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}
```

### File: `module/<module>/repository/postgresql/find_by_id.go`

```go
package postgresql

import (
	"context"

	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/object"
	"github.com/jackc/pgx/v5"
)

var queryFindById = ` + "`" + `
	SELECT
		id, <field1>, <field2>, created_at, updated_at
	FROM <table>
	WHERE id = @id AND deleted_at IS NULL
	LIMIT 1
` + "`" + `

// FindByID retrieves a single <entity> by their ID.
func (r *Repository) FindByID(ctx context.Context, id string) (*domain.<Entity>, error) {
	rows, err := r.db.Query(ctx, queryFindById, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	record, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.<Entity>Entity])
	if err != nil {
		return nil, err
	}

	return object.Parse[*domain.<Entity>Entity, *domain.<Entity>](object.TagDB, object.TagObject, record)
}
```

### File: `module/<module>/repository/postgresql/find_all.go`

```go
package postgresql

import (
	"context"

	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/object"
	"github.com/jackc/pgx/v5"
)

var queryFindAll = ` + "`" + `
	SELECT
		id, <field1>, <field2>, created_at, updated_at
	FROM <table>
	WHERE deleted_at IS NULL
	-- Add WHERE conditions based on filter
	ORDER BY created_at DESC
	LIMIT @limit OFFSET @offset
` + "`" + `

// FindAll retrieves a list of <entities> based on the provided filter.
func (r *Repository) FindAll(ctx context.Context, filter domain.<Entity>Filter) ([]*domain.<Entity>, error) {
	rows, err := r.db.Query(ctx, queryFindAll, pgx.NamedArgs{
		"limit":  filter.PageSize,
		"offset": filter.Offset,
		// Add additional filter parameters
	})
	if err != nil {
		return nil, err
	}

	records, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[domain.<Entity>Entity])
	if err != nil {
		return nil, err
	}

	return object.ParseSlice[*domain.<Entity>Entity, *domain.<Entity>](object.TagDB, object.TagObject, records)
}
```

### File: `module/<module>/repository/postgresql/store.go`

```go
package postgresql

import (
	"context"

	"github.com/<org>/<project>/module/<module>/domain"
	"github.com/jackc/pgx/v5"
)

var queryStore = ` + "`" + `
	INSERT INTO <table> (
		id, <field1>, <field2>, created_at, updated_at
	) VALUES (
		@id, @<field1>, @<field2>, NOW(), NOW()
	)
` + "`" + `

// Store persists a new <entity> to the database.
func (r *Repository) Store(ctx context.Context, <entity> *domain.<Entity>) error {
	_, err := r.db.Exec(ctx, queryStore, pgx.NamedArgs{
		"id":      <entity>.Id,
		"<field1>": <entity>.<Field1>,
		"<field2>": <entity>.<Field2>,
		// Add all fields
	})

	return err
}
```

### File: `module/<module>/repository/postgresql/update.go`

```go
package postgresql

import (
	"context"

	"github.com/<org>/<project>/module/<module>/domain"
	"github.com/jackc/pgx/v5"
)

var queryUpdate = ` + "`" + `
	UPDATE <table>
	SET
		<field1> = @<field1>,
		<field2> = @<field2>,
		updated_at = NOW()
	WHERE id = @id AND deleted_at IS NULL
` + "`" + `

// Update modifies an existing <entity> record.
func (r *Repository) Update(ctx context.Context, <entity> *domain.<Entity>) error {
	_, err := r.db.Exec(ctx, queryUpdate, pgx.NamedArgs{
		"id":      <entity>.Id,
		"<field1>": <entity>.<Field1>,
		"<field2>": <entity>.<Field2>,
		// Add all updatable fields
	})

	return err
}
```

### File: `module/<module>/repository/postgresql/delete.go`

```go
package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Soft delete (recommended)
var queryDelete = ` + "`" + `
	UPDATE <table>
	SET deleted_at = NOW()
	WHERE id = @id AND deleted_at IS NULL
` + "`" + `

// Hard delete (use only if soft delete is not required)
// var queryDelete = ` + "`DELETE FROM <table> WHERE id = @id`" + `

// Delete removes a <entity> from the database.
// This performs a soft delete by setting the deleted_at timestamp.
func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, queryDelete, pgx.NamedArgs{
		"id": id,
	})

	return err
}
```

---

## Tests Required

### File: `module/<module>/repository/postgresql/repository_test.go`

1. **Test NewRepository**: Verify constructor creates valid instance
2. **Test FindByID**: 
   - Happy path: existing record
   - Error case: non-existent ID (should return pgx.ErrNoRows)
   - Error case: invalid ID format
3. **Test FindAll**:
   - Empty result set
   - Multiple records with pagination
   - Filter application
4. **Test Store**:
   - Successful insert
   - Duplicate ID error
   - Constraint violation
5. **Test Update**:
   - Successful update
   - Non-existent ID
6. **Test Delete**:
   - Successful deletion
   - Already deleted record
   - Non-existent ID

**Setup Requirements**:
- Use test database or Docker container
- Implement database fixtures/seeders
- Clean up test data after each test
- Use table-driven tests for multiple scenarios

---

## Notes / Pitfalls

### Critical Points

1. **Layer Boundary**: Repository MUST NOT leak database types
   - ❌ BAD: `func FindByID() (*pgxpool.Row, error)`
   - ✅ GOOD: `func FindByID() (*domain.Entity, error)`

2. **Entity vs Domain Object**:
   - `domain.<Entity>Entity` = database representation (with db tags)
   - `domain.<Entity>` = domain model (with object tags)
   - Always transform Entity → Domain Object before returning

3. **Named Parameters**: Always use `@param` syntax, never positional `$1, $2`
   ```go
   // ✅ GOOD
   pgx.NamedArgs{"id": id, "name": name}
   
   // ❌ BAD
   db.Exec(query, id, name)
   ```

4. **Error Handling**: Return raw database errors
   - Don't wrap or transform errors in repository
   - Let usecase layer interpret pgx.ErrNoRows, constraint violations, etc.

5. **Soft Delete**: Always filter by `deleted_at IS NULL` in SELECT queries

6. **SQL Injection**: pgx protects against injection when using named parameters, but always validate input at usecase layer

7. **Transactions**: For multi-table operations, accept `pgx.Tx` instead of `*pgxpool.Pool`:
   ```go
   type Repository interface {
       StoreWithTx(ctx context.Context, tx pgx.Tx, entity *Entity) error
   }
   ```

8. **Null Handling**: Use pointers for nullable fields:
   ```go
   DeletedAt *time.Time `db:"deleted_at"`
   ```

### Common Mistakes

- Mixing business logic in repository (belongs in usecase)
- Not using the object.Parse helper for entity transformation
- Forgetting to check `deleted_at` in queries
- Hardcoding pagination limits (should come from filter)
- Not handling pgx.ErrNoRows separately from other errors
