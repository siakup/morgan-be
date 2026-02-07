package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

// UpdateRoomTypeRequest is the request body for updating a room type.
type UpdateRoomTypeRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// UpdateRoomType handles PUT /room-types/:id.
func (h *RoomTypeHandler) UpdateRoomType(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("Invalid ID"))
	}

	userId, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || userId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	var req UpdateRoomTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.BadRequest("Invalid request body"))
	}

	if req.Name == "" {
		return h.handleError(c, errors.BadRequest("field name is required"))
	}

	rt := domain.RoomType{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
		UpdatedBy:   &userId,
	}

	if err := h.useCase.Update(ctx, &rt); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success(rt, "Room type updated"))
}
