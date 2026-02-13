package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/domain"
)

// ShiftSessionsUseCaseMock is a mock for ShiftSessions UseCase
type ShiftSessionsUseCaseMock struct {
	mock.Mock
}

func (m *ShiftSessionsUseCaseMock) FindAll(ctx context.Context, filter domain.ShiftSessionFilter) ([]*domain.ShiftSession, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.ShiftSession), args.Get(1).(int64), args.Error(2)
}

func (m *ShiftSessionsUseCaseMock) Get(ctx context.Context, id string) (*domain.ShiftSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ShiftSession), args.Error(1)
}

func (m *ShiftSessionsUseCaseMock) Create(ctx context.Context, shiftSession *domain.ShiftSession) error {
	args := m.Called(ctx, shiftSession)
	return args.Error(0)
}

func (m *ShiftSessionsUseCaseMock) Update(ctx context.Context, shiftSession *domain.ShiftSession) error {
	args := m.Called(ctx, shiftSession)
	return args.Error(0)
}

func (m *ShiftSessionsUseCaseMock) Delete(ctx context.Context, id string, deletedBy string) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

// ShiftSessionsRepositoryMock is a mock for ShiftSessions Repository
type ShiftSessionsRepositoryMock struct {
	mock.Mock
}

func (m *ShiftSessionsRepositoryMock) FindAll(ctx context.Context, filter domain.ShiftSessionFilter) ([]*domain.ShiftSession, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.ShiftSession), args.Get(1).(int64), args.Error(2)
}

func (m *ShiftSessionsRepositoryMock) FindByID(ctx context.Context, id string) (*domain.ShiftSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ShiftSession), args.Error(1)
}

func (m *ShiftSessionsRepositoryMock) Store(ctx context.Context, shiftSession *domain.ShiftSession) error {
	args := m.Called(ctx, shiftSession)
	return args.Error(0)
}

func (m *ShiftSessionsRepositoryMock) Update(ctx context.Context, shiftSession *domain.ShiftSession) error {
	args := m.Called(ctx, shiftSession)
	return args.Error(0)
}

func (m *ShiftSessionsRepositoryMock) Delete(ctx context.Context, id string, deletedBy string) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}
