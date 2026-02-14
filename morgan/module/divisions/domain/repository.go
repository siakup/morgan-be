package domain

import (
	"context"
)

type DivisionRepository interface {
	FindAll(ctx context.Context, filter DivisionFilter) ([]*Division, int64, error)
	FindByID(ctx context.Context, id string) (*Division, error)
	Store(ctx context.Context, division *Division) error
	Update(ctx context.Context, division *Division) error
	Delete(ctx context.Context, id string, deletedBy string) error
}
