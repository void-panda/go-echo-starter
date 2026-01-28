package response

import (
	"github.com/labstack/echo/v4"
)

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Success returns a successful response
func Success(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error returns an error response
func Error(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, Response{
		Success: false,
		Message: message,
	})
}

// ErrorWithDetails returns an error response with additional details
func ErrorWithDetails(c echo.Context, statusCode int, message string, errors interface{}) error {
	return c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// ValidationError returns a validation error response
func ValidationError(c echo.Context, err error) error {
	return c.JSON(400, Response{
		Success: false,
		Message: "Validation failed",
		Errors:  err.Error(),
	})
}
