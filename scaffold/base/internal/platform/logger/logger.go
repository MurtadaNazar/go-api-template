package logger

import (
	"go_platform_template/internal/platform/config"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger wraps zap.Logger and zap.SugaredLogger for structured logging
type Logger struct {
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
}

// InitLogger initializes and returns a Logger instance
// It supports:
//   - JSON structured logs
//   - Log level from config (debug, info, warn, error)
//   - Console + file logging
//   - File rotation via lumberjack
func InitLogger() *Logger {
	cfg := config.GetConfig() // Load already initialized config

	// Determine log level
	var zapLevel zapcore.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Lumberjack log rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logs/app.log", // log file path
		MaxSize:    50,             // megabytes
		MaxBackups: 7,              // number of old log files to keep
		MaxAge:     30,             // days
		Compress:   true,           // compress old log files
	}

	// JSON encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Write logs to both console and file
	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(lumberjackLogger),
	)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		zapLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		Logger: logger,
		Sugar:  logger.Sugar(),
	}
}
