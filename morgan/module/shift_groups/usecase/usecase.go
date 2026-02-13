package usecase

import (
	"github.com/siakup/morgan-be/morgan/module/shift_groups/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var _ domain.UseCase = (*UseCase)(nil)

// UseCase implements the domain.UseCase interface.
type UseCase struct {
	repo   domain.ShiftGroupRepository
	tracer trace.Tracer
}

// NewUseCase creates a new instance of the UseCase.
func NewUseCase(repo domain.ShiftGroupRepository) *UseCase {
	return &UseCase{
		repo:   repo,
		tracer: otel.Tracer("morgan/module/shift_groups/usecase"),
	}
}
