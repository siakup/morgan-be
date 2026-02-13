package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/libraries/validator"
	"github.com/siakup/morgan-be/morgan/module/shift_groups/domain"
	"github.com/siakup/morgan-be/morgan/module/shift_groups/dto"
)

func (h *ShiftGroupHandler) UpdateShiftGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateShiftGroupRequest
	if err := validator.BindAndValidate(c, &req); err != nil {
		return h.handleError(c, err)
	}

	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	shiftGroup := domain.ShiftGroup{
		Id:        id,
		Name:      req.Name,
		Status:    req.Status,
		UpdatedBy: &userId,
	}

	err := h.useCase.Update(c.UserContext(), &shiftGroup)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.Success(shiftGroup, "Shift Group updated"))
}
