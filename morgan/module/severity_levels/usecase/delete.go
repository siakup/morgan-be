package usecase

import (
	"context"
)

func (u *UseCase) Delete(ctx context.Context, id string, deletedBy string) error {
	return u.repo.Delete(ctx, id, deletedBy)
}
