package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/service"
	"go-echo-starter/pkg/logger"
	"go-echo-starter/pkg/response"
	"go-echo-starter/pkg/validator"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService service.AuthService
	validator   *validator.Validator
	log         *logger.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService, v *validator.Validator, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   v,
		log:         log,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body domain.RegisterRequest true "Registration details"
// @Success 201 {object} response.Response{data=domain.TokenResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req domain.RegisterRequest
	if err := c.Bind(&req); err != nil {
		h.log.Warn().Err(err).Msg("Failed to bind register request")
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err)
	}

	token, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return response.Error(c, http.StatusConflict, "Email already exists")
		}
		return response.Error(c, http.StatusInternalServerError, "Failed to register user")
	}

	return response.Success(c, http.StatusCreated, "User registered successfully", token)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body domain.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=domain.TokenResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req domain.LoginRequest
	if err := c.Bind(&req); err != nil {
		h.log.Warn().Err(err).Msg("Failed to bind login request")
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err)
	}

	token, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return response.Error(c, http.StatusUnauthorized, "Invalid email or password")
		}
		return response.Error(c, http.StatusInternalServerError, "Failed to login")
	}

	return response.Success(c, http.StatusOK, "Login successful", token)
}

// GetMe godoc
// @Summary Get current user
// @Description Get the currently authenticated user's information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=domain.AuthUser}
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetMe(c echo.Context) error {
	val := c.Get("user")
	if val == nil {
		return response.Error(c, http.StatusUnauthorized, "User context not found")
	}

	user, ok := val.(*domain.AuthUser)
	if !ok {
		return response.Error(c, http.StatusInternalServerError, "Invalid user context")
	}

	return response.Success(c, http.StatusOK, "User retrieved successfully", user)
}
