package domain

import (
	"context"
)

type BuildingRepository interface {
	FindAll(ctx context.Context, filter BuildingFilter) ([]*Building, int64, error)
	FindByID(ctx context.Context, id string) (*Building, error)
	Store(ctx context.Context, building *Building) error
	Update(ctx context.Context, building *Building) error
	Delete(ctx context.Context, id string, deletedBy string) error
}
