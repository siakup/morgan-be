package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/responses"
)

func (h *DivisionHandler) GetDivisionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	division, err := h.useCase.Get(c.UserContext(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.Success(division, "Division retrieved successfully"))
}
