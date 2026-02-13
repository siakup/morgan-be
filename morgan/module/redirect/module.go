package redirect

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"github.com/siakup/morgan-be/morgan/module/redirect/delivery/http"
	"github.com/siakup/morgan-be/morgan/module/redirect/domain"
	"github.com/siakup/morgan-be/morgan/module/redirect/repository/postgresql"
	"github.com/siakup/morgan-be/morgan/module/redirect/usecase"
)

var Module = fx.Module(
	"redirect",
	fx.Provide(
		postgresql.NewRepository,
		fx.Annotate(
			postgresql.NewRepository,
			fx.As(new(domain.RedirectRepository)),
		),
		usecase.NewUseCase,
		fx.Annotate(
			usecase.NewUseCase,
			fx.As(new(domain.RedirectUseCase)),
		),
		http.NewHandler,
	),
	fx.Invoke(
		registerRoutes,
	),
)

func registerRoutes(h *http.RedirectHandler, app *fiber.App) {
	h.RegisterRoutes(app)
}
