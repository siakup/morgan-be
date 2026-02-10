package domain

import (
	"context"
)

// UseCase defines the business logic for the room type module.
type UseCase interface {
	FindAll(ctx context.Context, filter RoomTypeFilter) ([]*RoomType, int64, error)
	FindByID(ctx context.Context, id string) (*RoomType, error)
	Create(ctx context.Context, roomType *RoomType) error
	Update(ctx context.Context, roomType *RoomType) error
	Delete(ctx context.Context, id string, deletedBy string) error
}
