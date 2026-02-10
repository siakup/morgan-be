package usecase

import (
	"context"

	"github.com/rs/zerolog"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

func (u *UseCase) Delete(ctx context.Context, id string, deletedBy string) error {
	logger := zerolog.Ctx(ctx)

	existing, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return errors.InternalServerError("failed to find room type")
	}
	if existing == nil {
		return errors.NotFound("room type not found")
	}

	if err := u.repo.Delete(ctx, id, deletedBy); err != nil {
		logger.Error().Err(err).Str("func", "repository.Delete").Str("id", id).Msg("failed to delete room type")
		return errors.InternalServerError("failed to delete room type")
	}

	return nil
}
