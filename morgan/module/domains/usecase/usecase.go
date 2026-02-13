package usecase

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"github.com/siakup/morgan-be/morgan/module/domains/domain"
)

var _ domain.UseCase = (*UseCase)(nil)

// UseCase implements the logic for domains management.
type UseCase struct {
	repository domain.DomainRepository
	tracer     trace.Tracer
}

// NewUseCase creates a new instance of Domains UseCase.
func NewUseCase(repository domain.DomainRepository) *UseCase {
	return &UseCase{
		repository: repository,
		tracer:     otel.Tracer("domains"),
	}
}
