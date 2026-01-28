package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"go-echo-starter/internal/domain"
	"go-echo-starter/pkg/jwt"
	"go-echo-starter/pkg/response"
)

// JWTAuth creates a JWT authentication middleware
func JWTAuth(jwtService *jwt.JWT) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, http.StatusUnauthorized, "Missing authorization header")
			}

			// Check Bearer prefix
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Error(c, http.StatusUnauthorized, "Invalid authorization header format")
			}

			// Validate token
			claims, err := jwtService.Validate(parts[1])
			if err != nil {
				return response.Error(c, http.StatusUnauthorized, "Invalid or expired token")
			}

			// Set user in context
			user := &domain.AuthUser{
				ID:    claims.UserID,
				Name:  claims.Name,
				Email: claims.Email,
			}
			c.Set("user", user)

			return next(c)
		}
	}
}
