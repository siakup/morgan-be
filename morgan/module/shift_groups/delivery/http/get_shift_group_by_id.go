package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/responses"
)

func (h *ShiftGroupHandler) GetShiftGroupByID(c *fiber.Ctx) error {
	id := c.Params("id")
	shiftGroup, err := h.useCase.Get(c.UserContext(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.Success(shiftGroup, "Shift Group retrieved successfully"))
}
