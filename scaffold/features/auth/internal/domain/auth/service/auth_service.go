package service

import (
	"context"
	"go_platform_template/internal/domain/user/repo"
	apperrors "go_platform_template/internal/shared/errors"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   repo.UserRepo
	jwt        *JWTManager
	tokenStore *TokenStore
	logger     *zap.SugaredLogger
}

func NewAuthService(userRepo repo.UserRepo, jwt *JWTManager, store *TokenStore, logger *zap.SugaredLogger) *AuthService {
	return &AuthService{userRepo: userRepo, jwt: jwt, tokenStore: store, logger: logger}
}

func (s *AuthService) Login(ctx context.Context, emailOrUsername, password string) (string, string, error) {
	// Try to find user by email OR username
	user, err := s.userRepo.GetByEmailOrUsername(ctx, emailOrUsername)
	if err != nil {
		s.logger.Errorw("failed to fetch user", "email_or_username", emailOrUsername, "error", err)
		return "", "", apperrors.NewAppError(apperrors.UnauthorizedError, "Invalid credentials")
	}
	if user == nil {
		s.logger.Warnw("user not found", "email_or_username", emailOrUsername)
		return "", "", apperrors.NewAppError(apperrors.UnauthorizedError, "Invalid credentials")
	}

	// Check if user is active
	if !user.IsActive() {
		s.logger.Warnw("inactive user login attempt", "user_id", user.ID)
		return "", "", apperrors.NewAppError(apperrors.ForbiddenError, "Account is inactive")
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.logger.Warnw("invalid password", "user_id", user.ID)
		return "", "", apperrors.NewAppError(apperrors.UnauthorizedError, "Invalid credentials")
	}

	// Generate tokens
	access, refresh, err := s.jwt.GenerateTokens(user.ID, string(user.UserType))
	if err != nil {
		s.logger.Errorw("failed to generate tokens", "user_id", user.ID, "error", err)
		return "", "", apperrors.NewAppError(apperrors.InternalError, "Failed to generate authentication tokens")
	}

	// Save refresh token
	if err := s.tokenStore.Save(ctx, refresh, user.ID, string(user.UserType), time.Now().Add(s.jwt.refreshExpires)); err != nil {
		s.logger.Errorw("failed to save refresh token", "user_id", user.ID, "error", err)
		return "", "", apperrors.NewAppError(apperrors.InternalError, "Failed to save authentication token")
	}

	s.logger.Infow("user logged in", "user_id", user.ID)
	return access, refresh, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	data, err := s.tokenStore.Validate(ctx, refreshToken, true)
	if err != nil {
		s.logger.Errorw("failed to validate refresh token", "error", err)
		return "", "", apperrors.NewAppError(apperrors.UnauthorizedError, "Invalid or expired refresh token")
	}

	access, newRefresh, err := s.jwt.GenerateTokens(data.UserID, data.Role)
	if err != nil {
		s.logger.Errorw("failed to generate new tokens", "user_id", data.UserID, "error", err)
		return "", "", apperrors.NewAppError(apperrors.InternalError, "Failed to generate new tokens")
	}

	if err := s.tokenStore.Save(ctx, newRefresh, data.UserID, data.Role, time.Now().Add(s.jwt.refreshExpires)); err != nil {
		s.logger.Errorw("failed to save new refresh token", "user_id", data.UserID, "error", err)
		return "", "", apperrors.NewAppError(apperrors.InternalError, "Failed to save new token")
	}

	s.logger.Infow("tokens refreshed", "user_id", data.UserID)
	return access, newRefresh, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if err := s.tokenStore.Delete(ctx, refreshToken); err != nil {
		s.logger.Errorw("failed to logout", "error", err)
		return apperrors.NewAppError(apperrors.InternalError, "Failed to logout")
	}
	return nil
}
