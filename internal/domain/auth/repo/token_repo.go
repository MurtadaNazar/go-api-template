package repo

import (
	"context"
	"errors"
	"go_platform_template/internal/domain/auth/model"
	apperrors "go_platform_template/internal/shared/errors"
	"time"

	"gorm.io/gorm"
)

type TokenRepo interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	RevokeToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
	DeleteExpiredTokens(ctx context.Context) error
}

type tokenRepo struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) TokenRepo {
	return &tokenRepo{db: db}
}

func (r *tokenRepo) Create(ctx context.Context, token *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *tokenRepo) FindByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ? AND is_revoked = ? AND expires_at > ?",
		token, false, time.Now()).First(&refreshToken).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.ErrTokenNotFoundExpired
	}

	return &refreshToken, err
}

func (r *tokenRepo) RevokeToken(ctx context.Context, token string) error {
	result := r.db.WithContext(ctx).Model(&model.RefreshToken{}).
		Where("token = ?", token).
		Update("is_revoked", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrTokenNotFound
	}

	return nil
}

func (r *tokenRepo) RevokeAllUserTokens(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&model.RefreshToken{}).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Update("is_revoked", true).Error
}

func (r *tokenRepo) DeleteExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).
		Delete(&model.RefreshToken{}).Error
}
