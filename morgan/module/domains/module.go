package domains

import (
	gofiber "github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"github.com/siakup/morgan-be/morgan/module/domains/delivery/http"
	"github.com/siakup/morgan-be/morgan/module/domains/domain"
	"github.com/siakup/morgan-be/morgan/module/domains/repository/postgresql"
	"github.com/siakup/morgan-be/morgan/module/domains/usecase"
)

// Module exports the roles module for Fx.
var Module = fx.Options(
	fx.Provide(
		postgresql.NewRepository,
		fx.Annotate(
			postgresql.NewRepository,
			fx.As(new(domain.DomainRepository)),
		),
		usecase.NewUseCase,
		fx.Annotate(
			usecase.NewUseCase,
			fx.As(new(domain.UseCase)),
		),
		http.NewDomainHandler,
	),
	fx.Invoke(registerRoutes),
)

func registerRoutes(h *http.DomainHandler, app *gofiber.App) {
	h.RegisterRoutes(app)
}
