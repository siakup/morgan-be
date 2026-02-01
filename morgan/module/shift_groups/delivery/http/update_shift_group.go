package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

type UpdateShiftGroupRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}

func (h *ShiftGroupHandler) UpdateShiftGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateShiftGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Fail("BAD_REQUEST", "Invalid request body"))
	}

	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	shiftGroup := domain.ShiftGroup{
		Id:        id,
		Name:      req.Name,
		Status:    req.Status,
		UpdatedBy: &userId,
	}

	err := h.useCase.Update(c.Context(), &shiftGroup)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.Success(shiftGroup, "Shift Group updated"))
}
