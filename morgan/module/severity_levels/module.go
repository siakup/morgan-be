package severity_levels

import (
	gofiber "github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/delivery/http"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/repository/postgresql"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/usecase"
)

// Module exports the severity_levels module for Fx.
var Module = fx.Options(
	fx.Provide(
		postgresql.NewRepository,
		fx.Annotate(
			postgresql.NewRepository,
			fx.As(new(domain.SeverityLevelRepository)),
		),
		usecase.NewUseCase,
		fx.Annotate(
			usecase.NewUseCase,
			fx.As(new(domain.UseCase)),
		),
		http.NewSeverityLevelHandler,
	),
	fx.Invoke(registerRoutes),
)

func registerRoutes(h *http.SeverityLevelHandler, app *gofiber.App) {
	h.RegisterRoutes(app)
}
