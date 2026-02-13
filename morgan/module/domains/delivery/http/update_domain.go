package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/libraries/validation"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/morgan/module/domains/domain"
)

type (
	UpdateDomainRequest struct {
		Name   string `json:"name" validate:"required"`
		Status bool   `json:"status"`
	}
)

func (h *DomainHandler) UpdateDomain(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("Invalid ID"))
	}

	userId, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || userId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	// DTO parsed/validated by middleware
	raw := c.Locals(validation.ValidatedBodyKey)
	req, ok := raw.(*UpdateDomainRequest)
	if !ok || req == nil {
		return h.handleError(c, errors.BadRequest("invalid request body"))
	}

	updatedDomain := domain.Domain{
		Id:        id,
		Name:      req.Name,
		Status:    req.Status,
		UpdatedBy: &userId,
	}

	if err := h.useCase.Update(ctx, &updatedDomain); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success[any](nil, "Domain updated"))
}
