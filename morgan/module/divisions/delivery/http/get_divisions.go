package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/morgan/module/divisions/domain"
)

func (h *DivisionHandler) GetDivisions(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	search := c.Query("search", "")

	filter := domain.DivisionFilter{
		Limit:  limit,
		Offset: offset,
		Search: search,
	}

	divisions, count, err := h.useCase.GetAll(c.UserContext(), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	meta := &responses.Meta{
		Page:       offset/limit + 1,
		Size:       limit,
		Total:      count,
		TotalPages: int((count + int64(limit) - 1) / int64(limit)),
	}

	return c.JSON(responses.SuccessWithMeta(divisions, "Divisions retrieved successfully", meta))
}
