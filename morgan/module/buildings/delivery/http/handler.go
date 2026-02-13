package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/morgan/module/buildings/domain"
)

type BuildingHandler struct {
	useCase domain.UseCase
	auth    *middleware.AuthorizationMiddleware
}

func NewBuildingHandler(useCase domain.UseCase, auth *middleware.AuthorizationMiddleware) *BuildingHandler {
	return &BuildingHandler{
		useCase: useCase,
		auth:    auth,
	}
}

func (h *BuildingHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/buildings", middleware.TraceMiddleware)

	// Adjust permissions as needed
	group.Get("/", h.auth.Authenticate("buildings.view"), h.GetBuildings)
	group.Get("/:id", h.auth.Authenticate("buildings.view"), h.GetBuildingByID)
	group.Post("/", h.auth.Authenticate("buildings.edit"), h.CreateBuilding)
	group.Put("/:id", h.auth.Authenticate("buildings.edit"), h.UpdateBuilding)
	group.Delete("/:id", h.auth.Authenticate("buildings.edit"), h.DeleteBuilding)
}

func (h *BuildingHandler) handleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(responses.Fail(string(appErr.Type), appErr.Message))
	}

	return c.Status(500).JSON(responses.Fail("SYSTEM_ERROR", err.Error()))
}
