package shift_sessions

import (
	gofiber "github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/delivery/http"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/domain"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/repository/postgresql"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/usecase"
)

// Module exports the roles module for Fx.
var Module = fx.Options(
	fx.Provide(
		postgresql.NewRepository,
		fx.Annotate(
			postgresql.NewRepository,
			fx.As(new(domain.ShiftSessionRepository)),
		),
		usecase.NewUseCase,
		fx.Annotate(
			usecase.NewUseCase,
			fx.As(new(domain.UseCase)),
		),
		http.NewShiftSessionHandler,
	),
	fx.Invoke(registerRoutes),
)

func registerRoutes(h *http.ShiftSessionHandler, app *gofiber.App) {
	h.RegisterRoutes(app)
}
