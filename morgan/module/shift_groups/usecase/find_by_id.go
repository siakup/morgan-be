package usecase

import (
	"context"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

func (u *UseCase) FindByID(ctx context.Context, id string) (*domain.ShiftGroup, error) {
	return u.repo.FindByID(ctx, id)
}
