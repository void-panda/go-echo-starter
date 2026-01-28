package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator
type Validator struct {
	validator *validator.Validate
}

// New creates a new Validator instance
func New() *Validator {
	v := validator.New()
	return &Validator{validator: v}
}

// Validate validates a struct
func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return formatValidationErrors(err)
	}
	return nil
}

// formatValidationErrors formats validation errors into a user-friendly message
func formatValidationErrors(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errMsgs []string
		for _, e := range validationErrors {
			errMsgs = append(errMsgs, formatFieldError(e))
		}
		return fmt.Errorf("%s", strings.Join(errMsgs, "; "))
	}
	return err
}

// formatFieldError formats a single field error
func formatFieldError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// RegisterCustomValidation registers a custom validation function
func (v *Validator) RegisterCustomValidation(tag string, fn validator.Func) error {
	return v.validator.RegisterValidation(tag, fn)
}
