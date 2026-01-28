package service

import (
	"context"
	"errors"

	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/repository"
	"go-echo-starter/pkg/logger"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
	ErrInvalidInput = errors.New("invalid input")
)

type userService struct {
	userRepo repository.UserRepository
	log      *logger.Logger
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, log *logger.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		log:      log,
	}
}

// Create creates a new user
func (s *userService) Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error) {
	user := &domain.User{
		Name:  req.Name,
		Email: req.Email,
	}

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			s.log.Warn().Str("email", req.Email).Msg("Attempted to create user with existing email")
			return nil, ErrEmailExists
		}
		s.log.Error().Err(err).Msg("Failed to create user")
		return nil, err
	}

	s.log.Info().Str("user_id", user.ID.String()).Msg("User created successfully")
	return user.ToResponse(), nil
}

// GetByID gets a user by ID
func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		s.log.Error().Err(err).Str("user_id", id.String()).Msg("Failed to get user")
		return nil, err
	}

	return user.ToResponse(), nil
}

// GetAll gets all users
func (s *userService) GetAll(ctx context.Context) ([]*domain.UserResponse, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get users")
		return nil, err
	}

	responses := make([]*domain.UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToResponse()
	}

	return responses, nil
}

// Update updates a user
func (s *userService) Update(ctx context.Context, id uuid.UUID, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		s.log.Error().Err(err).Str("user_id", id.String()).Msg("Failed to get user for update")
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, ErrEmailExists
		}
		s.log.Error().Err(err).Str("user_id", id.String()).Msg("Failed to update user")
		return nil, err
	}

	s.log.Info().Str("user_id", id.String()).Msg("User updated successfully")
	return user.ToResponse(), nil
}

// Delete deletes a user
func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		s.log.Error().Err(err).Str("user_id", id.String()).Msg("Failed to delete user")
		return err
	}

	s.log.Info().Str("user_id", id.String()).Msg("User deleted successfully")
	return nil
}
