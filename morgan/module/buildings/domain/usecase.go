package domain

import (
	"context"
)

type UseCase interface {
	GetAll(ctx context.Context, filter BuildingFilter) ([]*Building, int64, error)
	Get(ctx context.Context, id string) (*Building, error)
	Create(ctx context.Context, building *Building) error
	Update(ctx context.Context, building *Building) error
	Delete(ctx context.Context, id string, deletedBy string) error
}
