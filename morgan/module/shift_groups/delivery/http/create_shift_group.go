package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/morgan/module/shift_groups/domain"
)

type CreateShiftGroupRequest struct {
	Name   string `json:"name" validate:"required"`
	Status bool   `json:"status"`
}

func (h *ShiftGroupHandler) CreateShiftGroup(c *fiber.Ctx) error {
	var req CreateShiftGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Fail("BAD_REQUEST", "Invalid request body"))
	}

	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	shiftGroup := domain.ShiftGroup{
		Name:      req.Name,
		Status:    req.Status,
		CreatedBy: &userId,
		UpdatedBy: &userId,
	}

	err := h.useCase.Create(c.Context(), &shiftGroup)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.Success(shiftGroup, "Shift Group created"))
}
