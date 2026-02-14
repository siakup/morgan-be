package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/siakup/morgan-be/libraries/object"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/domain"
)

// FindAll retrieves a list of shift sessions based on the provided filter.
func (r *Repository) FindAll(ctx context.Context, filter domain.ShiftSessionFilter) ([]*domain.ShiftSession, int64, error) {
	baseQuery := `
    FROM hr.shift_sessions
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
    start,
    "end",
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

	records, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[ShiftSessionEntity])
	if err != nil {
		return nil, 0, err
	}

	shiftSessions, err := object.ParseAll[*ShiftSessionEntity, *domain.ShiftSession](object.TagDB, object.TagObject, records)
	if err != nil {
		return nil, 0, err
	}

	return shiftSessions, total, nil
}
