package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/libraries/validator"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/dto"
)

func (h *SeverityLevelHandler) UpdateSeverityLevel(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateSeverityLevelRequest
	if err := validator.BindAndValidate(c, &req); err != nil {
		return h.handleError(c, err)
	}

	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	severityLevel := domain.SeverityLevel{
		Id:        id,
		Name:      req.Name,
		Status:    req.Status,
		UpdatedBy: &userId,
	}

	err := h.useCase.Update(c.UserContext(), &severityLevel)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.Success(severityLevel, "Severity Level updated"))
}
