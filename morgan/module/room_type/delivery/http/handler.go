package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

// RoomTypeHandler handles HTTP requests for the room type module.
type RoomTypeHandler struct {
	useCase domain.UseCase
	auth    *middleware.AuthorizationMiddleware
}

// NewRoomTypeHandler creates a new RoomTypeHandler.
func NewRoomTypeHandler(useCase domain.UseCase, auth *middleware.AuthorizationMiddleware) *RoomTypeHandler {
	return &RoomTypeHandler{
		useCase: useCase,
		auth:    auth,
	}
}

// RegisterRoutes registers the routes for the room type module.
func (h *RoomTypeHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/room-types", middleware.TraceMiddleware)

	group.Get("/", h.auth.Authenticate("room_types.view"), h.GetRoomTypes)
	group.Get("/:id", h.auth.Authenticate("room_types.view"), h.GetRoomTypeByID)
	group.Post("/", h.auth.Authenticate("room_types.edit"), h.CreateRoomType)
	group.Put("/:id", h.auth.Authenticate("room_types.edit"), h.UpdateRoomType)
	group.Delete("/:id", h.auth.Authenticate("room_types.edit"), h.DeleteRoomType)
}

func (h *RoomTypeHandler) handleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(responses.Fail(string(appErr.Type), appErr.Message))
	}
	return c.Status(http.StatusInternalServerError).JSON(responses.Fail("SYSTEM_ERROR", err.Error()))
}
