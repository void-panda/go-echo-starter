package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/repository"
	"go-echo-starter/pkg/jwt"
	"go-echo-starter/pkg/logger"
)

// Common auth errors
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// AuthService defines the interface for authentication
type AuthService interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.TokenResponse, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwt      *jwt.JWT
	log      *logger.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwt *jwt.JWT, log *logger.Logger) AuthService {
	return &authService{
		userRepo: userRepo,
		jwt:      jwt,
		log:      log,
	}
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.TokenResponse, error) {
	// Check if email already exists
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, repository.ErrNotFound) {
		s.log.Error().Err(err).Msg("Failed to check existing email")
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to hash password")
		return nil, err
	}

	// Create user
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.log.Error().Err(err).Msg("Failed to create user")
		return nil, err
	}

	// Generate token
	token, err := s.jwt.Generate(user)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to generate token")
		return nil, err
	}

	s.log.Info().Str("user_id", user.ID.String()).Msg("User registered successfully")

	return &domain.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.jwt.GetExpireTime(),
	}, nil
}

// Login authenticates a user
func (s *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		s.log.Error().Err(err).Msg("Failed to get user by email")
		return nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := s.jwt.Generate(user)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to generate token")
		return nil, err
	}

	s.log.Info().Str("user_id", user.ID.String()).Msg("User logged in successfully")

	return &domain.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.jwt.GetExpireTime(),
	}, nil
}
