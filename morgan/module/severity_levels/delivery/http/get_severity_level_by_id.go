package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/responses"
)

func (h *SeverityLevelHandler) GetSeverityLevelByID(c *fiber.Ctx) error {
	id := c.Params("id")
	severityLevel, err := h.useCase.Get(c.UserContext(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.Success(severityLevel, "Severity Level retrieved successfully"))
}
