package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/morgan/module/domains/domain"
)

// DomainHandler handles HTTP requests for domains module.
type DomainHandler struct {
	useCase domain.UseCase
	auth    *middleware.AuthorizationMiddleware
}

// NewDomaiHandler creates a new DomainHandler.
func NewDomaiHandler(useCase domain.UseCase, auth *middleware.AuthorizationMiddleware) *DomainHandler {
	return &DomainHandler{
		useCase: useCase,
		auth:    auth,
	}
}

// RegisterRoutes registers the routes for the domains module.
func (h *DomainHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/domains", middleware.TraceMiddleware)

	// group.Get("/", h.auth.Authenticate("domains.manage.all.view"), h.GetDomains)
	// group.Get("/:id", h.auth.Authenticate("domains.manage.all.view"), h.GetDomainSessionByID)
	// group.Post("/", h.auth.Authenticate("domains.manage.all.edit"), h.CreateDomain)
	// group.Put("/:id", h.auth.Authenticate("domains.manage.all.edit"), h.UpdateDomain)
	// group.Delete("/:id", h.auth.Authenticate("domains.manage.all.edit"), h.DeleteDomain)

	group.Get("/", h.GetDomains)
	group.Get("/:id", h.GetDomainByID)
	group.Post("/", h.CreateDomain)
	group.Put("/:id", h.UpdateDomain)
	group.Delete("/:id", h.DeleteDomain)
}

// handleError handles errors by mapping them to standardized responses.
func (h *DomainHandler) handleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(responses.Fail(string(appErr.Type), appErr.Message))
	}

	return c.Status(http.StatusInternalServerError).JSON(responses.Fail("SYSTEM_ERROR", err.Error()))
}
