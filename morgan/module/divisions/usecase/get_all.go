package usecase

import (
	"context"

	"github.com/siakup/morgan-be/morgan/module/divisions/domain"
)

func (u *UseCase) GetAll(ctx context.Context, filter domain.DivisionFilter) ([]*domain.Division, int64, error) {
	ctx, span := u.tracer.Start(ctx, "GetAll")
	defer span.End()

	return u.repository.FindAll(ctx, filter)
}
