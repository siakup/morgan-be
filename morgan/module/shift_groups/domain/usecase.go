package domain

import (
	"context"
)

// UseCase defines the business logic methods for shift groups.
type UseCase interface {
	FindAll(ctx context.Context, filter ShiftGroupFilter) ([]*ShiftGroup, int64, error)
	FindByID(ctx context.Context, id string) (*ShiftGroup, error)
	Create(ctx context.Context, shiftGroup *ShiftGroup) error
	Update(ctx context.Context, shiftGroup *ShiftGroup) error
	Delete(ctx context.Context, id string, deletedBy string) error
}
