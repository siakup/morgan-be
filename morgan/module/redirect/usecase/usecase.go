package usecase

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"github.com/siakup/morgan-be/libraries/errors"
	"github.com/siakup/morgan-be/libraries/idp"
	"github.com/siakup/morgan-be/libraries/types"
	"github.com/siakup/morgan-be/morgan/config"
	"github.com/siakup/morgan-be/morgan/module/redirect/domain"
)

var _ domain.RedirectUseCase = (*UseCase)(nil)

type UseCase struct {
	defaultRedirectUrl string
	repository         domain.RedirectRepository
	idp                idp.IDPProvider
	tracer             trace.Tracer
}

func NewUseCase(app *config.InternalAppConfig, repository domain.RedirectRepository, idp idp.IDPProvider) *UseCase {
	return &UseCase{
		defaultRedirectUrl: app.RedirectUrl,
		repository:         repository,
		idp:                idp,
		tracer:             otel.Tracer("redirect"),
	}
}

func (u *UseCase) Redirect(ctx context.Context, institutionId string, token string) (string, string, error) {
	ctx, span := u.tracer.Start(ctx, "Redirect")
	defer span.End()

	institution, err := u.repository.FindInstitutionByID(ctx, institutionId)
	if err != nil {
		return "", "", err
	}

	idpClient, err := u.idp.GetIDP(ctx, institutionId)
	if err != nil {
		return "", "", errors.InternalServerError("failed to get idp client")
	}

	authSession, err := idpClient.Check(ctx, token)
	if err != nil {
		return "", "", errors.New(errors.ErrorTypeUnauthorized, 403, "invalid token", nil)
	}

	user, err := u.repository.FindUserBySub(ctx, institutionId, authSession.Sub)
	if err != nil {
		return "", "", err
	}

	// Calculate session expiry based on authSession.ExpiresIn (int seconds)
	expiresAt := time.Now().Add(time.Duration(authSession.ExpiresIn) * time.Second)

	session := &domain.Session{
		SessionId:       types.GenerateID(),
		InstitutionId:   institutionId,
		UserId:          user.Id,
		ExternalSubject: user.ExternalSubject,
		Roles:           user.Roles,
		AccessToken:     token,
		ExpiresAt:       expiresAt,
	}

	if err := u.repository.StoreSession(ctx, session); err != nil {
		return "", "", errors.InternalServerError("failed to store session")
	}

	redirectUrl := institution.Settings.Url
	if redirectUrl == "" {
		redirectUrl = u.defaultRedirectUrl
	}

	return redirectUrl, session.SessionId, nil
}
