package usecase

import (
	"context"

	"github.com/siakup/morgan-be/morgan/module/shift_groups/domain"
)

func (u *UseCase) FindByID(ctx context.Context, id string) (*domain.ShiftGroup, error) {
	return u.repo.FindByID(ctx, id)
}
