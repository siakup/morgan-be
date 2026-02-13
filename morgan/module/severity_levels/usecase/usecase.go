package usecase

import (
	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var _ domain.UseCase = (*UseCase)(nil)

// UseCase implements the domain.UseCase interface.
type UseCase struct {
	repository domain.SeverityLevelRepository
	tracer     trace.Tracer
}

// NewUseCase creates a new UseCase.
func NewUseCase(repository domain.SeverityLevelRepository) domain.UseCase {
	return &UseCase{
		repository: repository,
		tracer:     otel.Tracer("morgan/module/severity_levels/usecase"),
	}
}
