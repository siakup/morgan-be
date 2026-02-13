package usecase

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	liberrors "github.com/siakup/morgan-be/libraries/errors" // Alias for custom errors
	"github.com/siakup/morgan-be/morgan/module/shift_groups/domain"
	"go.opentelemetry.io/otel/trace"
)

func (u *UseCase) Get(ctx context.Context, id string) (*domain.ShiftGroup, error) {
	ctx, span := u.tracer.Start(ctx, "usecase.Get", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	shiftGroup, err := u.repo.FindByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get shift group")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, liberrors.NotFound("Shift Group not found")
		}
		return nil, err
	}

	return shiftGroup, nil
}
