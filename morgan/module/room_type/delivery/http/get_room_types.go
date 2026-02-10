package http

import (
	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/types"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

// GetRoomTypesResponseItem is a single item in the list response.
type GetRoomTypesResponseItem struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsActive    bool    `json:"is_active"`
	CreatedBy   *string `json:"created_by,omitempty"`
	UpdatedBy   *string `json:"updated_by,omitempty"`
}

// GetRoomTypes handles GET /room-types.
func (h *RoomTypeHandler) GetRoomTypes(c *fiber.Ctx) error {
	ctx := c.UserContext()

	page := c.QueryInt("page", 1)
	size := c.QueryInt("size", 10)

	filter := domain.RoomTypeFilter{
		Pagination: types.Pagination{
			Page: page,
			Size: size,
		},
		Search: c.Query("search"),
	}

	list, total, err := h.useCase.FindAll(ctx, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	result := make([]GetRoomTypesResponseItem, len(list))
	for i, rt := range list {
		result[i] = GetRoomTypesResponseItem{
			Id:          rt.Id,
			Name:        rt.Name,
			Description: rt.Description,
			IsActive:    rt.IsActive,
			CreatedBy:   rt.CreatedBy,
			UpdatedBy:   rt.UpdatedBy,
		}
	}

	meta := &responses.Meta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: (int(total) + size - 1) / size,
	}

	return c.JSON(responses.SuccessWithMeta(result, "Room types retrieved", meta))
}
