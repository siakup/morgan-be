package postgresql

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

var _ domain.RoomTypeRepository = (*Repository)(nil)

// RoomTypeEntity represents the schema in the database.
type RoomTypeEntity struct {
	Id          string
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	CreatedBy   sql.NullString
	UpdatedAt   time.Time
	UpdatedBy   sql.NullString
	DeletedAt   sql.NullTime
	DeletedBy   sql.NullString
}

// Repository implements the domain.RoomTypeRepository interface for PostgreSQL.
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
