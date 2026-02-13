package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
)

// SeverityLevelHandler handles HTTP requests for severity levels module.
type SeverityLevelHandler struct {
	useCase domain.UseCase
	auth    *middleware.AuthorizationMiddleware
}

// NewSeverityLevelHandler creates a new SeverityLevelHandler.
func NewSeverityLevelHandler(useCase domain.UseCase, auth *middleware.AuthorizationMiddleware) *SeverityLevelHandler {
	return &SeverityLevelHandler{
		useCase: useCase,
		auth:    auth,
	}
}

// RegisterRoutes registers the routes for the severity levels module.
func (h *SeverityLevelHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/severity-levels", middleware.TraceMiddleware)

	// Adjust permissions as needed
	group.Get("/", h.auth.Authenticate("severity_levels.view"), h.GetSeverityLevels)
	group.Get("/:id", h.auth.Authenticate("severity_levels.view"), h.GetSeverityLevelByID)
	group.Post("/", h.auth.Authenticate("severity_levels.edit"), h.CreateSeverityLevel)
	group.Put("/:id", h.auth.Authenticate("severity_levels.edit"), h.UpdateSeverityLevel)
	group.Delete("/:id", h.auth.Authenticate("severity_levels.edit"), h.DeleteSeverityLevel)
}

// handleError handles errors by mapping them to standardized responses.
func (h *SeverityLevelHandler) handleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(responses.Fail(string(appErr.Type), appErr.Message))
	}

	return c.Status(http.StatusInternalServerError).JSON(responses.Fail("SYSTEM_ERROR", err.Error()))
}
