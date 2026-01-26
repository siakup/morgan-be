# Template: Message Consumer (Delivery Layer)

## Location
```
module/<module>/delivery/messaging/
└── consumer.go    # Consumer handler function
```

## Purpose
Implements the message consumer for processing asynchronous messages from RabbitMQ. The Consumer:
- Consumes messages from a specific queue/exchange
- Unmarshals message payload to event structure
- Extracts trace ID for distributed tracing
- Delegates processing to UseCase
- Handles acknowledgment (Ack) and negative acknowledgment (Nack)
- Implements retry logic based on error type

## Rules Applied
- Layer boundary enforcement: Consumer MUST NOT contain business logic
- No direct repository access (only through UseCase)
- Proper message acknowledgment based on error type
- Context enrichment with trace ID and structured logging
- Error classification (system errors → requeue, business errors → discard)
- Consumer returns a handler function, not a struct

---

## Code Skeleton

### File: `module/<module>/delivery/messaging/consumer.go`

```go
package messaging

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/consumer"
	errors2 "yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/helper"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

// <Entity>Event represents the structure of the message received from the queue.
// This should match the published event structure from other services.
type <Entity>Event struct {
	Event   string          ` + "`json:\"event\"`" + `
	Service string          ` + "`json:\"service\"`" + `
	Data    <EventData>     ` + "`json:\"data\"`" + `
	// Add fields matching your event payload
}

// <EventData> represents the actual payload data within the event
type <EventData> struct {
	Id      string  ` + "`json:\"id\"`" + `
	<Field1> <Type1> ` + "`json:\"<field1>\"`" + `
	<Field2> <Type2> ` + "`json:\"<field2>\"`" + `
	// Add fields as needed
}

// Consume processes messages from the <module> queue.
// It unmarshals the message, extracts the trace ID, and delegates processing to the use case.
func Consume(useCase domain.UseCase) consumer.Handler {
	return func(ctx context.Context, msg amqp.Delivery) error {
		// Step 1: Unmarshal message
		var event <Entity>Event
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Error().
				Err(err).
				Str("body", string(msg.Body)).
				Msg("failed to unmarshal message")
			
			// Invalid JSON = permanent failure, don't requeue
			_ = msg.Nack(false, false)
			return err
		}

		// Step 2: Extract trace ID and enrich context
		messageId := msg.MessageId
		ctx = helper.WithTraceID(ctx, messageId)
		logger := log.With().Str("request_id", messageId).Logger()
		ctx = logger.WithContext(ctx)

		// Step 3: Log incoming event
		logger.Info().
			Str("event", event.Event).
			Str("service", event.Service).
			Msg("processing event")

		// Step 4: Transform event to domain object (if needed)
		domainObj := domain.<Entity>{
			Id:      event.Data.Id,
			<Field1>: event.Data.<Field1>,
			<Field2>: event.Data.<Field2>,
		}

		// Step 5: Delegate to use case
		err := useCase.Process(ctx, domainObj)
		if err != nil {
			// Classify error type for retry logic
			var vErr *errors2.AppError
			if errors.As(err, &vErr) {
				// System errors should be retried (requeue)
				if vErr.Type == errors2.ErrorTypeSystem {
					logger.Error().
						Err(err).
						Str("error_type", string(vErr.Type)).
						Msg("system error, requeuing message")
					
					_ = msg.Nack(false, true) // requeue
					return err
				}
			}

			// Business errors should not be retried (discard)
			logger.Error().
				Err(err).
				Msg("business error, discarding message")
			
			_ = msg.Nack(false, false) // don't requeue
			return err
		}

		// Step 6: Acknowledge successful processing
		return msg.Ack(false)
	}
}
```

### Alternative: Consumer with Multiple Event Types

```go
package messaging

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/consumer"
	errors2 "yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/helper"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

const (
	EventTypeCreate<Entity> = "Create<Entity>"
	EventTypeUpdate<Entity> = "Update<Entity>"
	EventTypeDelete<Entity> = "Delete<Entity>"
)

// <Entity>Event represents the base message structure
type <Entity>Event struct {
	Event   string ` + "`json:\"event\"`" + `
	Service string ` + "`json:\"service\"`" + `
	Data    any    ` + "`json:\"data\"`" + ` // Use interface{} for polymorphic data
}

// Consume processes messages from the <module> queue with event type routing.
func Consume(useCase domain.UseCase) consumer.Handler {
	return func(ctx context.Context, msg amqp.Delivery) error {
		var event <Entity>Event
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			_ = msg.Nack(false, false)
			return err
		}

		messageId := msg.MessageId
		ctx = helper.WithTraceID(ctx, messageId)
		logger := log.With().Str("request_id", messageId).Logger()
		ctx = logger.WithContext(ctx)

		// Route based on event type
		var err error
		switch event.Event {
		case EventTypeCreate<Entity>:
			err = handleCreate<Entity>(ctx, useCase, event.Data)
		case EventTypeUpdate<Entity>:
			err = handleUpdate<Entity>(ctx, useCase, event.Data)
		case EventTypeDelete<Entity>:
			err = handleDelete<Entity>(ctx, useCase, event.Data)
		default:
			logger.Warn().
				Str("event", event.Event).
				Msg("unknown event type, discarding")
			_ = msg.Nack(false, false)
			return nil
		}

		if err != nil {
			var vErr *errors2.AppError
			if errors.As(err, &vErr) {
				if vErr.Type == errors2.ErrorTypeSystem {
					_ = msg.Nack(false, true)
					return err
				}
			}
			_ = msg.Nack(false, false)
			return err
		}

		return msg.Ack(false)
	}
}

func handleCreate<Entity>(ctx context.Context, useCase domain.UseCase, data any) error {
	// Parse data to specific structure
	// Call useCase.Create(...)
	return nil
}

func handleUpdate<Entity>(ctx context.Context, useCase domain.UseCase, data any) error {
	// Parse data to specific structure
	// Call useCase.Update(...)
	return nil
}

func handleDelete<Entity>(ctx context.Context, useCase domain.UseCase, data any) error {
	// Parse data to specific structure
	// Call useCase.Delete(...)
	return nil
}
```

### Integration: Register Consumer in Module

```go
// In module/<module>/module.go

package <module>

import (
	"github.com/<org>/<project>/module/<module>/delivery/messaging"
	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/consumer"
)

type Module struct {
	// ... other fields
}

// RegisterConsumers registers message consumers for this module
func (m *Module) RegisterConsumers(consumer *consumer.Consumer, useCase domain.UseCase) {
	consumer.Subscribe("<queue-name>", messaging.Consume(useCase))
	
	// Register multiple consumers if needed
	// consumer.Subscribe("<another-queue>", messaging.ConsumeOther(useCase))
}
```

---

## Tests Required

### File: `module/<module>/delivery/messaging/consumer_test.go`

1. **Test Consume - Valid Message**:
   - Valid JSON payload
   - UseCase processes successfully
   - Message is acknowledged (Ack)

2. **Test Consume - Invalid JSON**:
   - Malformed JSON
   - Message is negatively acknowledged without requeue

3. **Test Consume - System Error**:
   - UseCase returns system error (ErrorTypeSystem)
   - Message is negatively acknowledged WITH requeue

4. **Test Consume - Business Error**:
   - UseCase returns business error (not ErrorTypeSystem)
   - Message is negatively acknowledged WITHOUT requeue

5. **Test Consume - Event Type Routing** (if applicable):
   - Correct handler called based on event type
   - Unknown event type is discarded

**Test Structure**:
```go
func TestConsume(t *testing.T) {
	tests := []struct {
		name          string
		messageBody   string
		messageId     string
		mockSetup     func(*MockUseCase)
		expectAck     bool
		expectNack    bool
		expectRequeue bool
	}{
		{
			name:        "valid message",
			messageBody: ` + "`" + `{"event":"CreateUser","service":"user","data":{"id":"123"}}` + "`" + `,
			messageId:   "trace-123",
			mockSetup: func(m *MockUseCase) {
				m.On("Process", mock.Anything, mock.Anything).Return(nil)
			},
			expectAck: true,
		},
		{
			name:        "invalid json",
			messageBody: "invalid-json",
			messageId:   "trace-124",
			expectNack:  true,
			expectRequeue: false,
		},
		// Add more test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(MockUseCase)
			if tt.mockSetup != nil {
				tt.mockSetup(mockUC)
			}

			handler := Consume(mockUC)

			mockMsg := &MockDelivery{
				body:      []byte(tt.messageBody),
				messageId: tt.messageId,
			}

			err := handler(context.Background(), mockMsg)

			if tt.expectAck {
				assert.True(t, mockMsg.acked)
			}
			if tt.expectNack {
				assert.True(t, mockMsg.nacked)
				assert.Equal(t, tt.expectRequeue, mockMsg.requeue)
			}
		})
	}
}
```

---

## Notes / Pitfalls

### Critical Points

1. **Message Acknowledgment Strategy**:
   ```go
   // ✅ GOOD - Proper error classification
   if vErr.Type == errors2.ErrorTypeSystem {
       msg.Nack(false, true)  // Requeue for retry
   } else {
       msg.Nack(false, false) // Discard, don't retry
   }
   
   // ❌ BAD - Always requeuing
   msg.Nack(false, true)  // Will cause infinite retry loop
   ```

2. **Trace ID Propagation**:
   - Always use `msg.MessageId` as trace ID
   - Attach to context for downstream propagation
   - Include in structured logging

3. **Error Handling**:
   - Invalid JSON → Nack without requeue (permanent failure)
   - System errors → Nack with requeue (transient failure)
   - Business errors → Nack without requeue (validation failure)
   - Unknown event types → Nack without requeue (skip)

4. **Context Enrichment**:
   ```go
   ctx = helper.WithTraceID(ctx, messageId)
   logger := log.With().Str("request_id", messageId).Logger()
   ctx = logger.WithContext(ctx)
   ```

5. **Return Signature**:
   - Consumer function returns `consumer.Handler`
   - Handler is `func(context.Context, amqp.Delivery) error`
   - Don't create a struct; use function closure

6. **Idempotency**:
   - Messages may be delivered multiple times
   - UseCase should handle duplicate processing
   - Consider using message ID for deduplication

7. **Dead Letter Queue (DLQ)**:
   - Configure DLQ for messages that exceed retry count
   - System should handle poison messages
   - Monitor DLQ for investigation

### Common Mistakes

- Not distinguishing between system and business errors (causes infinite retry)
- Ignoring Ack/Nack errors (can lead to message loss)
- Not logging message details before discarding
- Putting business logic in consumer (belongs in usecase)
- Using generic `any` for data without proper type assertion
- Not handling context cancellation
- Missing trace ID propagation
- Hardcoding queue names (should come from config)
- Not setting up DLQ for failed messages

### RabbitMQ-Specific Notes

**amqp.Delivery Methods**:
- `msg.Ack(multiple bool)` - Acknowledge message
  - `false`: Ack only this message
  - `true`: Ack this and all prior unacked messages
  
- `msg.Nack(multiple bool, requeue bool)` - Negative acknowledge
  - `multiple=false`: Nack only this message
  - `requeue=true`: Put back in queue for retry
  - `requeue=false`: Discard or route to DLQ

**Message Properties**:
- `msg.MessageId` - Unique message identifier (use as trace ID)
- `msg.Body` - Message payload (JSON)
- `msg.Headers` - Message headers (metadata)
- `msg.ContentType` - Usually "application/json"

### Retry Strategy

Recommended approach:
1. **Transient failures** (DB connection, network): Requeue
2. **Permanent failures** (invalid data, validation): Discard
3. **Max retries**: Configure in RabbitMQ (x-max-retries header)
4. **Backoff**: Use RabbitMQ delayed message plugin or TTL

Example configuration:
```go
// In consumer setup
consumer.Subscribe("queue-name", handler, consumer.WithMaxRetries(3))
```

### Observability

Always log:
- Message received (with trace ID)
- Event type and service
- Processing result (success/failure)
- Error details (for failures)
- Ack/Nack decisions

Metrics to track:
- Messages consumed
- Messages acknowledged
- Messages requeued
- Processing latency
- Error rate by type
