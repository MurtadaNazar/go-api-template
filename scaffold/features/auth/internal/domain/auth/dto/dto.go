package dto

// LoginRequest represents the payload for user login
// swagger:model
type LoginRequest struct {
	// Email or username of the user
	// Required: true
	// Example: john.doe@example.com
	EmailOrUsername string `json:"email_or_username" validate:"required,min=3,max=255"`

	// Password for the user account
	// Required: true
	// Example: StrongP@ssw0rd
	Password string `json:"password" validate:"required,min=8"`
}

// LoginResponse represents the response after successful login
// swagger:model
type LoginResponse struct {
	// Access token for API authentication
	// Example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	AccessToken string `json:"access_token"`

	// Refresh token for obtaining new access tokens
	// Example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	RefreshToken string `json:"refresh_token"`

	// Type of the token
	// Example: Bearer
	TokenType string `json:"token_type" example:"Bearer"`

	// Expiration time in seconds
	// Example: 900
	ExpiresIn int64 `json:"expires_in" example:"900"`
}

// RefreshTokenRequest represents the payload for refreshing an access token
// swagger:model
type RefreshTokenRequest struct {
	// Refresh token obtained during login
	// Required: true
	// Example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	RefreshToken string `json:"refresh_token" validate:"required,min=10"`
}

// ErrorResponse represents an error response
// swagger:model
type ErrorResponse struct {
	// Error message
	// Example: authentication failed
	Error string `json:"error" example:"authentication failed"`

	// Detailed error information
	// Example: invalid credentials
	Details string `json:"details,omitempty" example:"invalid credentials"`
}
