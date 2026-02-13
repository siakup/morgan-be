package usecase

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/morgan/module/domains/domain"
)

// FindAll retrieves a list of domains based on filter criteria.
func (u *UseCase) FindAll(ctx context.Context, filter domain.DomainFilter) ([]*domain.Domain, int64, error) {
	ctx, span := u.tracer.Start(ctx, "FindAll")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	domains, total, err := u.repository.FindAll(ctx, filter)
	if err != nil {
		logger.Error().
			Str("func", "repository.FindAll").
			Err(err).
			Msg("failed to find domains")

		return nil, 0, errors.InternalServerError("failed to find domains")
	}

	return domains, total, nil
}
