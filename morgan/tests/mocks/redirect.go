package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/siakup/morgan-be/morgan/module/redirect/domain"
)

// RedirectUseCaseMock is a mock for RedirectUseCase
type RedirectUseCaseMock struct {
	mock.Mock
}

func (m *RedirectUseCaseMock) Redirect(ctx context.Context, institutionId string, token string) (string, string, error) {
	args := m.Called(ctx, institutionId, token)
	return args.String(0), args.String(1), args.Error(2)
}

// RedirectRepositoryMock is a mock for RedirectRepository
type RedirectRepositoryMock struct {
	mock.Mock
}

func (m *RedirectRepositoryMock) FindInstitutionByID(ctx context.Context, id string) (*domain.Institution, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Institution), args.Error(1)
}

func (m *RedirectRepositoryMock) FindUserBySub(ctx context.Context, institutionId string, sub string) (*domain.User, error) {
	args := m.Called(ctx, institutionId, sub)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *RedirectRepositoryMock) StoreSession(ctx context.Context, session *domain.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}
