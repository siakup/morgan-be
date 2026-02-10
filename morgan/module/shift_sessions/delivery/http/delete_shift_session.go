package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
)

// DeleteShiftSession handles DELETE /shift-sessions/:id
func (h *ShiftSessionHandler) DeleteShiftSession(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("Invalid ID"))
	}

	deletedBy, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || deletedBy == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	if err := h.useCase.Delete(ctx, id, deletedBy); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success[any](nil, "Shift session deleted"))
}
