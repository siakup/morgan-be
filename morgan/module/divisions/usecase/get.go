package usecase

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	liberrors "github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/morgan/module/divisions/domain"
)

func (u *UseCase) Get(ctx context.Context, id string) (*domain.Division, error) {
	ctx, span := u.tracer.Start(ctx, "Get")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	division, err := u.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, liberrors.NotFound("division not found")
		}

		logger.Error().
			Str("func", "repository.FindByID").
			Err(err).
			Msg("failed to find division by id")
		return nil, liberrors.InternalServerError("failed to find division by id")
	}

	return division, nil
}
