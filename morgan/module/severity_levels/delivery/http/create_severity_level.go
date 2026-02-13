package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
)

type CreateSeverityLevelRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}

func (h *SeverityLevelHandler) CreateSeverityLevel(c *fiber.Ctx) error {
	var req CreateSeverityLevelRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Fail("BAD_REQUEST", "Invalid request body"))
	}

	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	severityLevel := domain.SeverityLevel{
		Name:      req.Name,
		Status:    req.Status,
		CreatedBy: &userId,
		UpdatedBy: &userId,
	}

	err := h.useCase.Create(c.Context(), &severityLevel)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.Success(severityLevel, "Severity Level created"))
}
