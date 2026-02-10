package usecase

import (
	"context"
	errs "errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/domains/domain"
)

// Update modifies an existing domain.
func (u *UseCase) Update(ctx context.Context, domain *domain.Domain) error {
	ctx, span := u.tracer.Start(ctx, "Update")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	// Verify exists
	_, err := u.repository.FindByID(ctx, domain.Id)
	if err != nil {
		if errs.Is(err, pgx.ErrNoRows) {
			return errors.NotFound("domain not found")
		}
		return errors.InternalServerError("failed to find domain")
	}

	if err := u.repository.Update(ctx, domain); err != nil {
		logger.Error().
			Str("func", "repository.Update").
			Err(err).
			Msg("failed to update domain")

		return errors.InternalServerError("failed to update domain")
	}

	return nil
}
