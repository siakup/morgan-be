package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
)

type UpdateSeverityLevelRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}

func (h *SeverityLevelHandler) UpdateSeverityLevel(c *fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateSeverityLevelRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Fail("BAD_REQUEST", "Invalid request body"))
	}

	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	severityLevel := domain.SeverityLevel{
		Id:        id,
		Name:      req.Name,
		Status:    req.Status,
		UpdatedBy: &userId,
	}

	err := h.useCase.Update(c.Context(), &severityLevel)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.Success[any](nil, "Severity Level updated"))
}
