package errors

import "net/http"

// ErrorType represents the category of error
type ErrorType string

const (
	ValidationError    ErrorType = "VALIDATION"
	NotFoundError      ErrorType = "NOT_FOUND"
	ConflictError      ErrorType = "CONFLICT"
	UnauthorizedError  ErrorType = "UNAUTHORIZED"
	ForbiddenError     ErrorType = "FORBIDDEN"
	InternalError      ErrorType = "INTERNAL"
	BadRequestError    ErrorType = "BAD_REQUEST"
	AlreadyExistsError ErrorType = "ALREADY_EXISTS"
)

// AppError is the unified error type for the application
type AppError struct {
	Type       ErrorType `json:"type"`
	Message    string    `json:"message"`
	HTTPStatus int       `json:"-"` // Not exposed in JSON
	Details    string    `json:"details,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError with default HTTP status
func NewAppError(errType ErrorType, message string) *AppError {
	return &AppError{
		Type:       errType,
		Message:    message,
		HTTPStatus: mapErrorTypeToStatus(errType),
	}
}

// NewAppErrorWithDetails creates a new AppError with details
func NewAppErrorWithDetails(errType ErrorType, message string, details string) *AppError {
	return &AppError{
		Type:       errType,
		Message:    message,
		HTTPStatus: mapErrorTypeToStatus(errType),
		Details:    details,
	}
}

// mapErrorTypeToStatus maps error types to HTTP status codes
func mapErrorTypeToStatus(errType ErrorType) int {
	switch errType {
	case ValidationError, BadRequestError:
		return http.StatusBadRequest
	case NotFoundError:
		return http.StatusNotFound
	case ConflictError, AlreadyExistsError:
		return http.StatusConflict
	case UnauthorizedError:
		return http.StatusUnauthorized
	case ForbiddenError:
		return http.StatusForbidden
	case InternalError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}

// Predefined errors for err113 compliance
var (
	ErrTokenNotFound          = NewAppError(NotFoundError, "token not found")
	ErrTokenNotFoundExpired   = NewAppError(NotFoundError, "token not found or expired")
	ErrInvalidToken           = NewAppError(UnauthorizedError, "invalid token")
	ErrInvalidRefreshToken    = NewAppError(UnauthorizedError, "invalid or expired refresh token")
	ErrUsernameAlreadyTaken   = NewAppError(ConflictError, "username already taken")
	ErrEmailAlreadyRegistered = NewAppError(ConflictError, "email already registered")
	ErrUserNotFound           = NewAppError(NotFoundError, "user not found")
	ErrDatabaseError          = NewAppError(InternalError, "database error")
	ErrInvalidFileExtension   = NewAppError(ValidationError, "file must have a valid extension")
	ErrUnsupportedFileType    = NewAppError(ValidationError, "unsupported file type")
)
