package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/siakup/morgan-be/morgan/module/domains/domain"
)

var queryStore = `
    	INSERT INTO master.domains (
    		name, created_by, updated_by
    	) VALUES (
    		@name, @created_by, @updated_by
    	)
    	RETURNING id
`

// Store persists a new domain to the database.
func (r *Repository) Store(ctx context.Context, domain *domain.Domain) error {
	rows, err := r.db.Query(ctx, queryStore, pgx.NamedArgs{
		"name":       domain.Name,
		"created_by": domain.CreatedBy,
		"updated_by": domain.UpdatedBy,
	})
	if err != nil {
		return err
	}

	// Scan returning ID
	var id string
	if _, err := pgx.ForEachRow(rows, []any{&id}, func() error { return nil }); err != nil {
		return err
	}
	domain.Id = id
	return nil
}
