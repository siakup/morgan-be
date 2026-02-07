package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

// CreateRoomTypeRequest is the request body for creating a room type.
type CreateRoomTypeRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// CreateRoomTypeResponse is the response body after creating a room type.
type CreateRoomTypeResponse struct {
	Id string `json:"id"`
}

// CreateRoomType handles POST /room-types.
func (h *RoomTypeHandler) CreateRoomType(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userId, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || userId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	var req CreateRoomTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.BadRequest("Invalid request body"))
	}

	if req.Name == "" {
		return h.handleError(c, errors.BadRequest("field name is required"))
	}

	rt := domain.RoomType{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
		CreatedBy:   &userId,
		UpdatedBy:   &userId,
	}

	if err := h.useCase.Create(ctx, &rt); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.Success(CreateRoomTypeResponse{Id: rt.Id}, "Room type created"))
}
