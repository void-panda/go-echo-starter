package domain

import (
	"time"

	"github.com/google/uuid"
)

// Auth represents authentication credentials
type Auth struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

// TokenResponse represents the token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// AuthUser represents authenticated user data in JWT claims
type AuthUser struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
