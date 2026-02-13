package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/libraries/validator"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/dto"
)

func (h *SeverityLevelHandler) CreateSeverityLevel(c *fiber.Ctx) error {
	var req dto.CreateSeverityLevelRequest
	if err := validator.BindAndValidate(c, &req); err != nil {
		return h.handleError(c, err)
	}

	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	severityLevel := domain.SeverityLevel{
		Name:      req.Name,
		Status:    req.Status,
		CreatedBy: &userId,
		UpdatedBy: &userId,
	}

	err := h.useCase.Create(c.UserContext(), &severityLevel)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.Success(severityLevel, "Severity Level created"))
}
