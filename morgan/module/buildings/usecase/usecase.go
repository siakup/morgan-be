package usecase

import (
	"github.com/siakup/morgan-be/morgan/module/buildings/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type UseCase struct {
	repository domain.BuildingRepository
	tracer     trace.Tracer
}

func NewUseCase(repository domain.BuildingRepository) *UseCase {
	return &UseCase{
		repository: repository,
		tracer:     otel.Tracer("morgan/module/buildings/usecase"),
	}
}
