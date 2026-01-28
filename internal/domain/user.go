package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest represents request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=255"`
	Email string `json:"email" validate:"required,email,max=255"`
}

// UpdateUserRequest represents request body for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=255"`
	Email string `json:"email" validate:"omitempty,email,max=255"`
}

// UserResponse represents user response
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
