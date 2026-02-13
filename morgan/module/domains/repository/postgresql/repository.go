package postgresql

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siakup/morgan-be/morgan/module/domains/domain"
)

var _ domain.DomainRepository = (*Repository)(nil)

// DomainEntity maps to master.domains table.
type DomainEntity struct {
	Id        string     `db:"id" map:"Id"`
	Name      string     `db:"name" map:"Name"`
	Status    bool       `db:"status" map:"Status"`
	CreatedAt time.Time  `db:"created_at" map:"CreatedAt"`
	UpdatedAt time.Time  `db:"updated_at" map:"UpdatedAt"`
	DeletedAt *time.Time `db:"deleted_at" map:"DeletedAt"`
	CreatedBy *string    `db:"created_by" map:"CreatedBy"`
	UpdatedBy *string    `db:"updated_by" map:"UpdatedBy"`
	DeletedBy *string    `db:"deleted_by" map:"DeletedBy"`
}

// Repository implements domain.DomainRepository.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new Domain Repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}
