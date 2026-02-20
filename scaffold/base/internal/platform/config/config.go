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
	DBConnMaxLifetime int
	LogLevel          string
	JWT               JWTConfig
	MinIO             MinIOConfig
}

var (
	appConfig *Config
	once      sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		_ = godotenv.Load()

		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			log.Println("No .env file found, relying on environment variables")
		}

		var dbHost, dbPort, dbUser, dbPassword, dbName string

		dbURL := viper.GetString("DATABASE_URL")
		if dbURL != "" {
			parsedURL, err := url.Parse(dbURL)
			if err != nil {
				log.Fatalf("Invalid DATABASE_URL: %v", err)
			}
			dbUser = parsedURL.User.Username()
			dbPassword, _ = parsedURL.User.Password()
			dbHost = parsedURL.Hostname()
			dbPort = parsedURL.Port()
			dbName = parsedURL.Path
			if len(dbName) > 0 && dbName[0] == '/' {
				dbName = dbName[1:]
			}
		} else {
			dbHost = getEnvWithDefault("DB_HOST", "localhost")
			dbPort = getEnvWithDefault("DB_PORT", "5432")
			dbUser = getEnvWithDefault("DB_USER", "postgres")
			dbPassword = getEnvWithDefault("DB_PASSWORD", "postgres")
			dbName = getEnvWithDefault("DB_NAME", "test")
		}

		serverAddr := getEnvWithDefault("SERVER_ADDR", ":8080")
		apiVersion := getEnvWithDefault("API_VERSION", "v1")
		ginMode := getEnvWithDefault("GIN_MODE", "release")
		logLevel := getEnvWithDefault("LOG_LEVEL", "info")

		dbMaxOpenConns := viper.GetInt("DB_MAX_OPEN_CONNS")
		if dbMaxOpenConns == 0 {
			dbMaxOpenConns = 25
		}
		dbMaxIdleConns := viper.GetInt("DB_MAX_IDLE_CONNS")
		if dbMaxIdleConns == 0 {
			dbMaxIdleConns = 5
		}
		dbConnMaxLifetime := viper.GetInt("DB_CONN_MAX_LIFETIME")
		if dbConnMaxLifetime == 0 {
			dbConnMaxLifetime = 300
		}

		jwtSigningKey := viper.GetString("JWT_SIGNING_KEY")
		if jwtSigningKey == "" {
			jwtSigningKey = generateRandomKey()
			log.Println("[WARN] JWT_SIGNING_KEY not set, generated temporary key")
		}
		jwtRefreshKey := viper.GetString("JWT_REFRESH_KEY")
		if jwtRefreshKey == "" {
			jwtRefreshKey = generateRandomKey()
			log.Println("[WARN] JWT_REFRESH_KEY not set, generated temporary key")
		}
		jwtAccessExpiry := parseDurationOrDefault(viper.GetString("JWT_ACCESS_EXPIRY"), 15*time.Minute)
		jwtRefreshExpiry := parseDurationOrDefault(viper.GetString("JWT_REFRESH_EXPIRY"), 7*24*time.Hour)

		minioEndpoint := getEnvWithDefault("MINIO_ENDPOINT", "localhost:9000")
		minioAccessKey := getEnvWithDefault("MINIO_ACCESS_KEY", "minioadmin")
		minioSecretKey := getEnvWithDefault("MINIO_SECRET_KEY", "minioadmin")
		minioBucket := getEnvWithDefault("MINIO_BUCKET", "uploads")
		minioUseSSL := viper.GetBool("MINIO_SECURE")

		appConfig = &Config{
			ServerAddr:        serverAddr,
			APIVersion:        apiVersion,
			DBHost:            dbHost,
			DBPort:            dbPort,
			DBUser:            dbUser,
			DBPassword:        dbPassword,
			DBName:            dbName,
			GinMode:           ginMode,
			DBMaxOpenConns:    dbMaxOpenConns,
			DBMaxIdleConns:    dbMaxIdleConns,
			DBConnMaxLifetime: dbConnMaxLifetime,
			LogLevel:          logLevel,
			JWT: JWTConfig{
				SigningKey:       jwtSigningKey,
				RefreshKey:       jwtRefreshKey,
				AccessExpiresIn:  jwtAccessExpiry,
				RefreshExpiresIn: jwtRefreshExpiry,
			},
			MinIO: MinIOConfig{
				MinioEndpoint:  minioEndpoint,
				MinioAccessKey: minioAccessKey,
				MinioSecretKey: minioSecretKey,
				MinioBucket:    minioBucket,
				MinioUseSSL:    minioUseSSL,
			},
		}
	})

	return appConfig
}

func GetConfig() *Config {
	if appConfig == nil {
		panic("Config not loaded. Call LoadConfig() first")
	}
	return appConfig
}

func getEnvWithDefault(key, defaultVal string) string {
	val := viper.GetString(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func parseDurationOrDefault(val string, def time.Duration) time.Duration {
	if val == "" {
		return def
	}
	if d, err := time.ParseDuration(val); err == nil {
		return d
	}
	return def
}

func generateRandomKey() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Fatal("failed to generate random JWT key:", err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
