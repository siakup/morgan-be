package usecase

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_sessions/domain"
)

var _ domain.UseCase = (*UseCase)(nil)

// UseCase implements the logic for shift sessions management.
type UseCase struct {
	repository domain.ShiftSessionRepository
	tracer     trace.Tracer
}

// NewUseCase creates a new instance of ShiftSessions UseCase.
func NewUseCase(repository domain.ShiftSessionRepository) *UseCase {
	return &UseCase{
		repository: repository,
		tracer:     otel.Tracer("shift_sessions"),
	}
}
