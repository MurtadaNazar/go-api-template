package database

import (
	"context"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

// GormLogger wraps Zap logger to implement GORM logger.Interface
type GormLogger struct {
	Logger        *zap.SugaredLogger
	SlowThreshold time.Duration
	LogLevel      logger.LogLevel
	MaskSensitive bool
}

func NewGormLogger(zapLogger *zap.SugaredLogger, maskSensitive bool) *GormLogger {
	return &GormLogger{
		Logger:        zapLogger,
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Info, // Use the GORM LogLevel constant, not the method
		MaskSensitive: maskSensitive,
	}
}

// LogMode implements logger.Interface
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info logs info messages
func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Info {
		requestID := ExtractRequestID(ctx)
		l.Logger.With("request_id", requestID).Infof(msg, args...)
	}
}

// Warn logs warnings
func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Warn {
		requestID := ExtractRequestID(ctx)
		l.Logger.With("request_id", requestID).Warnf(msg, args...)
	}
}

// Error logs errors
func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Error {
		requestID := ExtractRequestID(ctx)
		l.Logger.With("request_id", requestID).Errorf(msg, args...)
	}
}

// maskQuery masks sensitive data in SQL queries
func maskQuery(query string) string {
	// Mask all string literals first: 'anything' -> '***'
	masked := regexp.MustCompile(`'[^']*'`).ReplaceAllString(query, "'***'")

	// List of sensitive keywords to look for
	sensitiveKeywords := []string{
		"password",
		"token",
		"secret",
		"api_key",
		"access_key",
	}

	// Mask values assigned to these keywords in SQL (basic approach)
	for _, key := range sensitiveKeywords {
		// regex: key\s*=\s*'...' or key='...'
		re := regexp.MustCompile("(?i)" + key + `\s*=\s*'[^']*'`)
		masked = re.ReplaceAllStringFunc(masked, func(match string) string {
			parts := strings.Split(match, "=")
			if len(parts) == 2 {
				return parts[0] + "= '***'"
			}
			return match
		})
	}

	return masked
}

// Trace logs SQL queries with optional masking
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// Mask sensitive data if enabled
	if l.MaskSensitive {
		sql = maskQuery(sql)
	}

	// Safe request_id extraction
	requestID := ExtractRequestID(ctx)

	fields := map[string]interface{}{
		"request_id": requestID,
		"elapsed":    elapsed,
		"rows":       rows,
		"sql":        sql,
	}

	switch {
	case err != nil && l.LogLevel >= logger.Error:
		l.Logger.With(fields).Error("SQL error")
	case elapsed > l.SlowThreshold && l.LogLevel >= logger.Warn:
		l.Logger.With(fields).Warn("Slow SQL query")
	case l.LogLevel >= logger.Info:
		l.Logger.With(fields).Info("SQL executed")
	}
}
