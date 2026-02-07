package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
)

// GetRoomTypeByIDResponse is the response for a single room type.
type GetRoomTypeByIDResponse struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsActive    bool    `json:"is_active"`
	CreatedBy   *string `json:"created_by,omitempty"`
	UpdatedBy   *string `json:"updated_by,omitempty"`
}

// GetRoomTypeByID handles GET /room-types/:id.
func (h *RoomTypeHandler) GetRoomTypeByID(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("Invalid ID"))
	}

	rt, err := h.useCase.FindByID(ctx, id)
	if err != nil {
		return h.handleError(c, err)
	}
	if rt == nil {
		return c.Status(http.StatusNotFound).JSON(responses.Fail("NOT_FOUND", "Room type not found"))
	}

	return c.Status(http.StatusOK).JSON(responses.Success(GetRoomTypeByIDResponse{
		Id:          rt.Id,
		Name:        rt.Name,
		Description: rt.Description,
		IsActive:    rt.IsActive,
		CreatedBy:   rt.CreatedBy,
		UpdatedBy:   rt.UpdatedBy,
	}, "Room type retrieved"))
}
