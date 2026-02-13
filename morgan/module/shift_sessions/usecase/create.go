package usecase

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/domain"
)

// Create persists a new shift session record.
func (u *UseCase) Create(ctx context.Context, shiftSession *domain.ShiftSession) error {
	ctx, span := u.tracer.Start(ctx, "Create")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	if err := u.repository.Store(ctx, shiftSession); err != nil {
		logger.Error().
			Str("func", "repository.Store").
			Err(err).
			Msg("failed to store shift session")

		return errors.InternalServerError("failed to store shift session")
	}

	return nil
}
