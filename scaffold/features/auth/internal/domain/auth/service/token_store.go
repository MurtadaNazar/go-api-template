package service

import (
	"context"
	"go_platform_template/internal/domain/auth/model"
	"go_platform_template/internal/domain/auth/repo"
	apperrors "go_platform_template/internal/shared/errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TokenStore struct {
	repo   repo.TokenRepo
	logger *zap.SugaredLogger
}

type RefreshTokenData struct {
	UserID    uuid.UUID
	Role      string
	ExpiresAt time.Time
}

func NewTokenStore(r repo.TokenRepo, logger *zap.SugaredLogger) *TokenStore {
	return &TokenStore{repo: r, logger: logger}
}

func (s *TokenStore) Save(ctx context.Context, token string, userID uuid.UUID, role string, expiresAt time.Time) error {
	rt := &model.RefreshToken{
		Token:     token,
		UserID:    userID,
		Role:      role,
		ExpiresAt: expiresAt,
		IsRevoked: false,
	}
	if err := s.repo.Create(ctx, rt); err != nil {
		s.logger.Errorf("Save refresh token failed: %v", err)
		return err
	}
	return nil
}

// Validate refresh token and optionally rotate (revoke old)
func (s *TokenStore) Validate(ctx context.Context, token string, rotate bool) (*RefreshTokenData, error) {
	rt, err := s.repo.FindByToken(ctx, token)
	if err != nil {
		return nil, apperrors.ErrInvalidRefreshToken
	}

	if rotate {
		if err := s.repo.RevokeToken(ctx, token); err != nil {
			return nil, err
		}
	}

	return &RefreshTokenData{
		UserID:    rt.UserID,
		Role:      rt.Role,
		ExpiresAt: rt.ExpiresAt,
	}, nil
}

func (s *TokenStore) Delete(ctx context.Context, token string) error {
	return s.repo.RevokeToken(ctx, token)
}

func (s *TokenStore) CleanupExpiredTokens(ctx context.Context) error {
	return s.repo.DeleteExpiredTokens(ctx)
}
