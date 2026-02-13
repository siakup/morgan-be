package usecase

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/morgan/module/domains/domain"
)

// Create persists a new domain record.
func (u *UseCase) Create(ctx context.Context, domain *domain.Domain) error {
	ctx, span := u.tracer.Start(ctx, "Create")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	if err := u.repository.Store(ctx, domain); err != nil {
		logger.Error().
			Str("func", "repository.Store").
			Err(err).
			Msg("failed to store domain")

		return errors.InternalServerError("failed to store domain")
	}

	return nil
}
