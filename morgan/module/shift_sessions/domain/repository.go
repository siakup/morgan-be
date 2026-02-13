package domain

import (
	"context"
	"time"

	"github.com/siakup/morgan-be/libraries/types"
)

// ShiftSession represents the domain Shift Session entity.
type ShiftSession struct {
	Id        string     `object:"id"`
	Name      string     `object:"name"`
	Start     string     `object:"start"`
	End       string     `object:"end"`
	Status    bool       `object:"status"`
	CreatedAt time.Time  `object:"created_at"`
	UpdatedAt time.Time  `object:"updated_at"`
	DeletedAt *time.Time `object:"deleted_at"` // Pointer for nullable
	CreatedBy *string    `object:"created_by"` // Nullable
	UpdatedBy *string    `object:"updated_by"` // Nullable
	DeletedBy *string    `object:"deleted_by"` // Nullable
}

// ShiftSessionFilter represents filter options for listing Shift Sessions.
type ShiftSessionFilter struct {
	types.Pagination
	Status string
	Search string // Search in name
}

// ShiftSessionRepository defines the persistence layer contract.
type ShiftSessionRepository interface {
	FindAll(ctx context.Context, filter ShiftSessionFilter) ([]*ShiftSession, int64, error)
	FindByID(ctx context.Context, id string) (*ShiftSession, error)
	Store(ctx context.Context, shiftSession *ShiftSession) error
	Update(ctx context.Context, shiftSession *ShiftSession) error
	Delete(ctx context.Context, id string, deletedBy string) error
}
