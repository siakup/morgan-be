package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
)

type (
	SyncUserRequest struct {
		Code string `json:"code" validate:"required"`
	}
	SyncUserResponse struct {
		Id string `json:"id"`
	}
)

// SyncUser handles POST /users
func (h *UserHandler) SyncUser(c *fiber.Ctx) error {
	ctx := c.UserContext()

	// Get institution_id from locals
	institutionId, ok := c.Locals(middleware.XInstitutionId).(string)
	if !ok || institutionId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing institution context"))
	}

	var req SyncUserRequest
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.BadRequest("Invalid request body"))
	}

	if req.Code == "" {
		return h.handleError(c, errors.BadRequest("code is required"))
	}

	// Get token from locals
	token, ok := c.Locals(middleware.XTokenKey).(string)
	if !ok || token == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing token"))
	}

	user, err := h.useCase.SyncUser(ctx, institutionId, token, req.Code)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.Success(SyncUserResponse{
		Id: user.Id,
	}, "User synced"))
}
