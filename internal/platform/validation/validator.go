package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	apperrors "go_platform_template/internal/shared/errors"
)

// Validator wraps the playground validator for tag-based validation
type Validator struct {
	validate *validator.Validate
}

// New creates a new Validator instance
func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// ValidateStruct validates a struct using its validation tags
func (v *Validator) ValidateStruct(data interface{}) error {
	if err := v.validate.Struct(data); err != nil {
		return formatValidationError(err)
	}
	return nil
}

// formatValidationError converts validator errors to AppError format
func formatValidationError(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return apperrors.NewAppError(apperrors.ValidationError, "validation failed")
	}

	if len(validationErrors) == 0 {
		return apperrors.NewAppError(apperrors.ValidationError, "validation failed")
	}

	// Build detailed error messages
	var details strings.Builder
	for i, fieldError := range validationErrors {
		if i > 0 {
			details.WriteString("; ")
		}
		details.WriteString(formatFieldError(fieldError))
	}

	return apperrors.NewAppErrorWithDetails(
		apperrors.ValidationError,
		"validation failed",
		details.String(),
	)
}

// formatFieldError formats a single field validation error
func formatFieldError(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()
	value := fe.Value()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "min":
		return fmt.Sprintf("%s must have a minimum length of %s", field, param)
	case "max":
		return fmt.Sprintf("%s must have a maximum length of %s", field, param)
	case "alphanum":
		return fmt.Sprintf("%s must be alphanumeric", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "len":
		return fmt.Sprintf("%s must have exactly %s characters", field, param)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	default:
		return fmt.Sprintf("%s validation failed: %s (value: %v)", field, tag, value)
	}
}
