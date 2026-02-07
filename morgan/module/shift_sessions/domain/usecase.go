package domain

import (
	"context"
)

// UseCase defines the business logic for the roles module.
type UseCase interface {
	FindAll(ctx context.Context, filter ShiftSessionFilter) ([]*ShiftSession, int64, error)
	Get(ctx context.Context, id string) (*ShiftSession, error)
	Create(ctx context.Context, role *ShiftSession) error
	Update(ctx context.Context, role *ShiftSession) error
	Delete(ctx context.Context, id string, deletedBy string) error
}
