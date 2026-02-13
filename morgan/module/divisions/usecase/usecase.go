package usecase

import (
	"github.com/siakup/morgan-be/morgan/module/divisions/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type UseCase struct {
	repository domain.DivisionRepository
	tracer     trace.Tracer
}

func NewUseCase(repository domain.DivisionRepository) *UseCase {
	return &UseCase{
		repository: repository,
		tracer:     otel.Tracer("morgan/module/divisions/usecase"),
	}
}
