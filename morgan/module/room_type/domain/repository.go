package domain

import (
	"context"
	"time"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/types"
)

// RoomType represents the domain object for Room Type Master.
type RoomType struct {
	Id          string     `object:"id"`
	Name        string     `object:"name"`
	Description string     `object:"description"`
	IsActive    bool       `object:"is_active"`
	CreatedAt   time.Time  `object:"created_at"`
	CreatedBy   *string    `object:"created_by"`
	UpdatedAt   time.Time  `object:"updated_at"`
	UpdatedBy   *string    `object:"updated_by"`
	DeletedAt   *time.Time `object:"deleted_at"`
	DeletedBy   *string    `object:"deleted_by"`
}

// RoomTypeFilter represents the filter options for fetching room types.
type RoomTypeFilter struct {
	types.Pagination
	Search string
}

// RoomTypeRepository defines the methods for interacting with the room types storage.
type RoomTypeRepository interface {
	FindAll(ctx context.Context, filter RoomTypeFilter) ([]*RoomType, int64, error)
	FindByID(ctx context.Context, id string) (*RoomType, error)
	FindByName(ctx context.Context, name string) (*RoomType, error)
	Store(ctx context.Context, roomType *RoomType) error
	Update(ctx context.Context, roomType *RoomType) error
	Delete(ctx context.Context, id string, deletedBy string) error
}
