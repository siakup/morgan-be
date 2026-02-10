package domain

import (
	"context"
	"time"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/types"
)

// SeverityLevel represents the domain object for a Severity Level.
type SeverityLevel struct {
	Id        string     `object:"id"`
	Name      string     `object:"name"`
	Status    bool       `object:"status"`
	CreatedAt time.Time  `object:"created_at"`
	CreatedBy *string    `object:"created_by"`
	UpdatedAt time.Time  `object:"updated_at"`
	UpdatedBy *string    `object:"updated_by"`
	DeletedAt *time.Time `object:"deleted_at"`
	DeletedBy *string    `object:"deleted_by"`
}

// SeverityLevelFilter represents the filter options for fetching severity levels.
type SeverityLevelFilter struct {
	types.Pagination
	Search string
}

// SeverityLevelRepository defines the methods for interacting with the severity levels storage.
type SeverityLevelRepository interface {
	FindAll(ctx context.Context, filter SeverityLevelFilter) ([]*SeverityLevel, int64, error)
	FindByID(ctx context.Context, id string) (*SeverityLevel, error)
	Store(ctx context.Context, severityLevel *SeverityLevel) error
	Update(ctx context.Context, severityLevel *SeverityLevel) error
	// Delete removes a severity level from storage.
	Delete(ctx context.Context, id string, deletedBy string) error
}
