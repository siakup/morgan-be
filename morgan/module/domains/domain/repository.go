package domain

import (
	"context"
	"time"

	"github.com/siakup/morgan-be/libraries/types"
)

// Domain represents the domain Domain entity.
type Domain struct {
	Id        string     `object:"id"`
	Name      string     `object:"name"`
	Status    bool       `object:"status"`
	CreatedAt time.Time  `object:"created_at"`
	UpdatedAt time.Time  `object:"updated_at"`
	DeletedAt *time.Time `object:"deleted_at"` // Pointer for nullable
	CreatedBy *string    `object:"created_by"` // Nullable
	UpdatedBy *string    `object:"updated_by"` // Nullable
	DeletedBy *string    `object:"deleted_by"` // Nullable
}

// DomainFilter represents filter options for listing Domains.
type DomainFilter struct {
	types.Pagination
	Status string
	Search string // Search in name
}

// DomainRepository defines the persistence layer contract.
type DomainRepository interface {
	FindAll(ctx context.Context, filter DomainFilter) ([]*Domain, int64, error)
	FindByID(ctx context.Context, id string) (*Domain, error)
	Store(ctx context.Context, Domain *Domain) error
	Update(ctx context.Context, Domain *Domain) error
	Delete(ctx context.Context, id string, deletedBy string) error
}
