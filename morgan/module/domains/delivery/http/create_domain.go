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
	CreateDomainRequest struct {
		Name string `json:"name" validate:"required"`
	}
	CreateDomainResponse struct {
		Id string `json:"id"`
	}
)

func (h *DomainHandler) CreateDomain(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userId, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || userId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	raw := c.Locals(validation.ValidatedBodyKey)
	req, ok := raw.(*CreateDomainRequest)
	if !ok || req == nil {
		return h.handleError(c, errors.BadRequest("invalid request body"))
	}

	newDomain := domain.Domain{
		Name:      req.Name,
		CreatedBy: &userId,
		UpdatedBy: &userId,
	}

	if err := h.useCase.Create(ctx, &newDomain); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.Success(CreateDomainResponse{
		Id: newDomain.Id,
	}, "Domain created"))
}
