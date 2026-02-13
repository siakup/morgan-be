package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
)

func (h *SeverityLevelHandler) DeleteSeverityLevel(c *fiber.Ctx) error {
	id := c.Params("id")
	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	err := h.useCase.Delete(c.UserContext(), id, userId)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.Success[any](nil, "Severity Level deleted"))
}
