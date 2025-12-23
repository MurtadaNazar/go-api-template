package repo

import (
	"context"
	"errors"
	"go_platform_template/internal/domain/user/model"
	apperrors "go_platform_template/internal/shared/errors"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.User, error)
	GetByEmailOrUsername(ctx context.Context, identifier string) (*model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

// handleConstraintError converts database constraint errors to user-friendly messages
func handleConstraintError(err error) error {
	// Check for PostgreSQL unique constraint violations
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" { // unique_violation
			if strings.Contains(pgErr.Constraint, "username") {
				return apperrors.ErrUsernameAlreadyTaken
			}
			if strings.Contains(pgErr.Constraint, "email") {
				return apperrors.ErrEmailAlreadyRegistered
			}
		}
	}
	// Return original error if not a constraint violation
	return err
}

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return handleConstraintError(result.Error)
	}
	return nil
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Unscoped().First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Unscoped().First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Unscoped().Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Unscoped().Save(user).Error
}

// Delete fetches user by ID and deletes it
func (r *userRepo) Delete(ctx context.Context, id string) error {
	user, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return apperrors.ErrUserNotFound
	}
	return r.db.WithContext(ctx).Delete(user).Error
}

// List fetches users with optional filters, pagination, and sorting
func (r *userRepo) List(ctx context.Context, offset, limit int, filters map[string]interface{}, sortBy, sortOrder string) ([]*model.User, error) {
	var users []*model.User
	query := r.db.WithContext(ctx).Unscoped().Model(&model.User{})

	// Apply filters
	for key, val := range filters {
		query = query.Where(key+" = ?", val)
	}

	// Apply sorting
	if sortBy != "" {
		normalizedSortBy := strings.ToLower(sortBy)
		var column string
		switch normalizedSortBy {
		case "created_at":
			column = "created_at"
		case "username":
			column = "username"
		case "email":
			column = "email"
		default:
			// Fallback to a safe default if an unsupported sort field is requested
			column = "created_at"
		}

		// Explicitly map sortOrder to a validated, hard-coded value to prevent SQL injection
		// Only allow "asc" or "desc" (case-insensitive), default to "asc"
		normalizedSortOrder := strings.ToLower(sortOrder)
		direction := "asc"
		if normalizedSortOrder == "desc" {
			direction = "desc"
		}

		order := column + " " + direction
		query = query.Order(order)
	}

	// Apply pagination
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) GetByEmailOrUsername(ctx context.Context, identifier string) (*model.User, error) {
	if identifier == "" {
		return nil, nil
	}

	var user model.User
	err := r.db.WithContext(ctx).Unscoped().
		Where("LOWER(email) = LOWER(?) OR LOWER(username) = LOWER(?)", identifier, identifier).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}
