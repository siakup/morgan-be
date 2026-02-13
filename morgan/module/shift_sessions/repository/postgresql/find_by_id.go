package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/siakup/morgan-be/libraries/object"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/domain"
)

var queryFindById = `
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
	FROM schedule.shift_sessions
	WHERE id = @id AND deleted_at IS NULL
	LIMIT 1
`

// FindByID retrieves a single shift session by its ID.
func (r *Repository) FindByID(ctx context.Context, id string) (*domain.ShiftSession, error) {
	rows, err := r.db.Query(ctx, queryFindById, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	record, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[ShiftSessionEntity])
	if err != nil {
		return nil, err
	}

	shiftSession, err := object.Parse[*ShiftSessionEntity, *domain.ShiftSession](object.TagDB, object.TagObject, record)
	if err != nil {
		return nil, err
	}

	return shiftSession, nil
}
