package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/service"
	"go-echo-starter/pkg/logger"
	"go-echo-starter/pkg/response"
	"go-echo-starter/pkg/validator"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService service.UserService
	validator   *validator.Validator
	log         *logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, v *validator.Validator, log *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   v,
		log:         log,
	}
}

// Create godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body domain.CreateUserRequest true "User details"
// @Success 201 {object} response.Response{data=domain.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c echo.Context) error {
	var req domain.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		h.log.Warn().Err(err).Msg("Failed to bind create user request")
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err)
	}

	user, err := h.userService.Create(c.Request().Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			return response.Error(c, http.StatusConflict, "Email already exists")
		}
		return response.Error(c, http.StatusInternalServerError, "Failed to create user")
	}

	return response.Success(c, http.StatusCreated, "User created successfully", user)
}

// GetByID godoc
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response{data=domain.UserResponse}
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.userService.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.Error(c, http.StatusNotFound, "User not found")
		}
		return response.Error(c, http.StatusInternalServerError, "Failed to get user")
	}

	return response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

// GetAll godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]domain.UserResponse}
// @Failure 500 {object} response.Response
// @Router /api/v1/users [get]
func (h *UserHandler) GetAll(c echo.Context) error {
	users, err := h.userService.GetAll(c.Request().Context())
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to get users")
	}

	return response.Success(c, http.StatusOK, "Users retrieved successfully", users)
}

// Update godoc
// @Summary Update a user
// @Description Update a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body domain.UpdateUserRequest true "User details"
// @Success 200 {object} response.Response{data=domain.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid user ID")
	}

	var req domain.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		h.log.Warn().Err(err).Msg("Failed to bind update user request")
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.ValidationError(c, err)
	}

	user, err := h.userService.Update(c.Request().Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.Error(c, http.StatusNotFound, "User not found")
		}
		if errors.Is(err, service.ErrEmailExists) {
			return response.Error(c, http.StatusConflict, "Email already exists")
		}
		return response.Error(c, http.StatusInternalServerError, "Failed to update user")
	}

	return response.Success(c, http.StatusOK, "User updated successfully", user)
}

// Delete godoc
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid user ID")
	}

	err = h.userService.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.Error(c, http.StatusNotFound, "User not found")
		}
		return response.Error(c, http.StatusInternalServerError, "Failed to delete user")
	}

	return response.Success(c, http.StatusOK, "User deleted successfully", nil)
}
