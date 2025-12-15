package database

import (
	"context"

	authModel "go_platform_template/internal/domain/auth/model"
	fileModel "go_platform_template/internal/domain/file/model"
	userModel "go_platform_template/internal/domain/user/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ctxKey is used for passing request-scoped values
type ctxKey string

const requestIDKey ctxKey = "RequestID"

// MigrateDB handles database migrations and ensures indexes are created
func MigrateDB(db *gorm.DB, log *zap.SugaredLogger) error {
	log.Info("Running database migrations...")

	// Enable uuid-ossp extension for UUID generation
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&userModel.User{},
		&authModel.RefreshToken{},
		&fileModel.File{},
	); err != nil {
		return err
	}

	log.Info("Database migration completed successfully.")

	// Create functional index (case-insensitive username search)
	if err := userModel.CreateFunctionalIndexes(db); err != nil {
		return err
	}

	log.Info("Functional indexes created successfully.")
	return nil
}

// WithRequestLogger returns a SugaredLogger enriched with request_id
// Extracts request ID from context using the proper custom key
func WithRequestLogger(ctx context.Context, logger *zap.SugaredLogger) *zap.SugaredLogger {
	// Try the custom key first
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok || requestID == "" {
		// Fallback to the string key
		requestID, ok = ctx.Value("RequestID").(string)
		if !ok || requestID == "" {
			requestID = "unknown"
		}
	}
	return logger.With("request_id", requestID)
}

// ExtractRequestID safely extracts request ID from context
// Returns the request ID or "unknown" if not found
func ExtractRequestID(ctx context.Context) string {
	// Try the custom key first
	if requestID, ok := ctx.Value(requestIDKey).(string); ok && requestID != "" {
		return requestID
	}
	// Fallback to the string key
	if requestID, ok := ctx.Value("RequestID").(string); ok && requestID != "" {
		return requestID
	}
	return "unknown"
}

// WithRequest returns a GORM DB instance with context containing request_id
// Extracts request ID from Gin context and propagates it to database operations
func WithRequest(c *gin.Context, db *gorm.DB) *gorm.DB {
	// Get request ID from Gin context (set by RequestIDMiddleware)
	requestID := c.GetString("RequestID")

	// Create new context with request ID using the proper custom key
	ctx := context.WithValue(c.Request.Context(), requestIDKey, requestID)

	// Return DB instance with the enriched context
	return db.WithContext(ctx)
}

// ApplyGlobalScopes applies global query filters
func ApplyGlobalScopes(db *gorm.DB) *gorm.DB {
	return db.Scopes(func(tx *gorm.DB) *gorm.DB {
		if tx.Statement.Schema != nil {
			if _, ok := tx.Statement.Schema.FieldsByDBName["deleted_at"]; ok {
				return tx.Where("deleted_at IS NULL")
			}
		}
		return tx
	})
}
