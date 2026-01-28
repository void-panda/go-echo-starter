package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-echo-starter/internal/domain"
	"go-echo-starter/pkg/logger"
	"go-echo-starter/pkg/validator"
)

type MockUserServiceReal struct {
	mock.Mock
}

func (m *MockUserServiceReal) Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

func (m *MockUserServiceReal) GetByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

func (m *MockUserServiceReal) GetAll(ctx context.Context) ([]*domain.UserResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.UserResponse), args.Error(1)
}

func (m *MockUserServiceReal) Update(ctx context.Context, id uuid.UUID, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserResponse), args.Error(1)
}

func (m *MockUserServiceReal) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserHandler_Create(t *testing.T) {
	e := echo.New()
	v := validator.New()
	log := logger.New("debug", true)

	t.Run("success", func(t *testing.T) {
		mockSvc := new(MockUserServiceReal)
		h := NewUserHandler(mockSvc, v, log)

		userJSON := `{"name":"John Doe","email":"john@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expectedRes := &domain.UserResponse{
			ID:    uuid.New(),
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockSvc.On("Create", mock.Anything, mock.Anything).Return(expectedRes, nil)

		if assert.NoError(t, h.Create(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			var res map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &res)
			assert.Equal(t, "User created successfully", res["message"])

			data := res["data"].(map[string]interface{})
			assert.NotEmpty(t, data["id"])
			assert.Equal(t, "John Doe", data["name"])
		}
	})
}
