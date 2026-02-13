package usecase

import (
	"context"

	"github.com/siakup/morgan-be/morgan/module/buildings/domain"
)

func (u *UseCase) GetAll(ctx context.Context, filter domain.BuildingFilter) ([]*domain.Building, int64, error) {
	ctx, span := u.tracer.Start(ctx, "GetAll")
	defer span.End()

	return u.repository.FindAll(ctx, filter)
}
