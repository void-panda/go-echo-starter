package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"go-echo-starter/internal/domain"
)

// ErrNotFound is returned when a record is not found
var ErrNotFound = errors.New("record not found")

// ErrDuplicateEmail is returned when email already exists
var ErrDuplicateEmail = errors.New("email already exists")

type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowxContext(ctx, query, user.Name, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

// GetByID gets a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetByEmail gets a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.GetContext(ctx, user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetAll gets all users
func (r *userRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY id DESC`

	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2
		WHERE id = $3
		RETURNING updated_at
	`

	err := r.db.QueryRowxContext(ctx, query, user.Name, user.Email, user.ID).
		Scan(&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		if isDuplicateKeyError(err) {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

// Delete deletes a user
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// isDuplicateKeyError checks if error is a duplicate key violation
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "duplicate key") || contains(errStr, "23505")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsSubstring(s, substr)
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
