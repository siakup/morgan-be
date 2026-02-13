package usecase

import (
	"context"

	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
)

func (u *UseCase) FindAll(ctx context.Context, filter domain.SeverityLevelFilter) ([]*domain.SeverityLevel, int64, error) {
	return u.repository.FindAll(ctx, filter)
}
