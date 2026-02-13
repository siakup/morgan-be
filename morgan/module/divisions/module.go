package divisions

import (
	gofiber "github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/morgan/module/divisions/delivery/http"
	"github.com/siakup/morgan-be/morgan/module/divisions/domain"
	"github.com/siakup/morgan-be/morgan/module/divisions/repository/postgresql"
	"github.com/siakup/morgan-be/morgan/module/divisions/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			postgresql.NewDivisionRepository,
			fx.As(new(domain.DivisionRepository)),
		),
		fx.Annotate(
			usecase.NewUseCase,
			fx.As(new(domain.UseCase)),
		),
		http.NewDivisionHandler,
	),
	fx.Invoke(registerRoutes),
)

func registerRoutes(h *http.DivisionHandler, app *gofiber.App) {
	h.RegisterRoutes(app)
}
