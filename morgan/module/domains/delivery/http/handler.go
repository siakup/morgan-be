package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/validation"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/morgan/module/domains/domain"
)

type DomainHandler struct {
	useCase domain.UseCase
	auth    *middleware.AuthorizationMiddleware
}

func NewDomainHandler(useCase domain.UseCase, auth *middleware.AuthorizationMiddleware) *DomainHandler {
	return &DomainHandler{
		useCase: useCase,
		auth:    auth,
	}
}

func (h *DomainHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/domains", middleware.TraceMiddleware)

	group.Get("/", h.auth.Authenticate("domains.organization.domains.view"), h.GetDomains)
	group.Get("/:id", h.auth.Authenticate("domains.organization.domains.view"), h.GetDomainByID)
	group.Post("/", h.auth.Authenticate("domains.organization.domains.create"), validation.ValidateBody(func() interface{} { return &CreateDomainRequest{} }), h.CreateDomain)
	group.Put("/:id", h.auth.Authenticate("domains.organization.domains.edit"), validation.ValidateBody(func() interface{} { return &UpdateDomainRequest{} }), h.UpdateDomain)
	group.Delete("/:id", h.auth.Authenticate("domains.organization.domains.edit"), h.DeleteDomain)

}

func (h *DomainHandler) handleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(responses.Fail(string(appErr.Type), appErr.Message))
	}

	return c.Status(http.StatusInternalServerError).JSON(responses.Fail("SYSTEM_ERROR", err.Error()))
}
