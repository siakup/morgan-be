package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/domains/domain"
)

type (
	CreateDomainRequest struct {
		Name string `json:"name" validate:"required"`
	}
	CreateDomainResponse struct {
		Id string `json:"id"`
	}
)

// CreateDomain handles POST /domains
func (h *DomainHandler) CreateDomain(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userId, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || userId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	var req CreateDomainRequest
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.BadRequest("Invalid request body"))
	}

	// Manual basic validation
	if req.Name == "" {
		return h.handleError(c, errors.BadRequest("field name is required"))
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
