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

func (h *BuildingHandler) UpdateBuilding(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateBuildingRequest
	if err := validator.BindAndValidate(c, &req); err != nil {
		return h.handleError(c, err)
	}

	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	building := domain.Building{
		Id:        id,
		Name:      req.Name,
		Status:    req.Status,
		UpdatedBy: &userId,
	}

	err := h.useCase.Update(c.UserContext(), &building)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success(building, "Building updated"))
}
