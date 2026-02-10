package usecase

import (
	"context"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/severity_levels/domain"
)

func (u *UseCase) FindAll(ctx context.Context, filter domain.SeverityLevelFilter) ([]*domain.SeverityLevel, int64, error) {
	return u.repo.FindAll(ctx, filter)
}
