package usecase

import (
	"context"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/severity_levels/domain"
)

func (u *UseCase) FindByID(ctx context.Context, id string) (*domain.SeverityLevel, error) {
	return u.repo.FindByID(ctx, id)
}
