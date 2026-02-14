package domain

import (
	"context"
)

type UseCase interface {
	GetAll(ctx context.Context, filter DivisionFilter) ([]*Division, int64, error)
	Get(ctx context.Context, id string) (*Division, error)
	Create(ctx context.Context, division *Division) error
	Update(ctx context.Context, division *Division) error
	Delete(ctx context.Context, id string, deletedBy string) error
}
