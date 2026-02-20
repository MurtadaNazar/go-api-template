package middleware

import (
	"go_platform_template/internal/domain/auth/service"
	apperrors "go_platform_template/internal/shared/errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth(jwtManager *service.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			_ = c.Error(apperrors.NewAppError(apperrors.UnauthorizedError, "missing authorization header"))
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if token == "" {
			_ = c.Error(apperrors.NewAppError(apperrors.UnauthorizedError, "invalid authorization header"))
			return
		}

		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			_ = c.Error(apperrors.NewAppError(apperrors.UnauthorizedError, "invalid or expired token"))
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
