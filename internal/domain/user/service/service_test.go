package service

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"go_platform_template/internal/domain/user/dto"
	"go_platform_template/internal/domain/user/model"
	apperrors "go_platform_template/internal/shared/errors"
	"go_platform_template/internal/testutil"
)

func TestUserService_Register_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &testutil.MockUserRepo{}
	logger := zap.NewNop().Sugar()
	service := NewUserService(mockRepo, logger)

	req := &dto.UserCreateRequest{
		Email:     "newuser@example.com",
		Username:  "newuser",
		Password:  "password123",
		FirstName: "New",
		LastName:  "User",
		UserType:  string(model.UserTypeRegular),
	}

	mockRepo.FindByUsernameFn = func(ctx context.Context, username string) (*model.User, error) {
		return nil, nil // Username doesn't exist
	}
	mockRepo.GetByEmailFn = func(ctx context.Context, email string) (*model.User, error) {
		return nil, nil // Email doesn't exist
	}
	mockRepo.CreateFn = func(ctx context.Context, user *model.User) error {
		return nil
	}

	// Act
	result, err := service.Register(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("Register() error = %v, want nil", err)
	}
	if result == nil {
		t.Fatal("Register() returned nil user")
	}
	if result.Email != req.Email {
		t.Errorf("Register() email = %s, want %s", result.Email, req.Email)
	}
	if result.Username != req.Username {
		t.Errorf("Register() username = %s, want %s", result.Username, req.Username)
	}
	// Verify password is hashed
	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(req.Password)); err != nil {
		t.Error("Register() password not properly hashed")
	}
}

func TestUserService_Register_DuplicateEmail(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &testutil.MockUserRepo{}
	logger := zap.NewNop().Sugar()
	service := NewUserService(mockRepo, logger)

	req := &dto.UserCreateRequest{
		Email:    "existing@example.com",
		Username: "newuser",
		Password: "password123",
	}

	mockRepo.FindByUsernameFn = func(ctx context.Context, username string) (*model.User, error) {
		return nil, nil
	}
	mockRepo.GetByEmailFn = func(ctx context.Context, email string) (*model.User, error) {
		return testutil.TestUserWithEmail(email), nil // Email exists
	}

	// Act
	result, err := service.Register(ctx, req)

	// Assert
	if err == nil {
		t.Fatal("Register() error = nil, want error for duplicate email")
	}
	if result != nil {
		t.Error("Register() should return nil user on error")
	}
	appErr, ok := apperrors.IsAppError(err)
	if !ok || appErr.Type != apperrors.ConflictError {
		t.Errorf("Register() error type = %v, want ConflictError", appErr.Type)
	}
}

func TestUserService_Register_DuplicateUsername(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &testutil.MockUserRepo{}
	logger := zap.NewNop().Sugar()
	service := NewUserService(mockRepo, logger)

	req := &dto.UserCreateRequest{
		Email:    "new@example.com",
		Username: "existinguser",
		Password: "password123",
	}

	mockRepo.FindByUsernameFn = func(ctx context.Context, username string) (*model.User, error) {
		return testutil.TestUser(), nil // Username exists
	}

	// Act
	result, err := service.Register(ctx, req)

	// Assert
	if err == nil {
		t.Fatal("Register() error = nil, want error for duplicate username")
	}
	if result != nil {
		t.Error("Register() should return nil user on error")
	}
	appErr, ok := apperrors.IsAppError(err)
	if !ok || appErr.Type != apperrors.ConflictError {
		t.Errorf("Register() error type = %v, want ConflictError", appErr.Type)
	}
}

func TestUserService_GetByID_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &testutil.MockUserRepo{}
	logger := zap.NewNop().Sugar()
	service := NewUserService(mockRepo, logger)

	testUser := testutil.TestUser()
	userIDStr := testUser.ID.String()

	mockRepo.FindByIDFn = func(ctx context.Context, id string) (*model.User, error) {
		return testUser, nil
	}

	// Act
	result, err := service.GetByID(ctx, userIDStr)

	// Assert
	if err != nil {
		t.Fatalf("GetByID() error = %v, want nil", err)
	}
	if result == nil {
		t.Fatal("GetByID() returned nil user")
	}
	if result.ID != testUser.ID {
		t.Errorf("GetByID() id = %s, want %s", result.ID, testUser.ID)
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &testutil.MockUserRepo{}
	logger := zap.NewNop().Sugar()
	service := NewUserService(mockRepo, logger)

	mockRepo.FindByIDFn = func(ctx context.Context, id string) (*model.User, error) {
		return nil, nil // Not found
	}

	// Act
	result, err := service.GetByID(ctx, "non-existent-id")

	// Assert
	if err == nil {
		t.Fatal("GetByID() error = nil, want NotFoundError")
	}
	if result != nil {
		t.Error("GetByID() should return nil user when not found")
	}
	appErr, ok := apperrors.IsAppError(err)
	if !ok || appErr.Type != apperrors.NotFoundError {
		t.Errorf("GetByID() error type = %v, want NotFoundError", appErr.Type)
	}
}

func TestUserService_Update_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &testutil.MockUserRepo{}
	logger := zap.NewNop().Sugar()
	service := NewUserService(mockRepo, logger)

	testUser := testutil.TestUser()
	userIDStr := testUser.ID.String()
	updateReq := &dto.UserUpdateRequest{
		FirstName: "Updated",
		LastName:  "Name",
	}

	mockRepo.FindByIDFn = func(ctx context.Context, id string) (*model.User, error) {
		return testUser, nil
	}
	mockRepo.UpdateFn = func(ctx context.Context, user *model.User) error {
		return nil
	}

	// Act
	result, err := service.Update(ctx, userIDStr, updateReq)

	// Assert
	if err != nil {
		t.Fatalf("Update() error = %v, want nil", err)
	}
	if result == nil {
		t.Fatal("Update() returned nil user")
	}
	if result.FirstName != updateReq.FirstName {
		t.Errorf("Update() first_name = %s, want %s", result.FirstName, updateReq.FirstName)
	}
	if result.LastName != updateReq.LastName {
		t.Errorf("Update() last_name = %s, want %s", result.LastName, updateReq.LastName)
	}
}

func TestUserService_Update_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &testutil.MockUserRepo{}
	logger := zap.NewNop().Sugar()
	service := NewUserService(mockRepo, logger)

	mockRepo.FindByIDFn = func(ctx context.Context, id string) (*model.User, error) {
		return nil, nil // Not found
	}

	// Act
	result, err := service.Update(ctx, "non-existent-id", &dto.UserUpdateRequest{})

	// Assert
	if err == nil {
		t.Fatal("Update() error = nil, want NotFoundError")
	}
	if result != nil {
		t.Error("Update() should return nil user when not found")
	}
	appErr, ok := apperrors.IsAppError(err)
	if !ok || appErr.Type != apperrors.NotFoundError {
		t.Errorf("Update() error type = %v, want NotFoundError", appErr.Type)
	}
}

func TestUserService_Delete_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &testutil.MockUserRepo{}
	logger := zap.NewNop().Sugar()
	service := NewUserService(mockRepo, logger)

	testUser := testutil.TestUser()
	userIDStr := testUser.ID.String()

	mockRepo.FindByIDFn = func(ctx context.Context, id string) (*model.User, error) {
		return testUser, nil
	}
	mockRepo.DeleteFn = func(ctx context.Context, id string) error {
		return nil
	}

	// Act
	err := service.Delete(ctx, userIDStr)

	// Assert
	if err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}
}

func TestUserService_Delete_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &testutil.MockUserRepo{}
	logger := zap.NewNop().Sugar()
	service := NewUserService(mockRepo, logger)

	mockRepo.FindByIDFn = func(ctx context.Context, id string) (*model.User, error) {
		return nil, nil // Not found
	}

	// Act
	err := service.Delete(ctx, "non-existent-id")

	// Assert
	if err == nil {
		t.Fatal("Delete() error = nil, want NotFoundError")
	}
	appErr, ok := apperrors.IsAppError(err)
	if !ok || appErr.Type != apperrors.NotFoundError {
		t.Errorf("Delete() error type = %v, want NotFoundError", appErr.Type)
	}
}

func TestUserService_Delete_RepositoryError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &testutil.MockUserRepo{}
	logger := zap.NewNop().Sugar()
	service := NewUserService(mockRepo, logger)

	testUser := testutil.TestUser()
	userIDStr := testUser.ID.String()

	mockRepo.FindByIDFn = func(ctx context.Context, id string) (*model.User, error) {
		return testUser, nil
	}
	mockRepo.DeleteFn = func(ctx context.Context, id string) error {
		return apperrors.ErrDatabaseError
	}

	// Act
	err := service.Delete(ctx, userIDStr)

	// Assert
	if err == nil {
		t.Fatal("Delete() error = nil, want InternalError")
	}
	appErr, ok := apperrors.IsAppError(err)
	if !ok || appErr.Type != apperrors.InternalError {
		t.Errorf("Delete() error type = %v, want InternalError", appErr.Type)
	}
}

func TestUserService_List_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := &testutil.MockUserRepo{}
	logger := zap.NewNop().Sugar()
	service := NewUserService(mockRepo, logger)

	users := []*model.User{testutil.TestUser(), testutil.TestUserAdmin()}
	mockRepo.ListFn = func(ctx context.Context, offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.User, error) {
		return users, nil
	}

	// Act
	result, err := service.List(ctx, 0, 10, nil, "created_at", "asc")

	// Assert
	if err != nil {
		t.Fatalf("List() error = %v, want nil", err)
	}
	if len(result) != len(users) {
		t.Errorf("List() returned %d users, want %d", len(result), len(users))
	}
}
