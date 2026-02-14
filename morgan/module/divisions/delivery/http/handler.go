package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/morgan/module/divisions/domain"
)

type DivisionHandler struct {
	useCase domain.UseCase
	auth    *middleware.AuthorizationMiddleware
}

func NewDivisionHandler(useCase domain.UseCase, auth *middleware.AuthorizationMiddleware) *DivisionHandler {
	return &DivisionHandler{
		useCase: useCase,
		auth:    auth,
	}
}

func (h *DivisionHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/divisions", middleware.TraceMiddleware)

	// Adjust permissions as needed - assuming similar pattern to severity_levels
	group.Get("/", h.auth.Authenticate("divisions.view"), h.GetDivisions)
	group.Get("/:id", h.auth.Authenticate("divisions.view"), h.GetDivisionByID)
	group.Post("/", h.auth.Authenticate("divisions.edit"), h.CreateDivision)
	group.Put("/:id", h.auth.Authenticate("divisions.edit"), h.UpdateDivision)
	group.Delete("/:id", h.auth.Authenticate("divisions.edit"), h.DeleteDivision)
}

func (h *DivisionHandler) handleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(responses.Fail(string(appErr.Type), appErr.Message))
	}

	return c.Status(http.StatusInternalServerError).JSON(responses.Fail("SYSTEM_ERROR", err.Error()))
}
