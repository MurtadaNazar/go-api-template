package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserType defines the type of user
// swagger:enum UserType
type UserType string

const (
	// UserTypeRegular regular user type
	UserTypeRegular UserType = "user"
	// UserTypeAdmin admin user type
	UserTypeAdmin UserType = "admin"
)

// User represents the user entity
// swagger:model User
type User struct {
	// ID is the unique identifier for the user
	// example: 123e4567-e89b-12d3-a456-426614174000
	// format: uuid
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	// First name of the user
	// example: John
	// required: true
	// min length: 1
	// max length: 100
	FirstName string `gorm:"size:100;not null" json:"first_name"`

	// Middle name of the user (optional)
	// example: Michael
	// max length: 100
	SecondName string `gorm:"size:100" json:"second_name,omitempty"`

	// Last name of the user
	// example: Doe
	// required: true
	// min length: 1
	// max length: 100
	LastName string `gorm:"size:100;not null" json:"last_name"`

	// Username for authentication and display
	// example: johndoe123
	// required: true
	// min length: 3
	// max length: 50
	// pattern: ^[a-zA-Z0-9_]+$
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"`

	// Email address of the user
	// example: john.doe@example.com
	// required: true
	// format: email
	// max length: 100
	Email string `gorm:"size:100;uniqueIndex;not null" json:"email"`

	// Password for authentication (never exposed in JSON responses)
	// required: true
	// min length: 8
	// writeOnly: true
	Password string `gorm:"not null" json:"-"`

	// Type of user account
	// enum: user,admin
	// example: user
	// default: user
	UserType UserType `gorm:"type:varchar(20);default:'user'" json:"user_type"`

	// Account status
	// enum: active,inactive,suspended
	// example: active
	// default: active
	Status string `gorm:"type:varchar(20);default:'active'" json:"status"`

	// CreatedAt indicates when the user account was created
	// example: 2023-10-05T14:30:00Z
	// format: date-time
	// readOnly: true
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// UpdatedAt shows when the user account was last updated
	// example: 2023-10-05T14:30:00Z
	// format: date-time
	// readOnly: true
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate hook to generate UUID before inserting
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	// Note: Functional indexes are created during migrations, not here
	return nil
}

// FullName returns the full name of the user
func (u *User) FullName() string {
	parts := []string{u.FirstName, u.SecondName, u.LastName}

	var cleanedParts []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			cleanedParts = append(cleanedParts, p)
		}
	}

	return strings.Join(cleanedParts, " ")
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.UserType == UserTypeAdmin
}

// TableName sets the insert table name for this struct type
func (User) TableName() string {
	return "users"
}

// CreateFunctionalIndexes ensures case-insensitive search index
func CreateFunctionalIndexes(db *gorm.DB) error {
	return db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_lower_username
		ON users (LOWER(username));
	`).Error
}
