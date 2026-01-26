package middleware

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/idp"
)

const (
	authCookieKey   = "access_token"
	PrefixAuthToken = "auth:token:"
)

type (
	AuthorizationMiddleware struct {
		idp   idp.IDPProvider
		db    *pgxpool.Pool
		cache redis.UniversalClient
	}
)

func NewAuthorizationMiddleware(idp idp.IDPProvider, db *pgxpool.Pool, cache redis.UniversalClient) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{idp: idp, db: db, cache: cache}
}

const (
	XTokenKey        = "X-Token"
	XUserIdKey       = "X-User-Id"
	XExternalSubject = "X-External-Subject"
	XGroupKey        = "X-Group"
	XInstitutionId   = "X-Institution-Id"
)

func (a *AuthorizationMiddleware) Authenticate(scopes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authKey := c.Cookies(authCookieKey)

		ctx := c.UserContext()
		logger := zerolog.Ctx(ctx).With().Str("component", "middleware.auth").Logger()

		if authKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized: Missing Token"})
		}

		auth, err := a.findFromCache(ctx, authKey)
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				logger.Error().Err(err).Msg("failed to get auth from cache")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
			}

			auth, err = a.findFromDB(ctx, authKey)
			if err != nil {
				// Don't log normal "not found" as error if it's just invalid session
				if !errors.Is(err, pgx.ErrNoRows) {
					logger.Error().Err(err).Msg("failed to get auth from db")
				}
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
			}

			if auth.ExpiresAt.Before(time.Now()) {
				logger.Info().Str("identity_provider", auth.IdentityProvider).Msg("token expired, attempting refresh")

				idpClient, err := a.idp.GetIDP(ctx, auth.IdentityProvider)
				if err != nil {
					logger.Error().Err(err).Str("idp", auth.IdentityProvider).Msg("failed to get IDP client for refresh")
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
				}

				newAccess, err := idpClient.Refresh(ctx, auth.AccessToken)
				if nil != err {
					logger.Warn().Err(err).Msg("failed to refresh token, removing session")
					_ = a.removeOnDB(ctx, authKey)
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
				}

				auth.AccessToken = newAccess.AccessToken
				expiryDuration := time.Duration(newAccess.ExpiresIn-30) * time.Second
				auth.ExpiresAt = time.Now().Add(expiryDuration)

				if err = a.updateOnDB(ctx, auth); nil != err {
					logger.Error().Err(err).Msg("failed to update session in db after refresh")
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
				}
			}

			if err = a.putOnCache(ctx, auth); nil != err {
				logger.Warn().Err(err).Msg("failed to update session cache after refresh")
			}
		}

		for _, scope := range scopes {
			if !slices.Contains(auth.Permissions(), scope) {
				logger.Warn().Str("required_scope", scope).Str("user_id", auth.UserId).Msg("insufficient permissions")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized: Insufficient Permissions"})
			}
		}

		// Context/Locals values (fiber way to pass data)
		//get user id by seasion id
		c.Locals(XTokenKey, auth.AccessToken)
		c.Locals(XUserIdKey, auth.UserId)
		c.Locals(XExternalSubject, auth.ExternalSubject)
		c.Locals(XInstitutionId, auth.InstitutionId)
		c.Locals(XGroupKey, auth.Groups())

		return c.Next()
	}
}

func (a *AuthorizationMiddleware) findFromCache(ctx context.Context, key string) (*UserRoles, error) {
	var result UserRoles
	if err := a.cache.Get(ctx, PrefixAuthToken+key).Scan(&result); nil != err {
		return nil, err
	}

	return &result, nil
}

func (a *AuthorizationMiddleware) findFromDB(ctx context.Context, key string) (*UserRoles, error) {
	const query = `
		select
		    session_id, institution_id, identity_provider, user_id, external_subject, roles, access_token, expires_at
		from auth.sessions
		where session_id=@session_id
		limit 1
	`

	rows, err := a.db.Query(ctx, query, pgx.NamedArgs{"session_id": key})
	if err != nil {
		return nil, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[UserRoles])
}

func (a *AuthorizationMiddleware) removeOnDB(ctx context.Context, sessionId string) error {
	const query = `DELETE FROM auth.sessions WHERE session_id=@session_id`
	_, err := a.db.Exec(ctx, query, pgx.NamedArgs{"session_id": sessionId})

	return err
}

func (a *AuthorizationMiddleware) updateOnDB(ctx context.Context, auth *UserRoles) error {
	const query = `UPDATE auth.sessions SET access_token=@access_token, expires_at=@expires_at WHERE session_id=@session_id`
	_, err := a.db.Exec(ctx, query, pgx.NamedArgs{"access_token": auth.AccessToken, "expires_at": auth.ExpiresAt, "session_id": auth.SessionId})

	return err
}

func (a *AuthorizationMiddleware) putOnCache(ctx context.Context, auth *UserRoles) error {
	return a.cache.Set(ctx, PrefixAuthToken+auth.SessionId, auth, time.Until(auth.ExpiresAt)).Err()
}
