package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/domains/domain"
)

// DomainsUseCaseMock is a mock for Domains UseCase
type DomainsUseCaseMock struct {
	mock.Mock
}

func (m *DomainsUseCaseMock) FindAll(ctx context.Context, filter domain.DomainFilter) ([]*domain.Domain, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Domain), args.Get(1).(int64), args.Error(2)
}

func (m *DomainsUseCaseMock) Get(ctx context.Context, id string) (*domain.Domain, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Domain), args.Error(1)
}

func (m *DomainsUseCaseMock) Create(ctx context.Context, domain *domain.Domain) error {
	args := m.Called(ctx, domain)
	return args.Error(0)
}

func (m *DomainsUseCaseMock) Update(ctx context.Context, domain *domain.Domain) error {
	args := m.Called(ctx, domain)
	return args.Error(0)
}

func (m *DomainsUseCaseMock) Delete(ctx context.Context, id string, deletedBy string) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

// DomainsRepositoryMock is a mock for Domains Repository
type DomainsRepositoryMock struct {
	mock.Mock
}

func (m *DomainsRepositoryMock) FindAll(ctx context.Context, filter domain.DomainFilter) ([]*domain.Domain, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Domain), args.Get(1).(int64), args.Error(2)
}

func (m *DomainsRepositoryMock) FindByID(ctx context.Context, id string) (*domain.Domain, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Domain), args.Error(1)
}

func (m *DomainsRepositoryMock) Store(ctx context.Context, domain *domain.Domain) error {
	args := m.Called(ctx, domain)
	return args.Error(0)
}

func (m *DomainsRepositoryMock) Update(ctx context.Context, domain *domain.Domain) error {
	args := m.Called(ctx, domain)
	return args.Error(0)
}

func (m *DomainsRepositoryMock) Delete(ctx context.Context, id string, deletedBy string) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}
