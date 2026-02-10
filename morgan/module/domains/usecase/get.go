package usecase

import (
	"context"
	errs "errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/domains/domain"
)

// Get finds a domain by their unique identifier.
func (u *UseCase) Get(ctx context.Context, id string) (*domain.Domain, error) {
	ctx, span := u.tracer.Start(ctx, "Get")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	domain, err := u.repository.FindByID(ctx, id)
	if err != nil {
		if errs.Is(err, pgx.ErrNoRows) {
			return nil, errors.NotFound("domain not found")
		}

		logger.Error().
			Str("func", "repository.FindByID").
			Err(err).
			Msg("failed to find domain by id")
		return nil, errors.InternalServerError("failed to find domain by id")
	}

	return domain, nil
}
