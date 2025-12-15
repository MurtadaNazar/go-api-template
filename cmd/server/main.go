package main

import (
	_ "go_platform_template/docs" // Important: import the generated docs
	bootstrap "go_platform_template/internal/app"
	"go_platform_template/internal/platform/config"
	"go_platform_template/internal/platform/logger"

	"github.com/gin-gonic/gin"
)

// @title           Go Platform Template API
// @version         1.0
// @description     Go Platform Template - Base platform with Auth, User, File, and RBAC features
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load config
	cfg := config.LoadConfig()

	// Init logger
	logr := logger.InitLogger()
	defer func() { _ = logr.Logger.Sync() }()
	logr.Sugar.Infof("Starting go-platform-template server on %s", cfg.ServerAddr)

	// Init DB
	db := bootstrap.InitDB(cfg, logr.Sugar)

	// Init Gin
	r := gin.New()
	bootstrap.SetupMiddleware(r, logr.Sugar)

	// Register domain routes
	bootstrap.RegisterRoutes(r, db, cfg, logr.Sugar)

	// Setup Swagger
	bootstrap.SetupSwagger(r, cfg, logr.Sugar)

	// Health check
	r.GET("/health", bootstrap.HealthCheckHandler(db, logr.Sugar))

	// Start server
	bootstrap.StartServer(r, cfg.ServerAddr, db, logr.Sugar)
}
