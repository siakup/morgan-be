package usecase

import (
	"context"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

func (u *UseCase) FindAll(ctx context.Context, filter domain.RoomTypeFilter) ([]*domain.RoomType, int64, error) {
	return u.repo.FindAll(ctx, filter)
}
