package usecase

import (
	"context"
	errs "errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/domain"
)

// Update modifies an existing shift session.
func (u *UseCase) Update(ctx context.Context, shiftSession *domain.ShiftSession) error {
	ctx, span := u.tracer.Start(ctx, "Update")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	// Verify exists
	_, err := u.repository.FindByID(ctx, shiftSession.Id)
	if err != nil {
		if errs.Is(err, pgx.ErrNoRows) {
			return errors.NotFound("shift session not found")
		}
		return errors.InternalServerError("failed to find shift session")
	}

	if err := u.repository.Update(ctx, shiftSession); err != nil {
		logger.Error().
			Str("func", "repository.Update").
			Err(err).
			Msg("failed to update shift session")

		return errors.InternalServerError("failed to update shift session")
	}

	return nil
}
