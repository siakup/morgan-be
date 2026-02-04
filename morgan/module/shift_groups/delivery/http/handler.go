package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

// ShiftGroupHandler handles HTTP requests for shift groups module.
type ShiftGroupHandler struct {
	useCase domain.UseCase
	auth    *middleware.AuthorizationMiddleware
}

// NewShiftGroupHandler creates a new ShiftGroupHandler.
func NewShiftGroupHandler(useCase domain.UseCase, auth *middleware.AuthorizationMiddleware) *ShiftGroupHandler {
	return &ShiftGroupHandler{
		useCase: useCase,
		auth:    auth,
	}
}

// RegisterRoutes registers the routes for the shift groups module.
func (h *ShiftGroupHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/shift-groups", middleware.TraceMiddleware)

	// Adjust permissions as needed
	group.Get("/", h.auth.Authenticate("shift_groups.view"), h.GetShiftGroups)
	group.Get("/:id", h.auth.Authenticate("shift_groups.view"), h.GetShiftGroupByID)
	group.Post("/", h.auth.Authenticate("shift_groups.edit"), h.CreateShiftGroup)
	group.Put("/:id", h.auth.Authenticate("shift_groups.edit"), h.UpdateShiftGroup)
	group.Delete("/:id", h.auth.Authenticate("shift_groups.edit"), h.DeleteShiftGroup)
}

// handleError handles errors by mapping them to standardized responses.
func (h *ShiftGroupHandler) handleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(responses.Fail(string(appErr.Type), appErr.Message))
	}

	return c.Status(http.StatusInternalServerError).JSON(responses.Fail("SYSTEM_ERROR", err.Error()))
}
