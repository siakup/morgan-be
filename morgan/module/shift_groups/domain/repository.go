package domain

import (
	"context"
	"time"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/types"
)

// ShiftGroup represents the domain object for a Shift Group.
// ShiftGroup represents the domain object for a Shift Group.
type ShiftGroup struct {
	Id        string     `object:"id"`
	Name      string     `object:"name"` // Enum: FM, IT, HK
	Status    bool       `object:"status"`
	CreatedAt time.Time  `object:"created_at"`
	CreatedBy *string    `object:"created_by"`
	UpdatedAt time.Time  `object:"updated_at"`
	UpdatedBy *string    `object:"updated_by"`
	DeletedAt *time.Time `object:"deleted_at"`
	DeletedBy *string    `object:"deleted_by"`
}

// ShiftGroupFilter represents the filter options for fetching shift groups.
type ShiftGroupFilter struct {
	types.Pagination
	Search string
}

// ShiftGroupRepository defines the methods for interacting with the shift groups storage.
type ShiftGroupRepository interface {
	FindAll(ctx context.Context, filter ShiftGroupFilter) ([]*ShiftGroup, int64, error)
	FindByID(ctx context.Context, id string) (*ShiftGroup, error)
	Store(ctx context.Context, shiftGroup *ShiftGroup) error
	Update(ctx context.Context, shiftGroup *ShiftGroup) error
	// Delete removes a shift group from storage.
	Delete(ctx context.Context, id string, deletedBy string) error
}

// ShiftGroupUseCase defines the business logic for shift groups.
