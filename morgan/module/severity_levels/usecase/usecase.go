package usecase

import (
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/severity_levels/domain"
)

var _ domain.UseCase = (*UseCase)(nil)

// UseCase implements the domain.UseCase interface.
type UseCase struct {
	repo domain.SeverityLevelRepository
}

// NewUseCase creates a new instance of the UseCase.
func NewUseCase(repo domain.SeverityLevelRepository) *UseCase {
	return &UseCase{repo: repo}
}
