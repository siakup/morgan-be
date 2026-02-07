package usecase

import (
	"context"
	errs "errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

// Update modifies an existing room type. Name must remain unique (excluding current record).
func (u *UseCase) Update(ctx context.Context, rt *domain.RoomType) error {
	logger := zerolog.Ctx(ctx)

	current, err := u.repo.FindByID(ctx, rt.Id)
	if err != nil {
		return errors.InternalServerError("failed to find room type")
	}
	if current == nil {
		return errors.NotFound("room type not found")
	}

	if current.Name != rt.Name {
		existing, err := u.repo.FindByName(ctx, rt.Name)
		if err != nil && !errs.Is(err, pgx.ErrNoRows) {
			logger.Error().Err(err).Str("name", rt.Name).Msg("failed to check room type name uniqueness")
			return errors.InternalServerError("failed to validate room type name")
		}
		if existing != nil && existing.Id != rt.Id {
			return errors.BadRequest("room type name already exists")
		}
	}

	rt.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, rt); err != nil {
		logger.Error().Err(err).Str("func", "repository.Update").Msg("failed to update room type")
		return errors.InternalServerError("failed to update room type")
	}

	return nil
}
