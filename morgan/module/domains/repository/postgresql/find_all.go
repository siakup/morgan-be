package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/siakup/morgan-be/libraries/object"
	"github.com/siakup/morgan-be/morgan/module/domains/domain"
)

// FindAll retrieves a list of domains based on the provided filter.
func (r *Repository) FindAll(ctx context.Context, filter domain.DomainFilter) ([]*domain.Domain, int64, error) {
	baseQuery := `
	FROM master.domains
	WHERE deleted_at IS NULL
	`
	args := pgx.NamedArgs{}

	// Simple search implementation
	if filter.Search != "" {
		baseQuery += " AND (name ILIKE @search)"
		args["search"] = "%" + filter.Search + "%"
	}

	// 1. Count Total
	var total int64
	countQuery := "SELECT count(id)" + baseQuery
	if err := r.db.QueryRow(ctx, countQuery, args).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 2. Select Data
	selectQuery := `
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
	` + baseQuery + " ORDER BY created_at DESC LIMIT @limit OFFSET @offset"

	args["limit"] = filter.Pagination.GetLimit()
	args["offset"] = filter.Pagination.GetOffset()

	rows, err := r.db.Query(ctx, selectQuery, args)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	records, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[DomainEntity])
	if err != nil {
		return nil, 0, err
	}

	domains, err := object.ParseAll[*DomainEntity, *domain.Domain](object.TagDB, object.TagObject, records)
	if err != nil {
		return nil, 0, err
	}

	return domains, total, nil
}
