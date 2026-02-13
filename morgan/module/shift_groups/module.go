package shift_groups

import (
	gofiber "github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"github.com/siakup/morgan-be/morgan/module/shift_groups/delivery/http"
	"github.com/siakup/morgan-be/morgan/module/shift_groups/domain"
	"github.com/siakup/morgan-be/morgan/module/shift_groups/repository/postgresql"
	"github.com/siakup/morgan-be/morgan/module/shift_groups/usecase"
)

// Module exports the shift_groups module for Fx.
var Module = fx.Options(
	fx.Provide(
		postgresql.NewRepository,
		fx.Annotate(
			postgresql.NewRepository,
			fx.As(new(domain.ShiftGroupRepository)),
		),
		usecase.NewUseCase,
		fx.Annotate(
			usecase.NewUseCase,
			fx.As(new(domain.UseCase)),
		),
		http.NewShiftGroupHandler,
	),
	fx.Invoke(registerRoutes),
)

func registerRoutes(h *http.ShiftGroupHandler, app *gofiber.App) {
	h.RegisterRoutes(app)
}
