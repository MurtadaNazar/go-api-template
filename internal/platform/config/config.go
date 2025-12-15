package config

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type JWTConfig struct {
	SigningKey       string
	RefreshKey       string
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
}

type MinIOConfig struct {
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

type Config struct {
	ServerAddr        string
	APIVersion        string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	GinMode           string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime int // in seconds
	LogLevel          string
	JWT               JWTConfig
	MinIO             MinIOConfig
}

var (
	appConfig *Config
	once      sync.Once
)

// LoadConfig initializes the configuration once and parses DATABASE_URL
// It loads configuration from environment variables with the following precedence:
//   - System environment variables
//   - .env file values
//   - Default fallback values
//
// The function handles:
//   - Database connection pooling settings
//   - JWT token configuration with automatic secret generation
//   - MinIO object storage configuration
//   - Server and logging configuration
//
// Returns:
//   - *Config: Fully populated configuration object
func LoadConfig() *Config {
	once.Do(func() {
		// Load .env if exists
		_ = godotenv.Load()

		// -------------------------
		// Viper setup
		// -------------------------
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		viper.AutomaticEnv() // also read system env

		if err := viper.ReadInConfig(); err != nil {
			log.Println("No .env file found, relying on environment variables")
		}

		// -------------------------
		// Read DATABASE_URL
		// -------------------------
		dbURL := viper.GetString("DATABASE_URL")
		if dbURL == "" {
			panic("DATABASE_URL must be set")
		}

		parsedURL, err := url.Parse(dbURL)
		if err != nil {
			log.Fatalf("Invalid DATABASE_URL: %v", err)
		}

		user := parsedURL.User.Username()
		password, _ := parsedURL.User.Password()
		host := parsedURL.Hostname()
		port := parsedURL.Port()
		dbName := parsedURL.Path
		if len(dbName) > 0 && dbName[0] == '/' {
			dbName = dbName[1:] // remove leading slash
		}

		// -------------------------
		// Read SERVER_ADDR with fallback
		// -------------------------
		serverAddr := viper.GetString("SERVER_ADDR")
		if serverAddr == "" {
			serverAddr = ":8080"
		}

		// -------------------------
		// Read API_VERSION with fallback
		// -------------------------
		apiVersion := viper.GetString("API_VERSION")
		if apiVersion == "" {
			apiVersion = "v1"
		}
		log.Printf("[INFO] Using API version: %s", apiVersion)

		// -------------------------
		// Read GIN_MODE with fallback
		// -------------------------
		ginMode := viper.GetString("GIN_MODE")
		if ginMode == "" {
			ginMode = "release" // default to release mode
		}

		// -------------------------
		// Read Database Connection Pool Settings with fallbacks
		// -------------------------
		dbMaxOpenConns := viper.GetInt("DB_MAX_OPEN_CONNS")
		if dbMaxOpenConns == 0 {
			dbMaxOpenConns = 25 // default max open connections
		}

		dbMaxIdleConns := viper.GetInt("DB_MAX_IDLE_CONNS")
		if dbMaxIdleConns == 0 {
			dbMaxIdleConns = 5 // default max idle connections
		}

		dbConnMaxLifetime := viper.GetInt("DB_CONN_MAX_LIFETIME")
		if dbConnMaxLifetime == 0 {
			dbConnMaxLifetime = 300 // default 5 minutes (300 seconds)
		}

		// -------------------------
		// Read Log Level with fallback
		// -------------------------
		logLevel := viper.GetString("LOG_LEVEL")
		if logLevel == "" {
			logLevel = "info" // default log level
		}

		// -------------------------
		// Read JWT Configuration
		// -------------------------
		jwtKey := viper.GetString("JWT_SIGNING_KEY")
		if jwtKey == "" {
			jwtKey = generateRandomKey()
			log.Printf("[WARN] JWT signing key not found in environment. Generated new temporary key: %s", jwtKey)
		}

		// JWT refresh token key
		jwtRefreshKey := viper.GetString("JWT_REFRESH_KEY")
		if jwtRefreshKey == "" {
			jwtRefreshKey = generateRandomKey()
			log.Printf("[WARN] JWT refresh key not found in environment. Generated new temporary key: %s", jwtRefreshKey)
		}

		jwtAccessExpiresIn := parseDurationOrDefault(viper.GetString("JWT_ACCESS_EXPIRES_IN"), 15*time.Minute)
		jwtRefreshExpiresIn := parseDurationOrDefault(viper.GetString("JWT_REFRESH_EXPIRES_IN"), 7*24*time.Hour)

		// -------------------------
		// Read MinIO Configuration
		// -------------------------
		minioEndpoint := viper.GetString("MINIO_ENDPOINT")
		if minioEndpoint == "" {
			minioEndpoint = "localhost:9000" // default MinIO endpoint
			log.Printf("[INFO] MINIO_ENDPOINT not set, using default: %s", minioEndpoint)
		}

		minioAccessKey := viper.GetString("MINIO_ACCESS_KEY")
		if minioAccessKey == "" {
			minioAccessKey = "minioadmin" // default MinIO access key
			log.Printf("[INFO] MINIO_ACCESS_KEY not set, using default: %s", minioAccessKey)
		}

		minioSecretKey := viper.GetString("MINIO_SECRET_KEY")
		if minioSecretKey == "" {
			minioSecretKey = "minioadmin" // default MinIO secret key
			log.Printf("[INFO] MINIO_SECRET_KEY not set, using default: %s", minioSecretKey)
		}

		minioBucket := viper.GetString("MINIO_BUCKET")
		if minioBucket == "" {
			minioBucket = "go_platform_template" // default bucket name
			log.Printf("[INFO] MINIO_BUCKET not set, using default: %s", minioBucket)
		}

		minioUseSSL := viper.GetBool("MINIO_USE_SSL")
		// If not explicitly set, default to false for local development
		if !viper.IsSet("MINIO_USE_SSL") {
			minioUseSSL = false
			log.Printf("[INFO] MINIO_USE_SSL not set, using default: %t", minioUseSSL)
		}

		// -------------------------
		// Build Config
		// -------------------------
		// Set appConfig
		appConfig = &Config{
			ServerAddr:        serverAddr,
			APIVersion:        apiVersion,
			DBHost:            host,
			DBPort:            port,
			DBUser:            user,
			DBPassword:        password,
			DBName:            dbName,
			GinMode:           ginMode,
			DBMaxOpenConns:    dbMaxOpenConns,
			DBMaxIdleConns:    dbMaxIdleConns,
			DBConnMaxLifetime: dbConnMaxLifetime,
			LogLevel:          logLevel,
			JWT: JWTConfig{
				SigningKey:       jwtKey,
				RefreshKey:       jwtRefreshKey,
				AccessExpiresIn:  jwtAccessExpiresIn,
				RefreshExpiresIn: jwtRefreshExpiresIn,
			},
			MinIO: MinIOConfig{
				MinioEndpoint:  minioEndpoint,
				MinioAccessKey: minioAccessKey,
				MinioSecretKey: minioSecretKey,
				MinioBucket:    minioBucket,
				MinioUseSSL:    minioUseSSL,
			},
		}

		validateConfig(appConfig)
	})

	return appConfig
}

// GetConfig returns the already loaded config
// This function should be called after LoadConfig() has been initialized
//
// Returns:
//   - *Config: The singleton configuration instance
//
// Panics:
//   - If config has not been loaded via LoadConfig() first
func GetConfig() *Config {
	if appConfig == nil {
		panic("Config not loaded. Call LoadConfig() first")
	}
	return appConfig
}

// validateConfig ensures required fields are set and performs configuration validation
// It checks:
//   - Database configuration completeness
//   - Server address presence
//   - Database connection pool constraints
//   - JWT secret requirements
//
// Parameters:
//   - cfg: The configuration object to validate
//
// Panics:
//   - If any required configuration is invalid or missing
func validateConfig(cfg *Config) {
	// Database validation
	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		panic("Incomplete database configuration")
	}

	// Server validation
	if cfg.ServerAddr == "" {
		panic("Server address is required")
	}

	// Database connection pool validation
	if cfg.DBMaxOpenConns <= 0 {
		panic("DB_MAX_OPEN_CONNS must be greater than 0")
	}
	if cfg.DBMaxIdleConns <= 0 {
		panic("DB_MAX_IDLE_CONNS must be greater than 0")
	}
	if cfg.DBMaxIdleConns > cfg.DBMaxOpenConns {
		panic("DB_MAX_IDLE_CONNS cannot be greater than DB_MAX_OPEN_CONNS")
	}

	// JWT validation
	if cfg.JWT.SigningKey == "" {
		panic("JWT signing key must be set")
	}
	if cfg.JWT.RefreshKey == "" {
		panic("JWT refresh key must be set")
	}

	// MinIO validation (warning only, as it might be optional)
	if cfg.MinIO.MinioEndpoint == "" {
		log.Println("[WARN] MINIO_ENDPOINT is not set - MinIO operations will fail")
	}
	if cfg.MinIO.MinioAccessKey == "" || cfg.MinIO.MinioSecretKey == "" {
		log.Println("[WARN] MinIO access credentials are not set - MinIO operations will fail")
	}
	if cfg.MinIO.MinioBucket == "" {
		log.Println("[WARN] MINIO_BUCKET is not set - MinIO operations will fail")
	}
}

// parseDurationOrDefault parses a duration string or returns a default value
// Supported duration units: ns, us, ms, s, m, h, d (if using time.ParseDuration)
//
// Parameters:
//   - val: The duration string to parse (e.g., "15m", "1h", "24h")
//   - def: The default duration to return if parsing fails or string is empty
//
// Returns:
//   - time.Duration: The parsed duration or default value
func parseDurationOrDefault(val string, def time.Duration) time.Duration {
	if val == "" {
		return def
	}
	if d, err := time.ParseDuration(val); err == nil {
		return d
	}
	return def
}

// generateRandomSecret generates a cryptographically secure random secret
// using 32 bytes of random data encoded in base64 URL encoding
//
// Returns:
//   - string: Base64 URL encoded random secret (43 characters)
//
// Panics:
//   - If unable to read from cryptographically secure random source
func generateRandomKey() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Fatal("failed to generate random JWT key:", err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
