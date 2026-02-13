package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
)

// DeleteDomain handles DELETE /domains/:id
func (h *DomainHandler) DeleteDomain(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("Invalid ID"))
	}

	deletedBy, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || deletedBy == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	if err := h.useCase.Delete(ctx, id, deletedBy); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success[any](nil, "Domain deleted"))
}
