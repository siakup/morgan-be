package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
)

type (
	GetDomainByIDResponse struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Status bool   `json:"status"`
	}
)

// GetDomainByID handles GET /domains/:id
func (h *DomainHandler) GetDomainByID(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("Invalid ID"))
	}

	domain, err := h.useCase.Get(ctx, id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success(GetDomainByIDResponse{
		Id:     domain.Id,
		Name:   domain.Name,
		Status: domain.Status,
	}, "Domain retrieved"))
}
