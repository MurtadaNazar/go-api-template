package middleware

import (
	apperrors "go_platform_template/internal/shared/errors"
	"go_platform_template/internal/shared/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandlerMiddleware handles errors consistently across the application
// It intercepts errors, logs them, and returns standardized error responses
func ErrorHandlerMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Handle errors if any occurred
		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last()
			requestID, _ := c.Get("RequestID")

			// Check if it's an AppError
			if appErr, ok := apperrors.IsAppError(lastErr.Err); ok {
				// Log the error with context
				logger.Errorw("request error",
					"request_id", requestID,
					"status", appErr.HTTPStatus,
					"error_type", appErr.Type,
					"message", appErr.Message,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)

				// Return standardized error response
				c.JSON(appErr.HTTPStatus, response.NewErrorResponse(
					appErr.Message,
					string(appErr.Type),
					appErr.Details,
					requestID.(string),
				))
				return
			}

			// Handle generic errors
			logger.Errorw("unexpected error",
				"request_id", requestID,
				"error", lastErr.Err.Error(),
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
			)

			c.JSON(http.StatusInternalServerError, response.NewErrorResponse(
				"An unexpected error occurred",
				"INTERNAL",
				lastErr.Err.Error(),
				requestID.(string),
			))
		}
	}
}
