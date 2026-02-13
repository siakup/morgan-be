package mocks

import (
	"context"

	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
	"github.com/stretchr/testify/mock"
)

// SeverityLevelsUseCaseMock is a mock implementation of domain.UseCase
type SeverityLevelsUseCaseMock struct {
	mock.Mock
}

func (m *SeverityLevelsUseCaseMock) FindAll(ctx context.Context, filter domain.SeverityLevelFilter) ([]*domain.SeverityLevel, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, int64(0), args.Error(2)
	}
	return args.Get(0).([]*domain.SeverityLevel), args.Get(1).(int64), args.Error(2)
}

func (m *SeverityLevelsUseCaseMock) Get(ctx context.Context, id string) (*domain.SeverityLevel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SeverityLevel), args.Error(1)
}

func (m *SeverityLevelsUseCaseMock) Create(ctx context.Context, severityLevel *domain.SeverityLevel) error {
	args := m.Called(ctx, severityLevel)
	return args.Error(0)
}

func (m *SeverityLevelsUseCaseMock) Update(ctx context.Context, severityLevel *domain.SeverityLevel) error {
	args := m.Called(ctx, severityLevel)
	return args.Error(0)
}

func (m *SeverityLevelsUseCaseMock) Delete(ctx context.Context, id string, deletedBy string) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}
