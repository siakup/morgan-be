# Template: Event Publisher

## Location
```
module/<module>/domain/
└── event.go    # Event structure and publisher interface implementation
```

## Purpose
Defines event structures for publishing domain events to message broker. The Event:
- Represents state changes or significant domain actions
- Implements publisher.Message interface
- Contains routing information (exchange, topic)
- Includes trace ID for distributed tracing
- Serializes to JSON for transport

## Rules Applied
- Events live in domain layer (not infrastructure)
- Event structure separate from internal domain objects
- Immutable event data (use struct, not pointers for data fields)
- Publisher.Message interface implementation required
- No framework-specific dependencies

---

## Code Skeleton

### File: `module/<module>/domain/event.go`

```go
package domain

import "encoding/json"

// <Entity>Event represents the event message for <entity> operations.
// This event is published to the message broker for other services to consume.
type <Entity>Event struct {
	QueueExchange string ` + "`json:\"-\"`" + `          // RabbitMQ exchange name
	QueueTopic    string ` + "`json:\"-\"`" + `          // RabbitMQ routing key
	ContentTypes  string ` + "`json:\"-\"`" + `          // Content type (application/json)
	MessageIds    string ` + "`json:\"-\"`" + `          // Trace ID for correlation
	Event         string ` + "`json:\"event\"`" + `      // Event type (e.g., "CreateUser")
	Service       string ` + "`json:\"service\"`" + `    // Originating service
	Data          any    ` + "`json:\"data\"`" + `       // Event payload (domain object)
}

// New<Entity>Event creates a new instance of <Entity>Event.
// Parameters:
//   - messageId: Trace ID for distributed tracing
//   - event: Event type (use constants from usecase)
//   - service: Service name (use constant from usecase)
//   - data: Domain object representing the event data
func New<Entity>Event(messageId string, event string, service string, data any) *<Entity>Event {
	return &<Entity>Event{
		QueueExchange: "",                    // Default exchange (empty = default)
		QueueTopic:    "<routing-key>",       // e.g., "user.created", "order.updated"
		ContentTypes:  "application/json",
		MessageIds:    messageId,
		Event:         event,
		Service:       service,
		Data:          data,
	}
}

// Exchange returns the exchange name for the event.
// Implements publisher.Message interface.
func (e *<Entity>Event) Exchange() string {
	return e.QueueExchange
}

// Topic returns the routing key for the event.
// Implements publisher.Message interface.
func (e *<Entity>Event) Topic() string {
	return e.QueueTopic
}

// MessageId returns the trace ID for the event.
// Implements publisher.Message interface.
func (e *<Entity>Event) MessageId() string {
	return e.MessageIds
}

// ContentType returns the content type for the event.
// Implements publisher.Message interface.
func (e *<Entity>Event) ContentType() string {
	return e.ContentTypes
}

// Body returns the JSON-encoded event body.
// Implements publisher.Message interface.
func (e *<Entity>Event) Body() []byte {
	body, _ := json.Marshal(e)
	return body
}
```

### Alternative: Strongly-Typed Event Data

For better type safety, use a specific struct instead of `any`:

```go
package domain

import "encoding/json"

// <Entity>EventData represents the typed payload for <entity> events.
type <Entity>EventData struct {
	Id      string  ` + "`json:\"id\"`" + `
	<Field1> <Type1> ` + "`json:\"<field1>\"`" + `
	<Field2> <Type2> ` + "`json:\"<field2>\"`" + `
	// Add relevant fields that consumers need
}

// <Entity>Event represents the event message.
type <Entity>Event struct {
	QueueExchange string            ` + "`json:\"-\"`" + `
	QueueTopic    string            ` + "`json:\"-\"`" + `
	ContentTypes  string            ` + "`json:\"-\"`" + `
	MessageIds    string            ` + "`json:\"-\"`" + `
	Event         string            ` + "`json:\"event\"`" + `
	Service       string            ` + "`json:\"service\"`" + `
	Data          <Entity>EventData ` + "`json:\"data\"`" + ` // Strongly typed
}

// New<Entity>Event creates a new instance with type safety.
func New<Entity>Event(messageId string, event string, service string, data <Entity>EventData) *<Entity>Event {
	return &<Entity>Event{
		QueueExchange: "",
		QueueTopic:    "<routing-key>",
		ContentTypes:  "application/json",
		MessageIds:    messageId,
		Event:         event,
		Service:       service,
		Data:          data,
	}
}

// ... implement publisher.Message interface methods
```

### Alternative: Multiple Event Types for Same Entity

```go
package domain

import "encoding/json"

// Base event structure (shared)
type baseEvent struct {
	QueueExchange string
	QueueTopic    string
	ContentTypes  string
	MessageIds    string
	Event         string
	Service       string
}

// <Entity>CreatedEvent for creation events
type <Entity>CreatedEvent struct {
	baseEvent
	Data <Entity> ` + "`json:\"data\"`" + `
}

// New<Entity>CreatedEvent creates a creation event
func New<Entity>CreatedEvent(messageId string, service string, data <Entity>) *<Entity>CreatedEvent {
	return &<Entity>CreatedEvent{
		baseEvent: baseEvent{
			QueueExchange: "",
			QueueTopic:    "<entity>.created",
			ContentTypes:  "application/json",
			MessageIds:    messageId,
			Event:         "Create<Entity>",
			Service:       service,
		},
		Data: data,
	}
}

// Implement publisher.Message interface...
func (e *<Entity>CreatedEvent) Exchange() string     { return e.QueueExchange }
func (e *<Entity>CreatedEvent) Topic() string        { return e.QueueTopic }
func (e *<Entity>CreatedEvent) MessageId() string    { return e.MessageIds }
func (e *<Entity>CreatedEvent) ContentType() string  { return e.ContentTypes }
func (e *<Entity>CreatedEvent) Body() []byte {
	body, _ := json.Marshal(e)
	return body
}

// <Entity>UpdatedEvent for update events
type <Entity>UpdatedEvent struct {
	baseEvent
	Data <Entity> ` + "`json:\"data\"`" + `
}

// New<Entity>UpdatedEvent creates an update event
func New<Entity>UpdatedEvent(messageId string, service string, data <Entity>) *<Entity>UpdatedEvent {
	return &<Entity>UpdatedEvent{
		baseEvent: baseEvent{
			QueueExchange: "",
			QueueTopic:    "<entity>.updated",
			ContentTypes:  "application/json",
			MessageIds:    messageId,
			Event:         "Update<Entity>",
			Service:       service,
		},
		Data: data,
	}
}

// Implement publisher.Message interface...
```

### Usage in UseCase

```go
// In module/<module>/usecase/create.go

package usecase

import (
	"context"
	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/helper"
)

func (u *UseCase) Create(ctx context.Context, <entity> *domain.<Entity>) error {
	// ... store entity logic ...

	// Publish event
	traceId := helper.GetTraceID(ctx)
	event := domain.New<Entity>Event(traceId, EventTypeCreate<Entity>, Service<Module>, <entity>)
	
	if err := u.eventPublisher.Publish(ctx, event); err != nil {
		logger.Error().
			Str("func", "eventPublisher.Publish").
			Err(err).
			Msg("failed to publish event")
		// Don't fail operation if publishing fails
	}

	return nil
}
```

---

## Tests Required

### File: `module/<module>/domain/event_test.go`

1. **Test New<Entity>Event**:
   - Creates event with correct fields
   - Sets default values properly
   - Accepts all required parameters

2. **Test Exchange()**:
   - Returns configured exchange name

3. **Test Topic()**:
   - Returns configured routing key

4. **Test MessageId()**:
   - Returns provided trace ID

5. **Test ContentType()**:
   - Returns "application/json"

6. **Test Body()**:
   - Returns valid JSON
   - Includes all public fields
   - Excludes JSON-tagged "-" fields
   - Correctly serializes nested data

**Test Structure**:
```go
func TestNew<Entity>Event(t *testing.T) {
	messageId := "trace-123"
	event := "CreateUser"
	service := "user"
	data := &domain.User{Id: "1", Name: "Test"}

	evt := domain.New<Entity>Event(messageId, event, service, data)

	assert.Equal(t, "", evt.Exchange())
	assert.Equal(t, "user.created", evt.Topic())
	assert.Equal(t, messageId, evt.MessageId())
	assert.Equal(t, "application/json", evt.ContentType())
}

func TestEventBody(t *testing.T) {
	evt := domain.New<Entity>Event("trace-1", "CreateUser", "user", 
		&domain.User{Id: "1", Name: "Test"})

	body := evt.Body()

	// Verify JSON structure
	var result map[string]any
	err := json.Unmarshal(body, &result)
	assert.NoError(t, err)
	assert.Equal(t, "CreateUser", result["event"])
	assert.Equal(t, "user", result["service"])
	assert.NotNil(t, result["data"])
	
	// Should NOT include metadata fields (tagged with "-")
	assert.NotContains(t, string(body), "QueueExchange")
	assert.NotContains(t, string(body), "MessageIds")
}
```

---

## Notes / Pitfalls

### Critical Points

1. **publisher.Message Interface**: Must implement ALL methods
   ```go
   type Message interface {
       Exchange() string
       Topic() string
       MessageId() string
       ContentType() string
       Body() []byte
   }
   ```

2. **JSON Tags**:
   - Use ` + "`json:\"-\"`" + ` for fields that should NOT be serialized
   - Metadata fields (QueueExchange, MessageIds, etc.) should not appear in body
   - Only Event, Service, and Data should be in the JSON body

3. **Routing Key Convention**:
   ```
   <entity>.<action>        // e.g., user.created, order.updated
   <service>.<entity>       // e.g., user.user, payment.transaction
   <domain>.<entity>.<action> // e.g., ecommerce.order.shipped
   ```

4. **Event Naming**:
   - Use PascalCase for event type: "CreateUser", "UpdateOrder"
   - Use past tense for completed actions: "UserCreated", "OrderShipped"
   - Be consistent across the system

5. **Data Payload**:
   - Include only necessary information (don't send entire entity if not needed)
   - Consider backwards compatibility when changing structure
   - Avoid sensitive data (passwords, tokens, etc.)

6. **Trace ID**:
   - Always use the request trace ID from context
   - Enables end-to-end tracing across services
   - Use `helper.GetTraceID(ctx)` to extract

7. **Error Handling in Body()**:
   ```go
   func (e *Event) Body() []byte {
       body, _ := json.Marshal(e)
       return body
       // Ignoring error is OK here since struct is known to be serializable
   }
   ```

### Common Mistakes

- Forgetting to implement all publisher.Message methods
- Including metadata fields in JSON body (QueueExchange, MessageIds)
- Not using trace ID for correlation
- Hardcoding exchange/routing keys (should be configurable or constant)
- Publishing events before state is persisted (publish AFTER)
- Sending entire entity when only ID is needed
- Not handling publisher errors in usecase
- Using mutable pointers in event data (events should be immutable)
- Inconsistent event naming across services

### Exchange and Routing Strategy

**Fanout Exchange** (all consumers receive):
```go
QueueExchange: "events.fanout"
QueueTopic:    ""  // Ignored in fanout
```

**Topic Exchange** (pattern matching):
```go
QueueExchange: "events.topic"
QueueTopic:    "user.created"
// Consumer binds with pattern like "user.*" or "*.created"
```

**Direct Exchange** (exact match):
```go
QueueExchange: "events.direct"
QueueTopic:    "user-service"
```

**Default Exchange** (queue name):
```go
QueueExchange: ""  // Default
QueueTopic:    "user.events"  // Queue name
```

### Event Schema Evolution

Consider versioning for breaking changes:

```go
type <Entity>EventV2 struct {
	QueueExchange string
	QueueTopic    string
	ContentTypes  string
	MessageIds    string
	Event         string
	Service       string
	Version       int    ` + "`json:\"version\"`" + ` // Add version field
	Data          any    ` + "`json:\"data\"`" + `
}

func New<Entity>EventV2(...) *<Entity>EventV2 {
	return &<Entity>EventV2{
		// ...
		Version: 2,
		// ...
	}
}
```

### Idempotency

Consumers should use MessageId for deduplication:
- Store processed message IDs in cache/database
- Skip processing if message ID already seen
- Set appropriate TTL for deduplication cache

### Dead Letter Exchange (DLX)

Configure in infrastructure, not in event structure:
```go
// Events that fail processing go to DLX
// Configure in RabbitMQ setup, not in domain code
```

### Monitoring

Track metrics for:
- Events published per type
- Publishing failures
- Publishing latency
- Message size
