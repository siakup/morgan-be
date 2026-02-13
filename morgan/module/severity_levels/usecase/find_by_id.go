package usecase

import (
	"context"

	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
)

func (u *UseCase) FindByID(ctx context.Context, id string) (*domain.SeverityLevel, error) {
	return u.repo.FindByID(ctx, id)
}
