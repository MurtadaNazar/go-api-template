package dto

// UserCreateRequest represents the payload for creating a user
// swagger:model
type UserCreateRequest struct {
	// FirstName of the user
	// Required: true
	// Example: John
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`

	// SecondName of the user
	// Example: Michael
	SecondName string `json:"second_name" validate:"omitempty,min=2,max=100"`

	// LastName of the user
	// Required: true
	// Example: Doe
	LastName string `json:"last_name" validate:"required,min=2,max=100"`

	// Username of the user
	// Required: true
	// Example: johndoe123
	Username string `json:"username" validate:"required,alphanum,min=3,max=50"`

	// Email of the user
	// Required: true
	// Example: john.doe@example.com
	Email string `json:"email" validate:"required,email"`

	// Password for the user account
	// Required: true
	// Example: StrongP@ssw0rd
	Password string `json:"password" validate:"required,min=8"`

	// UserType defines the role of the user
	// Enum: user, admin
	// Example: user
	UserType string `json:"user_type" validate:"omitempty,oneof=user admin"`
}

// UserUpdateRequest represents the payload for updating a user
// swagger:model
type UserUpdateRequest struct {
	// FirstName of the user
	// Example: John
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=100"`

	// SecondName of the user
	// Example: Michael
	SecondName string `json:"second_name" validate:"omitempty,min=2,max=100"`

	// LastName of the user
	// Example: Doe
	LastName string `json:"last_name" validate:"omitempty,min=2,max=100"`

	// Username of the user
	// Example: johndoe123
	Username string `json:"username" validate:"omitempty,alphanum,min=3,max=50"`

	// Email of the user
	// Example: john.doe@example.com
	Email string `json:"email" validate:"omitempty,email"`

	// Password for the user account
	// Example: StrongP@ssw0rd
	Password string `json:"password" validate:"omitempty,min=8"`

	// UserType defines the role of the user
	// Enum: user, admin
	// Example: user
	UserType string `json:"user_type" validate:"omitempty,oneof=user admin"`
}
