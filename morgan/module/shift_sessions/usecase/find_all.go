package usecase

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/domain"
)

// FindAll retrieves a list of shift sessions based on filter criteria.
func (u *UseCase) FindAll(ctx context.Context, filter domain.ShiftSessionFilter) ([]*domain.ShiftSession, int64, error) {
	ctx, span := u.tracer.Start(ctx, "FindAll")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	shiftSessions, total, err := u.repository.FindAll(ctx, filter)
	if err != nil {
		logger.Error().
			Str("func", "repository.FindAll").
			Err(err).
			Msg("failed to find shift sessions")

		return nil, 0, errors.InternalServerError("failed to find shift sessions")
	}

	return shiftSessions, total, nil
}
