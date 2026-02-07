package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
)

type (
	GetShiftSessionByIDResponse struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Start  string `json:"start"`
		End    string `json:"end"`
		Status bool   `json:"status"`
	}
)

// GetShiftSessionByID handles GET /shift-sessions/:id
func (h *ShiftSessionHandler) GetShiftSessionByID(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("Invalid ID"))
	}

	shiftSession, err := h.useCase.Get(ctx, id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success(GetShiftSessionByIDResponse{
		Id:     shiftSession.Id,
		Name:   shiftSession.Name,
		Start:  shiftSession.Start,
		End:    shiftSession.End,
		Status: shiftSession.Status,
	}, "Shift session retrieved"))
}
