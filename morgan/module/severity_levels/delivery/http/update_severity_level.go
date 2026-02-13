package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/libraries/validation"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
	"github.com/siakup/morgan-be/libraries/errors"
)

type UpdateSeverityLevelRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}

func (h *SeverityLevelHandler) UpdateSeverityLevel(c *fiber.Ctx) error {
	id := c.Params("id")
	raw := c.Locals(validation.ValidatedBodyKey)
	req, ok := raw.(*UpdateSeverityLevelRequest)
	if !ok || req == nil {
		return h.handleError(c, errors.BadRequest("invalid request body"))
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
