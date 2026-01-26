package domain

import (
	"context"
	"time"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/idp"
)

// Institution represents the institution settings for redirect.
type Institution struct {
	Id          string      `object:"id"`
	RedirectUrl string      `object:"redirect_url"`
	Settings    idp.Setting `object:"settings"`
}

// Session represents an authenticated user session.
type Session struct {
	SessionId       string    `object:"session_id"`
	InstitutionId   string    `object:"institution_id"`
	UserId          string    `object:"user_id"`
	ExternalSubject string    `object:"external_subject"`
	Roles           any       `object:"roles"` // Changed to any to support JSONB structure
	AccessToken     string    `object:"access_token"`
	ExpiresAt       time.Time `object:"expires_at"`
}

// User represents the user found by sub.
type User struct {
	Id               string `object:"id"`
	InstitutionId    string `object:"institution_id"`
	ExternalSubject  string `object:"external_subject"`
	IdentityProvider string `object:"identity_provider"`
	Metadata         any    `object:"metadata"`
	Roles            any    `object:"roles"`
}

// RedirectRepository defines the persistence layer contract.
type RedirectRepository interface {
	FindInstitutionByID(ctx context.Context, id string) (*Institution, error)
	FindUserBySub(ctx context.Context, institutionId string, sub string) (*User, error)
	StoreSession(ctx context.Context, session *Session) error
}
