package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/libraries/validation"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
	"github.com/siakup/morgan-be/libraries/errors"
)

type CreateSeverityLevelRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}

func (h *SeverityLevelHandler) CreateSeverityLevel(c *fiber.Ctx) error {
	raw := c.Locals(validation.ValidatedBodyKey)
	req, ok := raw.(*CreateSeverityLevelRequest)
	if !ok || req == nil {
		return h.handleError(c, errors.BadRequest("invalid request body"))
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
