package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/libraries/validator"
	"github.com/siakup/morgan-be/morgan/module/shift_groups/domain"
	"github.com/siakup/morgan-be/morgan/module/shift_groups/dto"
)

func (h *ShiftGroupHandler) CreateShiftGroup(c *fiber.Ctx) error {
	var req dto.CreateShiftGroupRequest
	if err := validator.BindAndValidate(c, &req); err != nil {
		return h.handleError(c, err)
	}

	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	shiftGroup := domain.ShiftGroup{
		Name:      req.Name,
		Status:    req.Status,
		CreatedBy: &userId,
		UpdatedBy: &userId,
	}

	err := h.useCase.Create(c.UserContext(), &shiftGroup)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.Success(shiftGroup, "Shift Group created"))
}
