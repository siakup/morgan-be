package redirect

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/redirect/delivery/http"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/redirect/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/redirect/repository/postgresql"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/redirect/usecase"
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
