package usecase

import (
	"context"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

func (u *UseCase) FindByID(ctx context.Context, id string) (*domain.RoomType, error) {
	return u.repo.FindByID(ctx, id)
}
