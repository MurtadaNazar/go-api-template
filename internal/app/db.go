package bootstrap

import (
	"database/sql"
	"fmt"
	"go_platform_template/internal/platform/config"
	"go_platform_template/internal/platform/database"

	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes the GORM DB connection with environment-aware logging
func InitDB(cfg *config.Config, log *zap.SugaredLogger) *gorm.DB {
	baseDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword)

	dbName := cfg.DBName

	// Ensure database exists before GORM connects
	if err := ensureDatabaseExists(baseDSN, dbName, log); err != nil {
		log.Errorf("failed to ensure database exists: %v", err)
		return nil
	}

	// Construct DSN with target DB
	dsn := fmt.Sprintf("%s dbname=%s", baseDSN, dbName)

	// Set GORM logger level based on environment
	var gormLogger logger.Interface
	if cfg.GinMode == "debug" || cfg.GinMode == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Warn)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Errorf("failed to connect database: %v", err)
		return nil
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLifetime) * time.Second)

	log.Info("Database connected successfully")

	// Run migrations and indexes
	if err := database.MigrateDB(db, log); err != nil {
		log.Errorf("Database migration failed: %v", err)
		return nil
	}

	// Seed admin user
	database.SeedAdminUser(db, log)
	log.Info("Database seeding completed")

	// Apply global scopes
	db = database.ApplyGlobalScopes(db)

	return db
}

// ensureDatabaseExists connects to Postgres without specifying dbname
// and creates the target database if it doesn't already exist.
func ensureDatabaseExists(baseDSN, dbName string, log *zap.SugaredLogger) error {
	conn, err := sql.Open("postgres", baseDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres (no db): %w", err)
	}
	defer conn.Close()

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	if err := conn.QueryRow(query, dbName).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		_, err := conn.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("failed to create database %s: %w", dbName, err)
		}
		log.Infof("Database %q created successfully", dbName)
	} else {
		log.Infof("Database %q already exists", dbName)
	}

	return nil
}
