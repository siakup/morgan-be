package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
)

func (h *ShiftGroupHandler) GetShiftGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	shiftGroup, err := h.useCase.FindByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}
	if shiftGroup == nil {
		return c.Status(http.StatusNotFound).JSON(responses.Fail("NOT_FOUND", "Shift Group not found"))
	}

	return c.JSON(responses.Success(shiftGroup, "Shift Group retrieved"))
}
