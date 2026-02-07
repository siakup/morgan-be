package room_type

import (
	gofiber "github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/delivery/http"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/repository/postgresql"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/usecase"
)

// Module exports the room_type module for Fx.
var Module = fx.Options(
	fx.Provide(
		postgresql.NewRepository,
		fx.Annotate(
			postgresql.NewRepository,
			fx.As(new(domain.RoomTypeRepository)),
		),
		usecase.NewUseCase,
		fx.Annotate(
			usecase.NewUseCase,
			fx.As(new(domain.UseCase)),
		),
		http.NewRoomTypeHandler,
	),
	fx.Invoke(registerRoutes),
)

func registerRoutes(h *http.RoomTypeHandler, app *gofiber.App) {
	h.RegisterRoutes(app)
}
