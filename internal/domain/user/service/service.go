package service

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"go_platform_template/internal/domain/user/dto"
	"go_platform_template/internal/domain/user/model"
	"go_platform_template/internal/domain/user/repo"
	apperrors "go_platform_template/internal/shared/errors"
)

type UserService interface {
	Register(ctx context.Context, req *dto.UserCreateRequest) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, id string, req *dto.UserUpdateRequest) (*model.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.User, error)
}

type userService struct {
	repo   repo.UserRepo
	logger *zap.SugaredLogger
}

func NewUserService(r repo.UserRepo, logger *zap.SugaredLogger) UserService {
	if logger == nil {
		logger = zap.NewNop().Sugar()
	}
	return &userService{
		repo:   r,
		logger: logger,
	}
}

// Register creates a new user with hashed password
func (s *userService) Register(ctx context.Context, req *dto.UserCreateRequest) (*model.User, error) {
	// Ensure username is unique
	if existing, err := s.repo.FindByUsername(ctx, req.Username); err != nil {
		s.logger.Errorw("failed to check username uniqueness", "username", req.Username, "error", err)
		return nil, apperrors.NewAppError(apperrors.InternalError, "Failed to register user")
	} else if existing != nil {
		s.logger.Warnw("duplicate username", "username", req.Username)
		return nil, apperrors.NewAppError(apperrors.ConflictError, "Username already taken")
	}

	// Ensure email is unique
	if existing, err := s.repo.GetByEmail(ctx, req.Email); err != nil {
		s.logger.Errorw("failed to check email uniqueness", "email", req.Email, "error", err)
		return nil, apperrors.NewAppError(apperrors.InternalError, "Failed to register user")
	} else if existing != nil {
		s.logger.Warnw("duplicate email", "email", req.Email)
		return nil, apperrors.NewAppError(apperrors.ConflictError, "Email already registered")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorw("failed to hash password", "error", err)
		return nil, apperrors.NewAppError(apperrors.InternalError, "Failed to register user")
	}

	user := &model.User{
		FirstName:  req.FirstName,
		SecondName: req.SecondName,
		LastName:   req.LastName,
		Username:   req.Username,
		Email:      req.Email,
		Password:   string(hashed),
		UserType:   model.UserType(req.UserType),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Errorw("failed to create user", "username", req.Username, "error", err)
		return nil, apperrors.NewAppError(apperrors.InternalError, "Failed to register user")
	}

	s.logger.Infow("user registered", "user_id", user.ID, "username", user.Username)
	return user, nil
}

// GetByID fetches a user by UUID
func (s *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorw("failed to fetch user", "user_id", id, "error", err)
		return nil, apperrors.NewAppError(apperrors.InternalError, "Failed to fetch user")
	}
	if user == nil {
		s.logger.Warnw("user not found", "user_id", id)
		return nil, apperrors.NewAppError(apperrors.NotFoundError, "User not found")
	}
	return user, nil
}

// Update modifies a user by ID with DTO input
func (s *userService) Update(ctx context.Context, id string, req *dto.UserUpdateRequest) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorw("failed to fetch user for update", "user_id", id, "error", err)
		return nil, apperrors.NewAppError(apperrors.InternalError, "Failed to update user")
	}
	if user == nil {
		s.logger.Warnw("user not found for update", "user_id", id)
		return nil, apperrors.NewAppError(apperrors.NotFoundError, "User not found")
	}

	// Only update fields provided in DTO
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.SecondName != "" {
		user.SecondName = req.SecondName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Errorw("failed to hash password", "user_id", id, "error", err)
			return nil, apperrors.NewAppError(apperrors.InternalError, "Failed to update user password")
		}
		user.Password = string(hashed)
	}
	if req.UserType != "" {
		user.UserType = model.UserType(req.UserType)
	}

	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Errorw("failed to update user", "user_id", id, "error", err)
		return nil, apperrors.NewAppError(apperrors.InternalError, "Failed to update user")
	}

	s.logger.Infow("user updated", "user_id", id)
	return user, nil
}

// Delete removes a user by ID
func (s *userService) Delete(ctx context.Context, id string) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorw("failed to fetch user for deletion", "user_id", id, "error", err)
		return apperrors.NewAppError(apperrors.InternalError, "Failed to delete user")
	}
	if user == nil {
		s.logger.Warnw("user not found for deletion", "user_id", id)
		return apperrors.NewAppError(apperrors.NotFoundError, "User not found")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Errorw("failed to delete user", "user_id", id, "error", err)
		return apperrors.NewAppError(apperrors.InternalError, "Failed to delete user")
	}

	s.logger.Infow("user deleted", "user_id", id)
	return nil
}

// List fetches users with pagination, filtering, and sorting
func (s *userService) List(ctx context.Context, offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.User, error) {
	return s.repo.List(ctx, offset, limit, filters, sortBy, sortOrder)
}
