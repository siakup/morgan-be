package postgresql

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

var _ domain.ShiftGroupRepository = (*Repository)(nil)

// ShiftGroupEntity represents the schema in the database.
type ShiftGroupEntity struct {
	Id        string         `db:"id"`
	Name      string         `db:"name"`
	Status    bool           `db:"status"`
	CreatedAt time.Time      `db:"created_at"`
	CreatedBy sql.NullString `db:"created_by"`
	UpdatedAt time.Time      `db:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by"`
}

// Repository implements the domain.ShiftGroupRepository interface for PostgreSQL.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new instance of the PostgreSQL Repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func nullStringToPointer(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullTimeToPointer(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
