package http

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/libraries/types"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/domain"
)

type (
	GetShiftSessionsResponse struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Start  string `json:"start"`
		End    string `json:"end"`
		Status bool   `json:"status"`
	}
)

// GetShiftSessions handles GET /shift-sessions
func (h *ShiftSessionHandler) GetShiftSessions(c *fiber.Ctx) error {
	ctx := c.UserContext()

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	filter := domain.ShiftSessionFilter{
		Pagination: types.Pagination{
			Page: page,
			Size: pageSize,
		},
		Search: c.Query("search"),
	}

	shiftSessions, total, err := h.useCase.FindAll(ctx, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	result := make([]GetShiftSessionsResponse, len(shiftSessions))
	for i, ss := range shiftSessions {
		result[i] = GetShiftSessionsResponse{
			Id:     ss.Id,
			Name:   ss.Name,
			Start:  ss.Start,
			End:    ss.End,
			Status: ss.Status,
		}
	}

	meta := &responses.Meta{
		Page:       page,
		Size:       pageSize,
		Total:      total,
		TotalPages: (int(total) + pageSize - 1) / pageSize,
	}

	return c.Status(http.StatusOK).JSON(responses.SuccessWithMeta(result, "Shift sessions retrieved", meta))
}
