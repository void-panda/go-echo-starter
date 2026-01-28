package handler

import (
	"go-echo-starter/internal/service"
	"go-echo-starter/pkg/logger"
	"go-echo-starter/pkg/validator"
)

// Handler holds all HTTP handlers
type Handler struct {
	User      *UserHandler
	Auth      *AuthHandler
	validator *validator.Validator
	log       *logger.Logger
}

// NewHandler creates a new handler
func NewHandler(
	userService service.UserService,
	authService service.AuthService,
	v *validator.Validator,
	log *logger.Logger,
) *Handler {
	return &Handler{
		User:      NewUserHandler(userService, v, log),
		Auth:      NewAuthHandler(authService, v, log),
		validator: v,
		log:       log,
	}
}
