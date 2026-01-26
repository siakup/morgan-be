# AI Templates - System Architecture Templates

## Overview

This directory contains **reusable, production-ready templates** for implementing features in this Go microservice architecture. These templates enforce architectural boundaries, best practices, and patterns observed in the existing codebase.

**Purpose**: Generate consistent, maintainable code by following established patterns
**Audience**: Developers and AI assistants implementing new features or modules

---

## 📁 Available Templates

| Template | File | Purpose |
|----------|------|---------|
| **Repository** | [repository.md](./repository.md) | Data access layer with PostgreSQL |
| **UseCase/Service** | [usecase.md](./usecase.md) | Business logic orchestration |
| **HTTP Handler** | [http_handler.md](./http_handler.md) | REST API endpoints with Fiber |
| **Consumer** | [consumer.md](./consumer.md) | RabbitMQ message consumers |
| **Publisher** | [publisher.md](./publisher.md) | Event publishing structures |
| **Cron Handler** | [cron_handler.md](./cron_handler.md) | Scheduled job handlers |

---

## 🏗️ Architecture Overview

This system follows **Clean Architecture** principles with clear layer separation:

```
┌─────────────────────────────────────────────────────────┐
│                    Delivery Layer                       │
│  (HTTP Handlers, Consumers, Cron Jobs)                 │
│  - Translates external protocols to domain calls       │
│  - No business logic                                    │
└────────────────┬────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────┐
│                   UseCase Layer                         │
│  - Business logic orchestration                         │
│  - Error transformation                                 │
│  - Event publishing                                     │
│  - Tracing                                              │
└────────────────┬────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────┐
│                 Repository Layer                        │
│  - Data access abstraction                             │
│  - PostgreSQL operations                               │
│  - Entity ↔ Domain object mapping                      │
└─────────────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────┐
│                  Domain Layer                           │
│  - Interfaces (Repository, UseCase)                    │
│  - Domain objects                                       │
│  - Events                                               │
│  - Pure business entities                              │
└─────────────────────────────────────────────────────────┘
```

### Layer Boundaries (CRITICAL)

**Domain Layer**:
- ✅ Pure Go structs and interfaces
- ✅ Business rules and validations
- ❌ NO framework dependencies (HTTP, DB, MQ)
- ❌ NO infrastructure concerns

**Repository Layer**:
- ✅ Implements domain.Repository interface
- ✅ Uses pgx/v5 for PostgreSQL
- ✅ Returns domain objects, NOT database entities
- ❌ NO business logic
- ❌ NO HTTP or messaging types

**UseCase Layer**:
- ✅ Orchestrates business operations
- ✅ Calls repository and external services
- ✅ Transforms errors (DB errors → AppErrors)
- ✅ Publishes domain events
- ❌ NO HTTP status codes, request/response DTOs
- ❌ NO direct database access

**Delivery Layer** (HTTP/Consumer/Cron):
- ✅ Protocol-specific handling (HTTP, AMQP, Cron)
- ✅ Request validation and response formatting
- ✅ Delegates to UseCase
- ❌ NO business logic
- ❌ NO direct repository access
- ❌ NO database operations

---

## 🚀 How to Use Templates

### Step 1: Choose Your Template

Identify which layer(s) you need to implement:
- **New feature with HTTP API**: Repository → UseCase → HTTP Handler
- **Background processing**: Repository → UseCase → Consumer
- **Scheduled task**: UseCase → Cron Handler
- **Event publishing**: UseCase → Publisher

### Step 2: Replace Placeholders

Each template uses consistent placeholders:

| Placeholder | Replace With | Example |
|-------------|-------------|---------|
| `<module>` | Module name (lowercase) | `user`, `order`, `payment` |
| `<Module>` | Module name (PascalCase) | `User`, `Order`, `Payment` |
| `<entity>` | Entity name (lowercase) | `user`, `product` |
| `<Entity>` | Entity name (PascalCase) | `User`, `Product` |
| `<Entities>` | Plural entity name | `Users`, `Products` |
| `<entities>` | Plural entity (lowercase) | `users`, `products` |
| `<field1>`, `<Field1>` | Struct field names | `name`, `Name` |
| `<Type1>` | Go type | `string`, `int`, `time.Time` |
| `<table>` | Database table name | `users`, `orders` |
| `<action>` | Action/operation name | `create`, `update`, `approve` |
| `<Action>` | Action (PascalCase) | `Create`, `Update`, `Approve` |
| `<org>` | GitHub organization | `beruang`, `mycompany` |
| `<project>` | Project name | `go-project`, `myapp` |
| `<routing-key>` | RabbitMQ routing key | `user.created`, `order.shipped` |
| `<job-name>` | Cron job identifier | `daily-cleanup`, `sync-data` |

### Step 3: Follow the Structure

Templates specify exact file locations:
```
module/<module>/
├── domain/
│   ├── repository.go      # Interface definitions
│   ├── usecase.go         # UseCase interface
│   └── event.go           # Event structures (if publishing)
├── repository/
│   └── postgresql/
│       ├── repository.go  # Struct and constructor
│       ├── find_*.go      # Query operations
│       ├── store.go       # Create
│       ├── update.go      # Update
│       └── delete.go      # Delete
├── usecase/
│   ├── usecase.go         # UseCase struct
│   ├── create.go          # Business logic per action
│   ├── get.go
│   ├── find_all.go
│   ├── put.go
│   └── delete.go
└── delivery/
    ├── http/
    │   ├── handler.go               # Handler struct and routing
    │   ├── create_<entity>.go       # Per-endpoint handlers
    │   ├── get_<entity>_by_id.go
    │   ├── get_<entities>.go
    │   ├── update_<entity>.go
    │   └── delete_<entity>.go
    ├── messaging/
    │   └── consumer.go              # Message consumer
    └── cron/
        └── handler.go               # Scheduled jobs
```

### Step 4: Implement Tests

Each template includes test requirements. Tests should cover:
- Happy path scenarios
- Error conditions
- Edge cases
- Mock dependencies appropriately

---

## 📋 Template Details

### 1. Repository Layer ([repository.md](./repository.md))

**When to use**: Accessing any persistent data store (PostgreSQL)

**Key principles**:
- Returns domain objects, not database entities
- One operation per file
- Uses named parameters (`@param`) for SQL
- No business logic

**Example**:
```go
// Domain interface
type Repository interface {
    FindByID(ctx context.Context, id string) (*User, error)
}

// Implementation
func (r *Repository) FindByID(ctx context.Context, id string) (*User, error) {
    rows, _ := r.db.Query(ctx, queryFindById, pgx.NamedArgs{"id": id})
    record, _ := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[UserEntity])
    return object.Parse[*UserEntity, *User](record)
}
```

---

### 2. UseCase Layer ([usecase.md](./usecase.md))

**When to use**: Implementing any business operation

**Key principles**:
- Orchestrates repository calls
- Transforms errors (pgx.ErrNoRows → NotFound)
- Publishes domain events after state changes
- Uses OpenTelemetry tracing
- No framework-specific types

**Example**:
```go
func (u *UseCase) Create(ctx context.Context, user *domain.User) error {
    ctx, span := u.tracer.Start(ctx, "Create")
    defer span.End()
    
    if err := u.repository.Store(ctx, user); err != nil {
        return errors.InternalServerError("failed to store user")
    }
    
    event := domain.NewUserEvent(helper.GetTraceID(ctx), "CreateUser", "user", user)
    u.eventPublisher.Publish(ctx, event)
    
    return nil
}
```

---

### 3. HTTP Handler ([http_handler.md](./http_handler.md))

**When to use**: Exposing REST API endpoints

**Key principles**:
- Uses Fiber framework
- Request/Response DTOs separate from domain
- Centralized error handling
- Context enrichment (trace ID, logger)
- No business logic

**Example**:
```go
func (h *Handler) CreateUser(c *gofiber.Ctx) error {
    var req CreateUserRequest
    c.BodyParser(&req)
    
    user := domain.User{Id: req.Id, Name: req.Name}
    if err := h.useCase.Create(c.UserContext(), &user); err != nil {
        return h.handleError(c, err)
    }
    
    return c.Status(201).JSON(responses.Success(user, "User created"))
}
```

---

### 4. Consumer ([consumer.md](./consumer.md))

**When to use**: Processing asynchronous messages from RabbitMQ

**Key principles**:
- Returns `consumer.Handler` function
- Proper Ack/Nack based on error type
- System errors → requeue
- Business errors → discard
- Trace ID from message ID

**Example**:
```go
func Consume(useCase domain.UseCase) consumer.Handler {
    return func(ctx context.Context, msg amqp.Delivery) error {
        var event Event
        json.Unmarshal(msg.Body, &event)
        
        ctx = helper.WithTraceID(ctx, msg.MessageId)
        
        if err := useCase.Process(ctx, event.Data); err != nil {
            if isSystemError(err) {
                msg.Nack(false, true)  // Requeue
            } else {
                msg.Nack(false, false) // Discard
            }
            return err
        }
        
        return msg.Ack(false)
    }
}
```

---

### 5. Publisher ([publisher.md](./publisher.md))

**When to use**: Publishing domain events to message broker

**Key principles**:
- Lives in domain layer
- Implements publisher.Message interface
- Includes routing information
- JSON serialization

**Example**:
```go
type UserEvent struct {
    QueueExchange string ` + "`json:\"-\"`" + `
    QueueTopic    string ` + "`json:\"-\"`" + `
    MessageIds    string ` + "`json:\"-\"`" + `
    Event         string ` + "`json:\"event\"`" + `
    Service       string ` + "`json:\"service\"`" + `
    Data          any    ` + "`json:\"data\"`" + `
}

func (e *UserEvent) Body() []byte {
    body, _ := json.Marshal(e)
    return body
}
```

---

### 6. Cron Handler ([cron_handler.md](./cron_handler.md))

**When to use**: Implementing scheduled/periodic tasks

**Key principles**:
- Returns `func()` for scheduler
- Never panic (graceful error handling)
- Context timeout
- Distributed lock for singleton jobs
- Trace ID generation

**Example**:
```go
func DailyCleanupHandler(useCase domain.UseCase) func() {
    return func() {
        traceId := uuid.New().String()
        ctx := helper.WithTraceID(context.Background(), traceId)
        ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
        defer cancel()
        
        if err := useCase.PerformCleanup(ctx); err != nil {
            log.Error().Err(err).Msg("cleanup failed")
            return // Don't panic
        }
        
        log.Info().Msg("cleanup completed")
    }
}
```

---

## ⚙️ Technology Stack

Based on `go.mod`:

| Component | Technology | Version |
|-----------|-----------|---------|
| **Language** | Go | 1.25.5 |
| **Web Framework** | Fiber | v2.52.9 |
| **Database** | PostgreSQL (pgx) | v5.7.4 |
| **Message Broker** | RabbitMQ (amqp091-go) | v1.10.0 |
| **Logging** | zerolog | v1.33.0 |
| **Tracing** | OpenTelemetry | v1.39.0 |
| **Dependency Injection** | Fx | v1.24.0 |

---

## 🎯 Common Patterns

### Error Handling

**Repository Layer** (return raw errors):
```go
return pgx.ErrNoRows
```

**UseCase Layer** (transform to app errors):
```go
if errors.Is(err, pgx.ErrNoRows) {
    return appErrors.NotFound("user not found")
}
return appErrors.InternalServerError("operation failed")
```

**HTTP Layer** (map to HTTP status):
```go
func (h *Handler) handleError(c *fiber.Ctx, err error) error {
    if appErr, ok := err.(*errors.AppError); ok {
        return c.Status(appErr.Code).JSON(responses.Fail(appErr.Type, appErr.Message))
    }
    return c.Status(500).JSON(responses.Fail("SYSTEM_ERROR", err.Error()))
}
```

### Context Propagation

Always pass context through the call chain:
```
HTTP Request → Handler (extract trace ID) 
  → UseCase (create span) 
    → Repository (use context)
```

### Tracing Pattern

```go
func (u *UseCase) Operation(ctx context.Context) error {
    ctx, span := u.tracer.Start(ctx, "Operation")
    defer span.End()
    
    logger := zerolog.Ctx(ctx)
    // ... logic
}
```

### Event Publishing Pattern

```go
// After successful state change
traceId := helper.GetTraceID(ctx)
event := domain.NewEvent(traceId, eventType, service, data)
if err := u.eventPublisher.Publish(ctx, event); err != nil {
    logger.Error().Err(err).Msg("failed to publish event")
    // Don't fail the operation
}
```

---

## ✅ Best Practices

### DO

- ✅ Follow layer boundaries strictly
- ✅ Use domain objects in interfaces
- ✅ Transform errors at each layer
- ✅ Publish events after state changes
- ✅ Use structured logging with context
- ✅ Implement comprehensive tests
- ✅ Use one file per operation/handler
- ✅ Apply distributed tracing
- ✅ Handle context cancellation

### DON'T

- ❌ Mix business logic in handlers/consumers
- ❌ Leak infrastructure types across layers
- ❌ Ignore error handling
- ❌ Skip event publishing
- ❌ Hardcode configuration values
- ❌ Use positional parameters in SQL
- ❌ Panic in cron jobs or consumers
- ❌ Return database entities from repository

---

## 🔍 Validation Checklist

Before considering code complete, verify:

- [ ] Templates comply with all layer boundaries
- [ ] No feature-specific logic in templates
- [ ] Placeholders used consistently
- [ ] Folder paths match repository structure
- [ ] Error handling at every layer
- [ ] Logging with structured context
- [ ] Tests cover happy path and errors
- [ ] Domain objects used in interfaces
- [ ] Events published after state changes
- [ ] Context propagated throughout
- [ ] No framework types in domain/usecase
- [ ] Safe to reuse across multiple features

---

## 📚 Additional Resources

### Module Structure Example

See existing modules for reference:
- `module/user/` - Full CRUD with HTTP handlers
- `module/audittrail/` - Consumer-based processing

### Helper Libraries

Located in `yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib`:
- `errors` - Standardized error types
- `responses` - HTTP response formatting
- `publisher` - Event publishing interface
- `consumer` - Message consumer interface
- `helper` - Utility functions (trace ID, etc.)
- `types` - Common types (Pagination, etc.)

### Configuration

Typically managed in:
- `config/config.go` - Configuration structure
- Environment variables for runtime config

---

## 🤝 Contributing

When adding new templates:
1. Extract patterns from 3+ existing implementations
2. Remove all feature-specific logic
3. Use consistent placeholders
4. Include comprehensive tests section
5. Document common pitfalls
6. Update this README

---

## 📝 License

These templates follow the same license as the main project (see [LICENSE](../LICENSE)).

---

**Generated**: December 17, 2025  
**Maintainer**: Engineering Team  
**Status**: Production-Ready
