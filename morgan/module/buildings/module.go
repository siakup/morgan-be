package buildings

import (
	gofiber "github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/morgan/module/buildings/delivery/http"
	"github.com/siakup/morgan-be/morgan/module/buildings/domain"
	"github.com/siakup/morgan-be/morgan/module/buildings/repository/postgresql"
	"github.com/siakup/morgan-be/morgan/module/buildings/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			postgresql.NewBuildingRepository,
			fx.As(new(domain.BuildingRepository)),
		),
		fx.Annotate(
			usecase.NewUseCase,
			fx.As(new(domain.UseCase)),
		),
		http.NewBuildingHandler,
	),
	fx.Invoke(registerRoutes),
)

func registerRoutes(h *http.BuildingHandler, app *gofiber.App) {
	h.RegisterRoutes(app)
}
