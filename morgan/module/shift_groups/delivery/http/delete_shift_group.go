package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
)

func (h *ShiftGroupHandler) DeleteShiftGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	err := h.useCase.Delete(c.Context(), id, userId)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success[any](nil, "Shift Group deleted"))
}
