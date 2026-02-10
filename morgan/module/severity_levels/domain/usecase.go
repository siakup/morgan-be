package domain

import (
	"context"
)

// UseCase defines the business logic for severity levels.
type UseCase interface {
	FindAll(ctx context.Context, filter SeverityLevelFilter) ([]*SeverityLevel, int64, error)
	FindByID(ctx context.Context, id string) (*SeverityLevel, error)
	Create(ctx context.Context, severityLevel *SeverityLevel) error
	Update(ctx context.Context, severityLevel *SeverityLevel) error
	Delete(ctx context.Context, id string, deletedBy string) error
}
