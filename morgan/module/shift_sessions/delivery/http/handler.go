package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/libraries/validation"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/domain"
)

// ShiftSessionHandler handles HTTP requests for shift sessions module.
type ShiftSessionHandler struct {
	useCase domain.UseCase
	auth    *middleware.AuthorizationMiddleware
}

// NewShiftSessionHandler creates a new ShiftSessionHandler.
func NewShiftSessionHandler(useCase domain.UseCase, auth *middleware.AuthorizationMiddleware) *ShiftSessionHandler {
	return &ShiftSessionHandler{
		useCase: useCase,
		auth:    auth,
	}
}

// RegisterRoutes registers the routes for the shift sessions module.
func (h *ShiftSessionHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/shift-sessions", middleware.TraceMiddleware)

	group.Get("/", h.auth.Authenticate("shift_sessions.schedule.shift_sessions.view"), h.GetShiftSessions)
	group.Get("/:id", h.auth.Authenticate("shift_sessions.schedule.shift_sessions.view"), h.GetShiftSessionByID)
	group.Post("/", h.auth.Authenticate("shift_sessions.schedule.shift_sessions.create"), validation.ValidateBody(func() interface{} { return &CreateShiftSessionRequest{} }), h.CreateShiftSession)
	group.Put("/:id", h.auth.Authenticate("shift_sessions.schedule.shift_sessions.edit"), validation.ValidateBody(func() interface{} { return &UpdateShiftSessionRequest{} }), h.UpdateShiftSession)
	group.Delete("/:id", h.auth.Authenticate("shift_sessions.schedule.shift_sessions.delete"), h.DeleteShiftSession)
}

// handleError handles errors by mapping them to standardized responses.
func (h *ShiftSessionHandler) handleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(responses.Fail(string(appErr.Type), appErr.Message))
	}

	return c.Status(http.StatusInternalServerError).JSON(responses.Fail("SYSTEM_ERROR", err.Error()))
}
