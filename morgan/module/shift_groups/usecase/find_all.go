package usecase

import (
	"context"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

func (u *UseCase) FindAll(ctx context.Context, filter domain.ShiftGroupFilter) ([]*domain.ShiftGroup, int64, error) {
	return u.repo.FindAll(ctx, filter)
}
