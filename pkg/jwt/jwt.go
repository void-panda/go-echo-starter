package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"go-echo-starter/internal/config"
	"go-echo-starter/internal/domain"
)

// Common errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims represents JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// JWT handles JWT operations
type JWT struct {
	secretKey  []byte
	expireTime time.Duration
}

// New creates a new JWT instance
func New(cfg *config.JWTConfig) *JWT {
	return &JWT{
		secretKey:  []byte(cfg.Secret),
		expireTime: cfg.ExpireTime,
	}
}

// Generate generates a new JWT token for a user
func (j *JWT) Generate(user *domain.User) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expireTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// Validate validates a JWT token and returns the claims
func (j *JWT) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetExpireTime returns the token expiration time in seconds
func (j *JWT) GetExpireTime() int64 {
	return int64(j.expireTime.Seconds())
}
