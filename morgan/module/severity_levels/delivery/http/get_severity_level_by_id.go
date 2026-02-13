package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/responses"
)

func (h *SeverityLevelHandler) GetSeverityLevelByID(c *fiber.Ctx) error {
	id := c.Params("id")
	severityLevel, err := h.useCase.FindByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}
	if severityLevel == nil {
		return c.Status(http.StatusNotFound).JSON(responses.Fail("NOT_FOUND", "Severity Level not found"))
	}
	return c.JSON(responses.Success(severityLevel, "Severity Level retrieved successfully"))
}
