# Template: UseCase/Service Layer

## Location
```
module/<module>/usecase/
├── usecase.go         # UseCase struct, interface implementation, constructor
├── <action1>.go       # Business logic for action 1 (e.g., create.go)
├── <action2>.go       # Business logic for action 2 (e.g., get.go)
├── <action3>.go       # Business logic for action 3 (e.g., update.go)
└── ...                # One file per action/method
```

Also requires:
```
module/<module>/domain/usecase.go    # Interface definition
```

## Purpose
Implements the core business logic of a module. The UseCase:
- Orchestrates operations between repository, external services, and publishers
- Implements domain.UseCase interface
- Handles error transformation (db errors → app errors)
- Publishes domain events for audit/integration
- Implements distributed tracing
- Contains NO framework-specific code (HTTP, messaging, etc.)

## Rules Applied
- Layer boundary enforcement: UseCase accepts and returns domain objects only
- No leaking of infrastructure concerns (HTTP status, MQ details, DB types)
- Error mapping: database/infrastructure errors → standardized app errors
- Event-driven: Publish events after state changes
- Observability: OpenTelemetry tracing for all operations
- Single responsibility: One file per business action

---

## Code Skeleton

### File: `module/<module>/domain/usecase.go`

```go
package domain

import (
	"context"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/types"
)

// <Entity> represents the domain object for a <entity>.
// It contains pure business data and logic, decoupled from database annotations.
type <Entity> struct {
	Id      string ` + "`object:\"id\"`" + `
	<Field1> <Type1> ` + "`object:\"<field1>\"`" + `
	<Field2> <Type2> ` + "`object:\"<field2>\"`" + `
	// Add business-relevant fields only
}

// <Entity>Filter represents the filter options for fetching <entities>.
type <Entity>Filter struct {
	types.Pagination
	// Add filter fields
	<FilterField1> *<Type1>
	<FilterField2> string
}

// UseCase defines the business logic for the <module> module.
type UseCase interface {
	FindAll(ctx context.Context, filter <Entity>Filter) ([]*<Entity>, error)
	Get(ctx context.Context, id string) (*<Entity>, error)
	Create(ctx context.Context, <entity> *<Entity>) error
	Put(ctx context.Context, <entity> *<Entity>) error
	Delete(ctx context.Context, id string) error
	// Add custom business methods as needed
	// <CustomAction>(ctx context.Context, params <Params>) (<Result>, error)
}
```

### File: `module/<module>/usecase/usecase.go`

```go
package usecase

import (
	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/publisher"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var _ domain.UseCase = (*UseCase)(nil)

const (
	Service<Module>         = "<module>"
	EventTypeCreate<Entity> = "Create<Entity>"
	EventTypeUpdate<Entity> = "Update<Entity>"
	EventTypeDelete<Entity> = "Delete<Entity>"
	// Add custom event types
)

// UseCase implements the logic for <module> management.
type UseCase struct {
	repository     domain.Repository
	eventPublisher *publisher.Publisher
	tracer         trace.Tracer
	// Add dependencies (other repositories, external services, etc.)
}

// NewUseCase creates a new instance of <Module> UseCase.
func NewUseCase(repository domain.Repository, eventPublisher *publisher.Publisher) *UseCase {
	return &UseCase{
		repository:     repository,
		eventPublisher: eventPublisher,
		tracer:         otel.Tracer("<module>"),
	}
}
```

### File: `module/<module>/usecase/create.go`

```go
package usecase

import (
	"context"

	auditDomain "github.com/<org>/<project>/module/audittrail/domain"
	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/helper"
	"github.com/rs/zerolog"
)

// Create persists a new <entity> record.
func (u *UseCase) Create(ctx context.Context, <entity> *domain.<Entity>) error {
	ctx, span := u.tracer.Start(ctx, "Create")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	// Step 1: Business validation (if needed)
	// if <entity>.<Field> == "" {
	//     return errors.BadRequest("<field> is required")
	// }

	// Step 2: Call repository to persist
	if err := u.repository.Store(ctx, <entity>); err != nil {
		logger.Error().
			Str("func", "repository.Store").
			Err(err).
			Msg("failed to store <entity>")

		return errors.InternalServerError("failed to store <entity>")
	}

	// Step 3: Publish audit event
	traceId := helper.GetTraceID(ctx)
	auditTrailEvent := auditDomain.NewAuditTrailEvent(traceId, EventTypeCreate<Entity>, Service<Module>, <entity>)
	if err := u.eventPublisher.Publish(ctx, auditTrailEvent); err != nil {
		logger.Error().
			Str("func", "eventPublisher.Publish").
			Err(err).
			Msg("failed to publish audit trail event")

		// Don't return error - event publishing failure shouldn't fail the operation
		// Consider implementing a dead-letter queue or retry mechanism
	}

	return nil
}
```

### File: `module/<module>/usecase/get.go`

```go
package usecase

import (
	"context"
	"errors"

	"github.com/<org>/<project>/module/<module>/domain"
	appErrors "yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

// Get finds a <entity> by their unique identifier.
func (u *UseCase) Get(ctx context.Context, id string) (*domain.<Entity>, error) {
	ctx, span := u.tracer.Start(ctx, "Get")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	<entity>, err := u.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.NotFound("<entity> not found")
		}

		logger.Error().
			Str("func", "repository.FindByID").
			Err(err).
			Msg("failed to find <entity> by id")
		return nil, appErrors.InternalServerError("failed to find <entity> by id")
	}

	return <entity>, nil
}
```

### File: `module/<module>/usecase/find_all.go`

```go
package usecase

import (
	"context"

	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"github.com/rs/zerolog"
)

// FindAll retrieves a list of <entities> based on filter criteria.
func (u *UseCase) FindAll(ctx context.Context, filter domain.<Entity>Filter) ([]*domain.<Entity>, error) {
	ctx, span := u.tracer.Start(ctx, "FindAll")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	<entities>, err := u.repository.FindAll(ctx, filter)
	if err != nil {
		logger.Error().
			Str("func", "repository.FindAll").
			Err(err).
			Msg("failed to find <entities>")

		return nil, errors.InternalServerError("failed to find <entities>")
	}

	return <entities>, nil
}
```

### File: `module/<module>/usecase/put.go`

```go
package usecase

import (
	"context"

	auditDomain "github.com/<org>/<project>/module/audittrail/domain"
	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/helper"
	"github.com/rs/zerolog"
)

// Put updates an existing <entity> record.
func (u *UseCase) Put(ctx context.Context, <entity> *domain.<Entity>) error {
	ctx, span := u.tracer.Start(ctx, "Put")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	// Step 1: Verify entity exists (optional but recommended)
	_, err := u.repository.FindByID(ctx, <entity>.Id)
	if err != nil {
		logger.Error().
			Str("func", "repository.FindByID").
			Err(err).
			Msg("failed to find <entity> by id")

		return errors.NotFound("<entity> not found")
	}

	// Step 2: Business validation (if needed)
	// Add your validation logic here

	// Step 3: Perform update
	if err := u.repository.Update(ctx, <entity>); err != nil {
		logger.Error().
			Str("func", "repository.Update").
			Err(err).
			Msg("failed to update <entity>")

		return errors.InternalServerError("failed to update <entity>")
	}

	// Step 4: Publish audit event
	traceId := helper.GetTraceID(ctx)
	auditTrailEvent := auditDomain.NewAuditTrailEvent(traceId, EventTypeUpdate<Entity>, Service<Module>, <entity>)
	if err := u.eventPublisher.Publish(ctx, auditTrailEvent); err != nil {
		logger.Error().
			Str("func", "eventPublisher.Publish").
			Err(err).
			Msg("failed to publish audit trail event")
	}

	return nil
}
```

### File: `module/<module>/usecase/delete.go`

```go
package usecase

import (
	"context"

	auditDomain "github.com/<org>/<project>/module/audittrail/domain"
	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/helper"
	"github.com/rs/zerolog"
)

// Delete removes a <entity> from the system.
func (u *UseCase) Delete(ctx context.Context, id string) error {
	ctx, span := u.tracer.Start(ctx, "Delete")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	// Step 1: Verify entity exists (optional but recommended)
	<entity>, err := u.repository.FindByID(ctx, id)
	if err != nil {
		logger.Error().
			Str("func", "repository.FindByID").
			Err(err).
			Msg("failed to find <entity> by id")

		return errors.NotFound("<entity> not found")
	}

	// Step 2: Perform deletion
	if err := u.repository.Delete(ctx, id); err != nil {
		logger.Error().
			Str("func", "repository.Delete").
			Err(err).
			Msg("failed to delete <entity>")

		return errors.InternalServerError("failed to delete <entity>")
	}

	// Step 3: Publish audit event
	traceId := helper.GetTraceID(ctx)
	auditTrailEvent := auditDomain.NewAuditTrailEvent(traceId, EventTypeDelete<Entity>, Service<Module>, <entity>)
	if err := u.eventPublisher.Publish(ctx, auditTrailEvent); err != nil {
		logger.Error().
			Str("func", "eventPublisher.Publish").
			Err(err).
			Msg("failed to publish audit trail event")
	}

	return nil
}
```

---

## Tests Required

### File: `module/<module>/usecase/usecase_test.go`

1. **Test Create**:
   - Successful creation
   - Repository failure (should return internal server error)
   - Event publishing failure (should not fail the operation)
   - Business validation failures

2. **Test Get**:
   - Existing entity
   - Non-existent entity (should return NotFound)
   - Repository error (should return InternalServerError)

3. **Test FindAll**:
   - Empty result
   - Multiple results
   - Filter application
   - Repository error

4. **Test Put**:
   - Successful update
   - Non-existent entity
   - Repository failure

5. **Test Delete**:
   - Successful deletion
   - Non-existent entity
   - Repository failure

**Mocking Requirements**:
- Mock `domain.Repository` interface
- Mock `publisher.Publisher`
- Use testify/mock or similar mocking framework
- Verify trace span creation (optional)

**Test Structure**:
```go
func TestUseCase_Create(t *testing.T) {
	tests := []struct {
		name          string
		input         *domain.Entity
		mockSetup     func(*MockRepository, *MockPublisher)
		expectedError error
	}{
		// Test cases
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			// Execute
			// Assert
		})
	}
}
```

---

## Notes / Pitfalls

### Critical Points

1. **Layer Boundary**: UseCase MUST NOT know about HTTP, messaging, or database specifics
   - ❌ BAD: `func Create(c *fiber.Ctx) error`
   - ❌ BAD: `func Create() (*pgx.Row, error)`
   - ✅ GOOD: `func Create(ctx context.Context, entity *domain.Entity) error`

2. **Error Mapping**: Transform infrastructure errors to application errors
   ```go
   // ✅ GOOD
   if errors.Is(err, pgx.ErrNoRows) {
       return appErrors.NotFound("entity not found")
   }
   return appErrors.InternalServerError("operation failed")
   
   // ❌ BAD - leaking pgx error
   return err
   ```

3. **Context Propagation**: Always pass context through the chain
   - Enables tracing, cancellation, and deadline propagation
   - Extract logger from context: `zerolog.Ctx(ctx)`

4. **Tracing Pattern**:
   ```go
   ctx, span := u.tracer.Start(ctx, "OperationName")
   defer span.End()
   ```

5. **Event Publishing**:
   - Always publish AFTER successful state change
   - Event publishing failure should be logged but NOT fail the operation
   - Use trace ID for correlation: `helper.GetTraceID(ctx)`

6. **Business Validation**:
   - Validate at usecase level, not repository
   - Return `errors.BadRequest()` for validation failures

7. **Idempotency**: Consider idempotency for Create operations
   - Check if entity exists before inserting
   - Use unique identifiers (UUID) generated by client

8. **Transaction Handling**: For multi-repository operations:
   ```go
   tx, err := u.db.Begin(ctx)
   defer tx.Rollback(ctx)
   
   // Multiple operations
   u.repository.StoreWithTx(ctx, tx, entity1)
   u.otherRepo.StoreWithTx(ctx, tx, entity2)
   
   tx.Commit(ctx)
   ```

### Common Mistakes

- Mixing presentation logic (HTTP status codes, response formatting) in usecase
- Not transforming database errors to app errors
- Making repository calls without error handling
- Failing operation when event publishing fails (should be resilient)
- Not using structured logging with context
- Returning database entities instead of domain objects
- Putting validation logic in repository layer
- Not creating separate trace spans for operations
- Ignoring context cancellation

### Dependencies

UseCase typically depends on:
- `domain.Repository` (same module)
- `publisher.Publisher` (for events)
- Other `domain.Repository` interfaces (cross-module)
- External service clients (if needed)

DO NOT depend on:
- HTTP handlers
- Message consumers
- Database connections directly (use Repository)
- Framework-specific types (fiber.Ctx, amqp.Delivery, etc.)
