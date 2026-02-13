package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/libraries/validator"
	"github.com/siakup/morgan-be/morgan/module/buildings/domain"
	"github.com/siakup/morgan-be/morgan/module/buildings/dto"
)

func (h *BuildingHandler) CreateBuilding(c *fiber.Ctx) error {
	var req dto.CreateBuildingRequest
	if err := validator.BindAndValidate(c, &req); err != nil {
		return h.handleError(c, err)
	}

	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	building := domain.Building{
		Name:      req.Name,
		Status:    req.Status,
		CreatedBy: &userId,
		UpdatedBy: &userId,
	}

	err := h.useCase.Create(c.UserContext(), &building)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.Success(building, "Building created"))
}
