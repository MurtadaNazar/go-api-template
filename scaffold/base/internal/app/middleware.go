package bootstrap

import (
	"go_platform_template/internal/platform/http/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupMiddleware adds all your prebuilt middlewares to the Gin engine
func SetupMiddleware(r *gin.Engine, log *zap.SugaredLogger) {
	r.Use(
		middleware.RequestIDMiddleware(),
		middleware.LoggerMiddleware(log),
		middleware.RecoveryMiddleware(log),
		middleware.ErrorHandlerMiddleware(log), // Global error handler
		middleware.CORSMiddleware(),
		middleware.RateLimitMiddleware(),
		// middleware.JWTAuthMiddleware(nil), // for global JWT if needed, or per-route
	)
}
