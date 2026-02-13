package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/libraries/types"
	"github.com/siakup/morgan-be/morgan/module/shift_groups/domain"
)

func (h *ShiftGroupHandler) GetShiftGroups(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	size := c.QueryInt("size", 10)

	filter := domain.ShiftGroupFilter{
		Pagination: types.Pagination{
			Page: page,
			Size: size,
		},
		Search: c.Query("search"),
	}

	shiftGroups, total, err := h.useCase.FindAll(c.UserContext(), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	meta := &responses.Meta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: (int(total) + size - 1) / size,
	}

	return c.JSON(responses.SuccessWithMeta(shiftGroups, "Shift Groups retrieved", meta))
}
