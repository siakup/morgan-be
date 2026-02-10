package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/object"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/domains/domain"
)

var queryFindById = `
	SELECT
		id,
    name,
    status,
    created_at,
    updated_at,
    deleted_at,
    created_by,
    updated_by,
    deleted_by
	FROM organization.domains
	WHERE id = @id AND deleted_at IS NULL
	LIMIT 1
`

// FindByID retrieves a domain by its ID.
func (r *Repository) FindByID(ctx context.Context, id string) (*domain.Domain, error) {
	rows, err := r.db.Query(ctx, queryFindById, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	record, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[DomainEntity])
	if err != nil {
		return nil, err
	}

	domain, err := object.Parse[*DomainEntity, *domain.Domain](object.TagDB, object.TagObject, record)
	if err != nil {
		return nil, err
	}

	return domain, nil
}
