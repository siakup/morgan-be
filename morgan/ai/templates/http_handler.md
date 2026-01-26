# Template: HTTP Handler (Delivery Layer)

## Location
```
module/<module>/delivery/http/
├── handler.go              # Handler struct, constructor, route registration
├── <action1>_<entity>.go  # Handler for action 1 (e.g., create_user.go)
├── <action2>_<entity>.go  # Handler for action 2 (e.g., get_user_by_id.go)
├── <action3>_<entity>.go  # Handler for action 3 (e.g., update_user.go)
└── ...                     # One file per endpoint
```

## Purpose
Implements the HTTP delivery layer using Fiber framework. The HTTP Handler:
- Translates HTTP requests to domain operations
- Delegates business logic to UseCase
- Transforms domain errors to HTTP responses
- Handles request validation and response formatting
- Manages request context (trace ID, logging)
- Registers routes with appropriate middleware

## Rules Applied
- Layer boundary enforcement: Handler MUST NOT contain business logic
- No direct repository access (only through UseCase)
- Use standardized response format (responses.Success/Fail)
- Error mapping: AppError → HTTP status codes
- Context enrichment: request ID, structured logging
- Request/Response DTOs separate from domain objects

---

## Code Skeleton

### File: `module/<module>/delivery/http/handler.go`

```go
package http

import (
	"net/http"

	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/helper"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	gofiber "github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// <Entity>Handler handles HTTP requests for <module> module.
type <Entity>Handler struct {
	useCase domain.UseCase
}

// New<Entity>Handler creates a new <Entity>Handler.
func New<Entity>Handler(useCase domain.UseCase) *<Entity>Handler {
	return &<Entity>Handler{
		useCase: useCase,
	}
}

// RegisterRoutes registers the routes for the <module> module.
// It implements the fiber.Router interface.
func (h *<Entity>Handler) RegisterRoutes(app *gofiber.App) {
	group := app.Group("/<entities>", func(c *gofiber.Ctx) error {
		rid, ok := c.Locals("requestid").(string)
		if !ok || rid == "" {
			rid = c.Get(gofiber.HeaderXRequestID)
		}

		ctx := helper.WithTraceID(c.UserContext(), rid)
		logger := log.With().Str("request_id", rid).Logger()
		ctx = logger.WithContext(ctx)

		c.SetUserContext(ctx)
		return c.Next()
	})

	// CRUD routes
	group.Get("/", h.Get<Entities>)
	group.Get("/:id", h.Get<Entity>ByID)
	group.Post("/", h.Create<Entity>)
	group.Put("/:id", h.Update<Entity>)
	group.Delete("/:id", h.Delete<Entity>)
	
	// Add custom routes as needed
	// group.Post("/:id/<action>", h.<CustomAction>)
}

// handleError handles errors by mapping them to standardized responses.
func (h *<Entity>Handler) handleError(c *gofiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(responses.Fail(string(appErr.Type), appErr.Message))
	}

	return c.Status(http.StatusInternalServerError).JSON(responses.Fail("SYSTEM_ERROR", err.Error()))
}
```

### File: `module/<module>/delivery/http/create_<entity>.go`

```go
package http

import (
	"net/http"

	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	gofiber "github.com/gofiber/fiber/v2"
)

type (
	Create<Entity>Request struct {
		Id      string  ` + "`json:\"id\" validate:\"required\"`" + `
		<Field1> <Type1> ` + "`json:\"<field1>\" validate:\"required\"`" + `
		<Field2> <Type2> ` + "`json:\"<field2>\"`" + `
		// Add all fields needed for creation
	}
	Create<Entity>Response struct {
		Id      string  ` + "`json:\"id\"`" + `
		<Field1> <Type1> ` + "`json:\"<field1>\"`" + `
		<Field2> <Type2> ` + "`json:\"<field2>\"`" + `
		// Match domain object or customize as needed
	}
)

// Create<Entity> handles POST /<entities>
func (h *<Entity>Handler) Create<Entity>(c *gofiber.Ctx) error {
	ctx := c.UserContext()

	var req Create<Entity>Request
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.BadRequest("Invalid request body"))
	}

	// Optional: Add request validation
	// if err := validate.Struct(req); err != nil {
	//     return h.handleError(c, errors.BadRequest(err.Error()))
	// }

	// Map request DTO to domain object
	<entity> := domain.<Entity>{
		Id:      req.Id,
		<Field1>: req.<Field1>,
		<Field2>: req.<Field2>,
	}

	// Delegate to use case
	if err := h.useCase.Create(ctx, &<entity>); err != nil {
		return h.handleError(c, err)
	}

	// Return success response
	return c.Status(http.StatusCreated).JSON(responses.Success(Create<Entity>Response{
		Id:      <entity>.Id,
		<Field1>: <entity>.<Field1>,
		<Field2>: <entity>.<Field2>,
	}, "<Entity> created"))
}
```

### File: `module/<module>/delivery/http/get_<entity>_by_id.go`

```go
package http

import (
	"net/http"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	gofiber "github.com/gofiber/fiber/v2"
)

type (
	Get<Entity>ByIDResponse struct {
		Id      string  ` + "`json:\"id\"`" + `
		<Field1> <Type1> ` + "`json:\"<field1>\"`" + `
		<Field2> <Type2> ` + "`json:\"<field2>\"`" + `
	}
)

// Get<Entity>ByID handles GET /<entities>/:id
func (h *<Entity>Handler) Get<Entity>ByID(c *gofiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("ID is required"))
	}

	<entity>, err := h.useCase.Get(ctx, id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success(Get<Entity>ByIDResponse{
		Id:      <entity>.Id,
		<Field1>: <entity>.<Field1>,
		<Field2>: <entity>.<Field2>,
	}, "<Entity> retrieved"))
}
```

### File: `module/<module>/delivery/http/get_<entities>.go`

```go
package http

import (
	"net/http"
	"strconv"

	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/types"
	gofiber "github.com/gofiber/fiber/v2"
)

type (
	Get<Entities>Response struct {
		Id      string  ` + "`json:\"id\"`" + `
		<Field1> <Type1> ` + "`json:\"<field1>\"`" + `
		<Field2> <Type2> ` + "`json:\"<field2>\"`" + `
	}
)

// Get<Entities> handles GET /<entities>
func (h *<Entity>Handler) Get<Entities>(c *gofiber.Ctx) error {
	ctx := c.UserContext()

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	// Build filter from query parameters
	filter := domain.<Entity>Filter{
		Pagination: types.Pagination{
			Page:     page,
			PageSize: pageSize,
			Offset:   (page - 1) * pageSize,
		},
		// Add custom filters
		// <FilterField1>: parseFilterField(c.Query("<filter_field1>")),
	}

	<entities>, err := h.useCase.FindAll(ctx, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	// Map domain objects to response DTOs
	result := make([]Get<Entities>Response, len(<entities>))
	for i, <entity> := range <entities> {
		result[i] = Get<Entities>Response{
			Id:      <entity>.Id,
			<Field1>: <entity>.<Field1>,
			<Field2>: <entity>.<Field2>,
		}
	}

	return c.Status(http.StatusOK).JSON(responses.Success(result, "<Entities> retrieved"))
}
```

### File: `module/<module>/delivery/http/update_<entity>.go`

```go
package http

import (
	"net/http"

	"github.com/<org>/<project>/module/<module>/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	gofiber "github.com/gofiber/fiber/v2"
)

type (
	Update<Entity>Request struct {
		<Field1> <Type1> ` + "`json:\"<field1>\" validate:\"required\"`" + `
		<Field2> <Type2> ` + "`json:\"<field2>\"`" + `
		// Fields allowed to be updated
	}
	Update<Entity>Response struct {
		Id      string  ` + "`json:\"id\"`" + `
		<Field1> <Type1> ` + "`json:\"<field1>\"`" + `
		<Field2> <Type2> ` + "`json:\"<field2>\"`" + `
	}
)

// Update<Entity> handles PUT /<entities>/:id
func (h *<Entity>Handler) Update<Entity>(c *gofiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("ID is required"))
	}

	var req Update<Entity>Request
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.BadRequest("Invalid request body"))
	}

	<entity> := domain.<Entity>{
		Id:      id,
		<Field1>: req.<Field1>,
		<Field2>: req.<Field2>,
	}

	if err := h.useCase.Put(ctx, &<entity>); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success(Update<Entity>Response{
		Id:      <entity>.Id,
		<Field1>: <entity>.<Field1>,
		<Field2>: <entity>.<Field2>,
	}, "<Entity> updated"))
}
```

### File: `module/<module>/delivery/http/delete_<entity>.go`

```go
package http

import (
	"net/http"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	gofiber "github.com/gofiber/fiber/v2"
)

// Delete<Entity> handles DELETE /<entities>/:id
func (h *<Entity>Handler) Delete<Entity>(c *gofiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("ID is required"))
	}

	if err := h.useCase.Delete(ctx, id); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success(nil, "<Entity> deleted"))
}
```

---

## Tests Required

### File: `module/<module>/delivery/http/handler_test.go`

1. **Test Create<Entity>**:
   - Valid request → 201 Created
   - Invalid JSON → 400 Bad Request
   - Missing required fields → 400 Bad Request
   - UseCase error → Appropriate error response

2. **Test Get<Entity>ByID**:
   - Valid ID → 200 OK with data
   - Missing ID → 400 Bad Request
   - Non-existent ID → 404 Not Found
   - UseCase error → 500 Internal Server Error

3. **Test Get<Entities>**:
   - No filters → 200 OK with list
   - With pagination → Correct offset/limit
   - Empty result → 200 OK with empty array
   - UseCase error → 500 Internal Server Error

4. **Test Update<Entity>**:
   - Valid update → 200 OK
   - Invalid request body → 400 Bad Request
   - Non-existent entity → 404 Not Found

5. **Test Delete<Entity>**:
   - Valid deletion → 200 OK
   - Non-existent entity → 404 Not Found

**Test Structure**:
```go
func TestHandler_Create<Entity>(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*MockUseCase)
		expectedStatus int
		expectedBody   string
	}{
		// Test cases
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Fiber app
			app := fiber.New()
			
			// Setup mock usecase
			mockUC := new(MockUseCase)
			tt.mockSetup(mockUC)
			
			// Create handler and register routes
			handler := NewHandler(mockUC)
			handler.RegisterRoutes(app)
			
			// Create test request
			req := httptest.NewRequest("POST", "/<entities>", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			
			// Execute request
			resp, _ := app.Test(req)
			
			// Assert
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
```

---

## Notes / Pitfalls

### Critical Points

1. **Layer Boundary**: Handler MUST NOT contain business logic
   - ❌ BAD: Validation logic beyond request parsing
   - ❌ BAD: Direct repository calls
   - ❌ BAD: Data transformation/computation
   - ✅ GOOD: Parse request → Call usecase → Format response

2. **Request/Response DTOs**:
   - Always separate from domain objects
   - Use JSON tags for serialization
   - Include validation tags if using validator
   - Can be different from domain model (subset, superset, transformed)

3. **Error Handling Pattern**:
   ```go
   // ✅ GOOD - Use centralized error handler
   if err := h.useCase.Create(ctx, entity); err != nil {
       return h.handleError(c, err)
   }
   
   // ❌ BAD - Manual error mapping per endpoint
   if err != nil {
       return c.Status(500).JSON(...)
   }
   ```

4. **Context Enrichment**:
   - Extract request ID from header or locals
   - Create structured logger with request ID
   - Attach to context for downstream propagation

5. **Route Registration**:
   - Use route groups for logical organization
   - Apply middleware at group level (auth, logging, etc.)
   - Keep route paths RESTful

6. **Response Format**: Use standardized response wrapper
   ```go
   responses.Success(data, message)
   responses.Fail(errorType, message)
   ```

7. **HTTP Status Codes**:
   - 200 OK: Successful GET, PUT, DELETE
   - 201 Created: Successful POST
   - 400 Bad Request: Invalid input, validation failure
   - 404 Not Found: Resource doesn't exist
   - 500 Internal Server Error: System errors

8. **Pagination**:
   - Default page size (e.g., 10)
   - Maximum page size limit (e.g., 100)
   - Calculate offset: `(page - 1) * pageSize`

### Common Mistakes

- Including business logic in handler (belongs in usecase)
- Returning domain objects directly (should use DTOs)
- Not validating path parameters (`:id`)
- Ignoring context from `c.UserContext()`
- Hardcoding error messages instead of using errors library
- Not handling missing required fields
- Direct repository injection (should use usecase only)
- Using generic `error` interface instead of `*errors.AppError`
- Not setting appropriate HTTP status codes
- Forgetting to propagate request ID to context

### Security Considerations

1. **Input Validation**: Always validate and sanitize input
2. **Parameter Extraction**: Validate path/query parameters before use
3. **Authentication/Authorization**: Apply middleware at group level
4. **Rate Limiting**: Consider adding rate limiting middleware
5. **Request Size Limits**: Configure max body size in Fiber
6. **CORS**: Configure CORS if needed for cross-origin requests

### Fiber-Specific Notes

- `c.Params("id")` - Extract path parameter
- `c.Query("key", "default")` - Extract query parameter
- `c.BodyParser(&struct)` - Parse JSON request body
- `c.Status(code).JSON(data)` - Set status and return JSON
- `c.UserContext()` - Get request context (important for tracing)
- `c.Locals("key")` - Get middleware-set values
- `c.Get(header)` - Get request header
