package usecase

import (
	"context"
	errs "errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

// Create persists a new room type. Name must be unique (no duplicate).
func (u *UseCase) Create(ctx context.Context, rt *domain.RoomType) error {
	logger := zerolog.Ctx(ctx)

	existing, err := u.repo.FindByName(ctx, rt.Name)
	if err != nil && !errs.Is(err, pgx.ErrNoRows) {
		logger.Error().Err(err).Str("name", rt.Name).Msg("failed to check room type name uniqueness")
		return errors.InternalServerError("failed to validate room type name")
	}
	if existing != nil {
		return errors.BadRequest("room type name already exists")
	}

	rt.Id = uuid.NewString()
	now := time.Now()
	rt.CreatedAt = now
	rt.UpdatedAt = now

	if err := u.repo.Store(ctx, rt); err != nil {
		logger.Error().Err(err).Str("func", "repository.Store").Msg("failed to store room type")
		return errors.InternalServerError("failed to store room type")
	}

	return nil
}
