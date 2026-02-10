package domains

import (
	gofiber "github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/domains/delivery/http"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/domains/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/domains/repository/postgresql"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/domains/usecase"
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
		http.NewDomaiHandler,
	),
	fx.Invoke(registerRoutes),
)

func registerRoutes(h *http.DomainHandler, app *gofiber.App) {
	h.RegisterRoutes(app)
}
