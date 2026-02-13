package usecase

import (
	"context"

	"github.com/siakup/morgan-be/morgan/module/shift_groups/domain"
)

func (u *UseCase) FindAll(ctx context.Context, filter domain.ShiftGroupFilter) ([]*domain.ShiftGroup, int64, error) {
	return u.repo.FindAll(ctx, filter)
}
