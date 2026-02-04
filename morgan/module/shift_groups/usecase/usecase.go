package usecase

import (
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

var _ domain.UseCase = (*UseCase)(nil)

// UseCase implements the domain.UseCase interface.
type UseCase struct {
	repo domain.ShiftGroupRepository
}

// NewUseCase creates a new instance of the UseCase.
func NewUseCase(repo domain.ShiftGroupRepository) *UseCase {
	return &UseCase{repo: repo}
}
