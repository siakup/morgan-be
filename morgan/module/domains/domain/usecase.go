package domain

import (
	"context"
)

// UseCase defines the business logic for the roles module.
type UseCase interface {
	FindAll(ctx context.Context, filter DomainFilter) ([]*Domain, int64, error)
	Get(ctx context.Context, id string) (*Domain, error)
	Create(ctx context.Context, role *Domain) error
	Update(ctx context.Context, role *Domain) error
	Delete(ctx context.Context, id string, deletedBy string) error
}
