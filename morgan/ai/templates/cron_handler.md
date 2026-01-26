# Template: Cron Handler (Scheduled Jobs)

## Location
```
module/<module>/delivery/cron/
└── handler.go    # Cron job handlers
```

## Purpose
Implements scheduled jobs that run at specific intervals or times. The Cron Handler:
- Executes periodic or scheduled tasks
- Delegates business logic to UseCase
- Implements proper error handling and logging
- Uses distributed tracing
- Prevents concurrent execution if needed
- Returns handler function compatible with cron library

## Rules Applied
- Layer boundary enforcement: Handler MUST NOT contain business logic
- No direct repository access (only through UseCase)
- Context enrichment with trace ID and structured logging
- Graceful error handling (don't crash the scheduler)
- Idempotency for safe retry/re-execution
- Job locking for single-instance execution (if required)

---

## Code Skeleton

### File: `module/<module>/delivery/cron/handler.go`

```go
package cron

import (
	"context"
	"time"

	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/helper"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// <Job>Handler executes the scheduled <job> task.
// This function is called by the cron scheduler at configured intervals.
func <Job>Handler(useCase domain.UseCase) func() {
	return func() {
		// Generate unique trace ID for this job execution
		traceId := uuid.New().String()
		ctx := helper.WithTraceID(context.Background(), traceId)
		
		// Setup structured logger
		logger := log.With().
			Str("job", "<job-name>").
			Str("trace_id", traceId).
			Time("started_at", time.Now()).
			Logger()
		ctx = logger.WithContext(ctx)

		logger.Info().Msg("cron job started")

		// Optional: Set timeout for job execution
		ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		// Execute business logic through use case
		if err := useCase.<JobMethod>(ctx); err != nil {
			logger.Error().
				Err(err).
				Msg("cron job failed")
			// Don't panic - let scheduler continue
			return
		}

		logger.Info().
			Dur("duration", time.Since(time.Now())).
			Msg("cron job completed successfully")
	}
}
```

### Alternative: Handler with Parameters

```go
package cron

import (
	"context"
	"time"

	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/helper"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// <Job>HandlerWithParams executes a parameterized scheduled task.
func <Job>HandlerWithParams(useCase domain.UseCase, param1 string, param2 int) func() {
	return func() {
		traceId := uuid.New().String()
		ctx := helper.WithTraceID(context.Background(), traceId)
		
		logger := log.With().
			Str("job", "<job-name>").
			Str("trace_id", traceId).
			Str("param1", param1).
			Int("param2", param2).
			Logger()
		ctx = logger.WithContext(ctx)

		logger.Info().Msg("cron job started")

		startTime := time.Now()
		ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
		defer cancel()

		// Pass parameters to use case
		if err := useCase.<JobMethod>(ctx, param1, param2); err != nil {
			logger.Error().
				Err(err).
				Dur("duration", time.Since(startTime)).
				Msg("cron job failed")
			return
		}

		logger.Info().
			Dur("duration", time.Since(startTime)).
			Msg("cron job completed successfully")
	}
}
```

### Alternative: Handler with Distributed Lock

For jobs that should run on only ONE instance in a cluster:

```go
package cron

import (
	"context"
	"time"

	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/helper"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// <Job>HandlerWithLock ensures only one instance executes the job at a time.
func <Job>HandlerWithLock(useCase domain.UseCase, redisClient *redis.Client) func() {
	const (
		lockKey    = "cron:lock:<job-name>"
		lockTTL    = 5 * time.Minute
		maxRetries = 3
	)

	return func() {
		traceId := uuid.New().String()
		ctx := helper.WithTraceID(context.Background(), traceId)
		
		logger := log.With().
			Str("job", "<job-name>").
			Str("trace_id", traceId).
			Logger()
		ctx = logger.WithContext(ctx)

		// Try to acquire distributed lock
		locked, err := redisClient.SetNX(ctx, lockKey, traceId, lockTTL).Result()
		if err != nil {
			logger.Error().Err(err).Msg("failed to acquire lock")
			return
		}

		if !locked {
			logger.Info().Msg("job already running on another instance, skipping")
			return
		}

		// Ensure lock is released
		defer func() {
			if err := redisClient.Del(context.Background(), lockKey).Err(); err != nil {
				logger.Error().Err(err).Msg("failed to release lock")
			}
		}()

		logger.Info().Msg("lock acquired, starting job")

		startTime := time.Now()
		ctx, cancel := context.WithTimeout(ctx, lockTTL-30*time.Second) // Leave buffer
		defer cancel()

		// Execute with retry logic
		var lastErr error
		for attempt := 1; attempt <= maxRetries; attempt++ {
			logger.Info().Int("attempt", attempt).Msg("executing job")

			if err := useCase.<JobMethod>(ctx); err != nil {
				lastErr = err
				logger.Warn().
					Err(err).
					Int("attempt", attempt).
					Msg("job attempt failed")

				if attempt < maxRetries {
					time.Sleep(time.Duration(attempt) * time.Second) // Exponential backoff
					continue
				}
			} else {
				logger.Info().
					Dur("duration", time.Since(startTime)).
					Int("attempts", attempt).
					Msg("job completed successfully")
				return
			}
		}

		logger.Error().
			Err(lastErr).
			Dur("duration", time.Since(startTime)).
			Msg("job failed after all retries")
	}
}
```

### Integration: Register Cron Jobs

```go
// In cmd/serve.go or main application bootstrap

package main

import (
	"github.com/<org>/<project>/module/<module>/delivery/cron"
	"github.com/robfig/cron/v3"
)

func registerCronJobs(scheduler *cron.Cron, <module>UseCase domain.UseCase) {
	// Every day at midnight
	scheduler.AddFunc("0 0 * * *", cron.<DailyJob>Handler(<module>UseCase))

	// Every 15 minutes
	scheduler.AddFunc("*/15 * * * *", cron.<FrequentJob>Handler(<module>UseCase))

	// Every hour
	scheduler.AddFunc("0 * * * *", cron.<HourlyJob>Handler(<module>UseCase))

	// Custom schedule with parameters
	scheduler.AddFunc("0 2 * * *", cron.<JobWithParams>Handler(<module>UseCase, "param", 123))
}
```

### Cron Schedule Examples

```
# Format: minute hour day month weekday

# Every minute
* * * * *

# Every 5 minutes
*/5 * * * *

# Every hour at minute 0
0 * * * *

# Every day at midnight
0 0 * * *

# Every day at 3:30 AM
30 3 * * *

# Every Monday at 9 AM
0 9 * * 1

# First day of every month at midnight
0 0 1 * *

# Weekdays at 6 PM
0 18 * * 1-5

# Every 30 minutes between 9 AM and 5 PM
*/30 9-17 * * *
```

---

## Tests Required

### File: `module/<module>/delivery/cron/handler_test.go`

1. **Test Handler Execution**:
   - Handler executes without panic
   - UseCase method is called
   - Context is properly configured

2. **Test Error Handling**:
   - UseCase error doesn't crash handler
   - Error is logged appropriately
   - Handler returns gracefully

3. **Test Context Timeout**:
   - Timeout is respected
   - Context cancellation works

4. **Test Distributed Lock** (if applicable):
   - Lock is acquired before execution
   - Lock prevents concurrent execution
   - Lock is released after completion
   - Lock is released on error

**Test Structure**:
```go
func Test<Job>Handler(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockUseCase)
		wantPanic bool
	}{
		{
			name: "successful execution",
			mockSetup: func(m *MockUseCase) {
				m.On("<JobMethod>", mock.Anything).Return(nil)
			},
			wantPanic: false,
		},
		{
			name: "usecase error",
			mockSetup: func(m *MockUseCase) {
				m.On("<JobMethod>", mock.Anything).Return(errors.New("error"))
			},
			wantPanic: false, // Should not panic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := new(MockUseCase)
			tt.mockSetup(mockUC)

			handler := <Job>Handler(mockUC)

			if tt.wantPanic {
				assert.Panics(t, handler)
			} else {
				assert.NotPanics(t, handler)
			}

			mockUC.AssertExpectations(t)
		})
	}
}
```

---

## Notes / Pitfalls

### Critical Points

1. **Error Handling**: NEVER panic in cron handler
   ```go
   // ✅ GOOD
   if err != nil {
       logger.Error().Err(err).Msg("job failed")
       return // Continue scheduling
   }
   
   // ❌ BAD
   if err != nil {
       panic(err) // Will crash scheduler
   }
   ```

2. **Trace ID Generation**:
   - Always generate unique ID for each execution
   - Use UUID or similar unique identifier
   - Include in all logs for correlation

3. **Context Timeout**:
   - Set appropriate timeout based on job duration
   - Prevent jobs from running indefinitely
   - Leave buffer time if using distributed lock

4. **Logging Best Practices**:
   ```go
   logger := log.With().
       Str("job", "job-name").
       Str("trace_id", traceId).
       Time("started_at", time.Now()).
       Logger()
   ```

5. **Idempotency**:
   - Jobs may execute multiple times (failures, restarts)
   - UseCase logic should handle duplicate execution safely
   - Consider checking last execution timestamp

6. **Distributed Lock**:
   - Use Redis or similar for cluster-wide coordination
   - Set appropriate TTL (longer than expected execution time)
   - Always release lock in defer
   - Handle lock acquisition failure gracefully

7. **Graceful Shutdown**:
   - Don't start new jobs during shutdown
   - Allow running jobs to complete
   - Set context cancellation

8. **Return Signature**:
   ```go
   // Handler must return func() for cron scheduler
   func HandlerName(deps...) func() {
       return func() {
           // Job logic
       }
   }
   ```

### Common Mistakes

- Panicking on errors (crashes entire scheduler)
- Not setting context timeout (jobs run indefinitely)
- Forgetting to generate trace ID (can't correlate logs)
- Not using distributed lock for singleton jobs (multiple instances run)
- Hardcoding schedule in handler (should be in registration)
- Not logging start/completion/duration
- Accessing shared state without synchronization
- Not handling context cancellation
- Using request-scoped context (use context.Background())

### Distributed Lock Patterns

**Redis SET NX with TTL**:
```go
locked, err := redisClient.SetNX(ctx, lockKey, value, ttl).Result()
```

**Redis with Lua Script** (atomic check-and-set):
```go
script := redis.NewScript(`
    if redis.call("get", KEYS[1]) == ARGV[1] then
        return redis.call("del", KEYS[1])
    else
        return 0
    end
`)
```

**Database Advisory Lock** (PostgreSQL):
```sql
SELECT pg_try_advisory_lock(hash_id);
-- Execute job
SELECT pg_advisory_unlock(hash_id);
```

### Scheduling Strategies

**Fixed Interval**: Run every N minutes
```go
scheduler.AddFunc("*/5 * * * *", handler)
```

**Fixed Time**: Run at specific time daily
```go
scheduler.AddFunc("0 2 * * *", handler) // 2 AM daily
```

**Delayed Start**: Wait before first execution
```go
time.Sleep(30 * time.Second)
scheduler.Start()
```

**One-Time Execution**:
```go
time.AfterFunc(duration, handler)
```

### Monitoring & Observability

**Metrics to Track**:
- Job execution count
- Job success/failure rate
- Job execution duration
- Lock acquisition failures
- Concurrent execution attempts (if prevented)

**Alerting Conditions**:
- Job hasn't run in expected interval
- Job failure rate exceeds threshold
- Job duration exceeds normal range
- Lock contention too high

**Logging Requirements**:
- Job start (with trace ID, timestamp)
- Job completion (with duration, result)
- Job failure (with error details, retry attempts)
- Lock status (acquired, skipped, released)

### UseCase Method Pattern

The use case should provide a dedicated method for the cron job:

```go
// In module/<module>/usecase/scheduled_task.go

func (u *UseCase) <JobMethod>(ctx context.Context) error {
	ctx, span := u.tracer.Start(ctx, "<JobMethod>")
	defer span.End()

	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("executing scheduled task")

	// Business logic
	// - Fetch records to process
	// - Apply transformations
	// - Update state
	// - Publish events if needed

	return nil
}
```

### Cron Library Options

**robfig/cron/v3** (most common):
```go
import "github.com/robfig/cron/v3"

c := cron.New()
c.AddFunc("0 0 * * *", handler)
c.Start()
defer c.Stop()
```

**Standard library** (time.Ticker):
```go
ticker := time.NewTicker(5 * time.Minute)
defer ticker.Stop()

for range ticker.C {
    handler()
}
```

**Alternatives**:
- `go-co-op/gocron` - More features, fluent API
- `beevik/cron` - Simple implementation
- Custom implementation with time.AfterFunc
