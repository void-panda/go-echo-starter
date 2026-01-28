package service

import (
	"context"
	"errors"
	"testing"

	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/repository"
	"go-echo-starter/pkg/logger"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserService_Create(t *testing.T) {
	log := logger.New("debug", true)

	t.Run("success", func(t *testing.T) {
		repo := new(MockUserRepository)
		svc := NewUserService(repo, log)

		req := &domain.CreateUserRequest{
			Name:  "Test User",
			Email: "test@example.com",
		}

		repo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.Name == req.Name && u.Email == req.Email
		})).Return(nil)

		res, err := svc.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, req.Name, res.Name)
		assert.Equal(t, req.Email, res.Email)
		repo.AssertExpectations(t)
	})

	t.Run("duplicate email", func(t *testing.T) {
		repo := new(MockUserRepository)
		svc := NewUserService(repo, log)

		req := &domain.CreateUserRequest{
			Name:  "Test User",
			Email: "test@example.com",
		}

		repo.On("Create", mock.Anything, mock.Anything).Return(repository.ErrDuplicateEmail)

		res, err := svc.Create(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.True(t, errors.Is(err, ErrEmailExists))
		repo.AssertExpectations(t)
	})
}
