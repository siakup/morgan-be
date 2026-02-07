package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_sessions/domain"
)

type (
	UpdateShiftSessionRequest struct {
		Name   string `json:"name" validate:"required"`
		Start  string `json:"start" validate:"required"`
		End    string `json:"end" validate:"required"`
		Status bool   `json:"status"`
	}
)

// UpdateShiftSession handles PUT /shift-sessions/:id
func (h *ShiftSessionHandler) UpdateShiftSession(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("Invalid ID"))
	}

	userId, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || userId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	var req UpdateShiftSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.BadRequest("Invalid request body"))
	}

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
		Id:        id,
		Name:      req.Name,
		Start:     req.Start,
		End:       req.End,
		Status:    req.Status,
		UpdatedBy: &userId,
	}

	if err := h.useCase.Update(ctx, &shiftSession); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success[any](nil, "Shift session updated"))
}
