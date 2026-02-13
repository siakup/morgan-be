package usecase

import (
	"context"
)

func (u *UseCase) Delete(ctx context.Context, id string, deletedBy string) error {
	ctx, span := u.tracer.Start(ctx, "Delete")
	defer span.End()

	return u.repository.Delete(ctx, id, deletedBy)
}
