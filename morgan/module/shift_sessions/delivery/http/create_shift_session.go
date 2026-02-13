package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/domain"
)

type (
	CreateShiftSessionRequest struct {
		Name  string `json:"name" validate:"required"`
		Start string `json:"start" validate:"required"`
		End   string `json:"end" validate:"required"`
	}
	CreateShiftSessionResponse struct {
		Id string `json:"id"`
	}
)

// CreateShiftSession handles POST /shift-sessions
func (h *ShiftSessionHandler) CreateShiftSession(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userId, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || userId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	var req CreateShiftSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.BadRequest("Invalid request body"))
	}

	// Manual basic validation
	if req.Name == "" {
		return h.handleError(c, errors.BadRequest("field name is required"))
	}
	if req.Start == "" {
		return h.handleError(c, errors.BadRequest("field start is required"))
	}
	if req.End == "" {
		return h.handleError(c, errors.BadRequest("field end is required"))
	}

	shiftSession := domain.ShiftSession{
		Name:      req.Name,
		Start:     req.Start,
		End:       req.End,
		CreatedBy: &userId,
		UpdatedBy: &userId,
	}

	if err := h.useCase.Create(ctx, &shiftSession); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.Success(CreateShiftSessionResponse{
		Id: shiftSession.Id,
	}, "Shift session created"))
}
