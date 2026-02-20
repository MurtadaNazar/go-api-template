package testutil

import (
	"context"
	"go_platform_template/internal/domain/user/model"
	"go_platform_template/internal/domain/user/repo"
	"time"

	"github.com/google/uuid"
)

// MockUserRepo is a mock implementation of UserRepo for testing
type MockUserRepo struct {
	CreateFn               func(ctx context.Context, user *model.User) error
	FindByIDFn             func(ctx context.Context, id string) (*model.User, error)
	UpdateFn               func(ctx context.Context, user *model.User) error
	DeleteFn               func(ctx context.Context, id string) error
	GetByEmailFn           func(ctx context.Context, email string) (*model.User, error)
	FindByUsernameFn       func(ctx context.Context, username string) (*model.User, error)
	GetByEmailOrUsernameFn func(ctx context.Context, emailOrUsername string) (*model.User, error)
	ListFn                 func(ctx context.Context, offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.User, error)
}

// Verify MockUserRepo implements UserRepo interface
var _ repo.UserRepo = (*MockUserRepo)(nil)

func (m *MockUserRepo) Create(ctx context.Context, user *model.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, user)
	}
	return nil
}

func (m *MockUserRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepo) Update(ctx context.Context, user *model.User) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, user)
	}
	return nil
}

func (m *MockUserRepo) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	if m.FindByUsernameFn != nil {
		return m.FindByUsernameFn(ctx, username)
	}
	return nil, nil
}

func (m *MockUserRepo) GetByEmailOrUsername(ctx context.Context, emailOrUsername string) (*model.User, error) {
	if m.GetByEmailOrUsernameFn != nil {
		return m.GetByEmailOrUsernameFn(ctx, emailOrUsername)
	}
	return nil, nil
}

func (m *MockUserRepo) List(ctx context.Context, offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.User, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, offset, limit, filters, sortBy, sortOrder)
	}
	return make([]*model.User, 0), nil
}

// TestUser creates a test user with default values
// Password is hashed with bcrypt (cost 10): "password" => "$2a$10$V.1lMHmJnhH7fB8VXBa5Zeq8N/Ygpg0hW.Qvz5OVSvE9A.6MzUqQm"
func TestUser() *model.User {
	return &model.User{
		ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		Email:      "test@example.com",
		Username:   "testuser",
		Password:   "$2a$10$V.1lMHmJnhH7fB8VXBa5Zeq8N/Ygpg0hW.Qvz5OVSvE9A.6MzUqQm", // bcrypt hash of "password"
		FirstName:  "Test",
		SecondName: "User",
		LastName:   "Account",
		UserType:   model.UserTypeRegular,
		Status:     "active",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// TestUserWithEmail creates a test user with specified email
func TestUserWithEmail(email string) *model.User {
	user := TestUser()
	user.Email = email
	return user
}

// TestUserWithID creates a test user with specified ID
func TestUserWithID(id uuid.UUID) *model.User {
	user := TestUser()
	user.ID = id
	return user
}

// TestUserAdmin creates a test admin user
func TestUserAdmin() *model.User {
	user := TestUser()
	user.UserType = model.UserTypeAdmin
	return user
}
