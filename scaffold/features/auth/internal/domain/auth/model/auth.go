package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LoginRequest represents the payload to log in a user
// swagger:model LoginRequest
type LoginRequest struct {
	// User email or username for login
	// required: true
	// example: john.doe@example.com OR johndoe123
	// min length: 3
	// max length: 100
	EmailOrUsername string `json:"email_or_username" binding:"required"`

	// User password for authentication
	// required: true
	// example: mySecurePassword123!
	// min length: 8
	// format: password
	// writeOnly: true
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response after successful login
// swagger:model LoginResponse
type LoginResponse struct {
	// Access JWT token for API authorization
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
	AccessToken string `json:"access_token"`

	// Refresh JWT token for obtaining new access tokens
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
	RefreshToken string `json:"refresh_token"`
}

// RefreshRequest represents the payload to refresh tokens
// swagger:model RefreshRequest
type RefreshRequest struct {
	// Refresh token obtained during login
	// required: true
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshResponse represents the response after token rotation
// swagger:model RefreshResponse
type RefreshResponse struct {
	// New access token for continued API access
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
	AccessToken string `json:"access_token"`

	// New refresh token for future token rotations
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken represents a refresh token stored in DB
// swagger:model RefreshToken
type RefreshToken struct {
	// ID is the unique identifier for the refresh token record
	// example: 123e4567-e89b-12d3-a456-426614174000
	// format: uuid
	// readOnly: true
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	// Token is the hashed refresh token value
	// required: true
	// writeOnly: true
	Token string `gorm:"uniqueIndex;not null" json:"-"`

	// UserID is the UUID of the user who owns this refresh token
	// example: 123e4567-e89b-12d3-a456-426614174000
	// format: uuid
	// required: true
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	// Role of the user for authorization context
	// example: user
	// enum: user,admin
	// required: true
	Role string `gorm:"type:varchar(50);not null" json:"role"`

	// ExpiresAt indicates when the token becomes invalid
	// example: 2023-10-12T14:30:00Z
	// format: date-time
	// required: true
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`

	// IsRevoked indicates if the token has been manually revoked
	// example: false
	// default: false
	IsRevoked bool `gorm:"default:false;index" json:"is_revoked"`

	// CreatedAt indicates when the token was created
	// example: 2023-10-05T14:30:00Z
	// format: date-time
	// readOnly: true
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt shows when the token was last updated
	// example: 2023-10-05T14:30:00Z
	// format: date-time
	// readOnly: true
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the default table name
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
