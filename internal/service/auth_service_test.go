package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"go-echo-starter/internal/config"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/repository"
	"go-echo-starter/pkg/jwt"
	"go-echo-starter/pkg/logger"
)

// MockJWT is a mock implementation of jwt.JWT
// Actually, jwt.JWT is a struct, we might need an interface or just mock it.
// Since jwt.JWT is a struct, I'll mock the internal behavior or just use a real one with a test secret.
// Alternatively, I can create an interface for JWT if it was missing.

func TestAuthService_Register(t *testing.T) {
	log := logger.New("debug", true)
	jwtSvc := jwt.New(&config.JWTConfig{Secret: "test-secret", ExpireTime: 24 * time.Hour})

	t.Run("success", func(t *testing.T) {
		repo := new(MockUserRepository)
		svc := NewAuthService(repo, jwtSvc, log)

		req := &domain.RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		repo.On("GetByEmail", mock.Anything, req.Email).Return(nil, repository.ErrNotFound)
		repo.On("Create", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			return u.Name == req.Name && u.Email == req.Email
		})).Return(nil)

		res, err := svc.Register(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.AccessToken)
		repo.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		repo := new(MockUserRepository)
		svc := NewAuthService(repo, jwtSvc, log)

		req := &domain.RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		repo.On("GetByEmail", mock.Anything, req.Email).Return(&domain.User{ID: uuid.New()}, nil)

		res, err := svc.Register(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.True(t, errors.Is(err, ErrEmailAlreadyExists))
		repo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	log := logger.New("debug", true)
	jwtSvc := jwt.New(&config.JWTConfig{Secret: "test-secret", ExpireTime: 24 * time.Hour})

	t.Run("success", func(t *testing.T) {
		repo := new(MockUserRepository)
		svc := NewAuthService(repo, jwtSvc, log)

		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := &domain.User{
			ID:       uuid.New(),
			Email:    "test@example.com",
			Password: string(hashedPassword),
		}

		req := &domain.LoginRequest{
			Email:    user.Email,
			Password: password,
		}

		repo.On("GetByEmail", mock.Anything, req.Email).Return(user, nil)

		res, err := svc.Login(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.AccessToken)
		repo.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		repo := new(MockUserRepository)
		svc := NewAuthService(repo, jwtSvc, log)

		req := &domain.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		repo.On("GetByEmail", mock.Anything, req.Email).Return(nil, repository.ErrNotFound)

		res, err := svc.Login(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.True(t, errors.Is(err, ErrInvalidCredentials))
		repo.AssertExpectations(t)
	})
}
