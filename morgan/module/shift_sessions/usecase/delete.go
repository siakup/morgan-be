package usecase

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/siakup/morgan-be/libraries/errors"
)

// Delete soft removes a shift session from the system.
func (u *UseCase) Delete(ctx context.Context, id string, deletedBy string) error {
	ctx, span := u.tracer.Start(ctx, "Delete")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	if err := u.repository.Delete(ctx, id, deletedBy); err != nil {
		logger.Error().
			Str("func", "repository.Delete").
			Err(err).
			Msg("failed to delete shift session")

		return errors.InternalServerError("failed to delete shift session")
	}

	return nil
}
