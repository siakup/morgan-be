package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/redirect/domain"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindInstitutionByID(ctx context.Context, id string) (*domain.Institution, error) {
	sql := `
		SELECT
			id, settings
		FROM auth.institutions
		WHERE id = $1
	`

	var inst domain.Institution
	// Using manual scan for safety/simplicity given uncertainties about tags
	if err := r.db.QueryRow(ctx, sql, id).Scan(&inst.Id, &inst.Settings); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("institution not found")
		}
		return nil, err
	}

	return &inst, nil
}

func (r *Repository) FindUserBySub(ctx context.Context, institutionId string, sub string) (*domain.User, error) {
	sql := `
		SELECT
		    u.id,
		    u.institution_id,
		    u.identity_provider,
		    u.external_subject,
		    u.metadata,
		    jsonb_agg(
		        jsonb_build_object(
		            'role_id', r.id,
		            'role_name', r.name,
		            'groups', (
		                SELECT array_agg(DISTINCT ur2.group_id)
		                FROM iam.user_roles ur2
		                WHERE ur2.user_id = u.id
		                  AND ur2.role_id = r.id
		            ),
		            'permissions', (
                        SELECT array_agg(DISTINCT p.code)
                        FROM iam.role_permissions rp
                        JOIN iam.permissions p ON p.id = rp.permission_id
                        WHERE rp.role_id = r.id
                    )
		        )
		    ) AS roles
		FROM auth.users u
		JOIN iam.user_roles ur ON u.id = ur.user_id
		JOIN iam.roles r ON ur.role_id = r.id
		WHERE u.external_subject = @subject
		AND u.institution_id = @institution_id
		AND u.deleted_at IS NULL
		GROUP BY u.id, u.external_subject
		LIMIT 1
	`

	args := pgx.NamedArgs{
		"subject":        sub,
		"institution_id": institutionId,
	}

	rows, err := r.db.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("user not found")
		}
		return nil, err
	}

	user.InstitutionId = institutionId
	return user, nil
}

func (r *Repository) StoreSession(ctx context.Context, session *domain.Session) error {
	sql := `
		INSERT INTO auth.sessions (
			session_id, institution_id, user_id, external_subject, roles, access_token, expires_at
		) VALUES (
			@session_id, @institution_id, @user_id, @external_subject, @roles, @access_token, @expires_at
		)
	`
	args := pgx.NamedArgs{
		"session_id":       session.SessionId,
		"institution_id":   session.InstitutionId,
		"user_id":          session.UserId,
		"external_subject": session.ExternalSubject,
		"roles":            session.Roles,
		"access_token":     session.AccessToken,
		"expires_at":       session.ExpiresAt,
	}

	_, err := r.db.Exec(ctx, sql, args)
	return err
}
