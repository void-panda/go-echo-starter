package service

import (
	"context"

	"github.com/google/uuid"

	"go-echo-starter/internal/domain"
)

// UserService defines the interface for user business logic
type UserService interface {
	Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error)
	GetAll(ctx context.Context) ([]*domain.UserResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *domain.UpdateUserRequest) (*domain.UserResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
