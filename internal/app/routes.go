package bootstrap

import (
	"time"

	"go_platform_template/internal/platform/config"
	"go_platform_template/internal/platform/http/middleware"

	authApi "go_platform_template/internal/domain/auth/api"
	authRepo "go_platform_template/internal/domain/auth/repo"
	authService "go_platform_template/internal/domain/auth/service"

	userApi "go_platform_template/internal/domain/user/api"
	userRepo "go_platform_template/internal/domain/user/repo"
	userService "go_platform_template/internal/domain/user/service"

	fileApi "go_platform_template/internal/domain/file/api"
	fileRepo "go_platform_template/internal/domain/file/repo"
	fileService "go_platform_template/internal/domain/file/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

/*
===========================
DEVELOPER NOTES: AUTH & ROLES
===========================

1. Roles (UserType):
    - "admin"      : System admin, full access.
    - "user"       : Regular user with standard access.

2. JWT Authentication Flow:
   - User logs in via `/login` and receives:
     a) Access Token  (short-lived, used in Authorization header)
     b) Refresh Token (long-lived, stored in DB, used to get new access token)
   - Access Token claims include:
       - user_id
       - role (matches UserType)
       - exp (expiration timestamp)
   - Protected endpoints require a valid access token in the header:
       Authorization: Bearer <access_token>

3. Protecting Endpoints:
   - Use `middleware.JWTAuth(jwtManager)` in your route group.
   - Handlers can get user info via context:
       userID := c.GetString("userID")
       role   := c.GetString("role")
   - Optional: enforce role-based access in handlers:
       if role != "admin" {
           c.JSON(403, gin.H{"error": "forbidden"})
           return
       }

4. Token Rotation & Logout:
   - Refresh tokens are stored in DB (tokenStore) and can be revoked.
   - Logout revokes the refresh token.
   - Refresh endpoint rotates refresh tokens for better security.

5. Example Usage:
    - Admin-only route:
        protected.GET("/users", requireRole("admin"), uHandler.ListUsers)
    - User route:
        protected.GET("/profile", requireRole("user"), uHandler.GetUser)

6. Notes:
   - Always pass JWTManager to middleware and AuthService to handlers.
   - Use Refresh token only for `/refresh` endpoint.
*/

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config, log *zap.SugaredLogger) {
	// -----------------------
	// JWT & Auth setup
	// -----------------------
	jwtManager := authService.NewJWTManager(
		cfg.JWT.SigningKey,
		cfg.JWT.RefreshKey,
		cfg.JWT.AccessExpiresIn,
		cfg.JWT.RefreshExpiresIn,
	)

	uRepo := userRepo.NewUserRepo(db)
	uService := userService.NewUserService(uRepo, log)
	uHandler := userApi.NewUserHandler(uService, log)

	tRepo := authRepo.NewTokenRepo(db)
	tStore := authService.NewTokenStore(tRepo, log)
	aService := authService.NewAuthService(uRepo, jwtManager, tStore, log)
	aHandler := authApi.NewAuthHandler(aService, log)

	// Start background job to clean up expired tokens every 24 hours
	go authService.StartTokenCleanupJob(tStore, 24*time.Hour)

	fRepo := fileRepo.NewFileRepo(db)
	var fileHandler *fileApi.FileHandler
	fSvc, err := fileService.NewFileService(fRepo, cfg, log)
	if err != nil {
		log.Warnf("FileService initialization failed (MinIO unavailable): %v", err)
		log.Warn("File upload/download endpoints will be unavailable")
		// Continue without file service - file endpoints won't be registered
	} else {
		fileHandler = fileApi.NewFileHandler(fSvc, log)
	}

	// -----------------------
	// API Versioning: v1
	// -----------------------
	v1 := r.Group("/api/v1")
	{
		// -----------------------
		// Auth routes
		// -----------------------
		auth := v1.Group("/")
		{
			auth.POST("/login", aHandler.Login)
			auth.POST("/refresh", aHandler.Refresh)
			auth.POST("/logout", aHandler.Logout)
		}

		// -----------------------
		// User routes
		// -----------------------
		users := v1.Group("/users")
		{
			users.POST("/", uHandler.Register)
			users.GET("/", middleware.JWTAuth(jwtManager), uHandler.ListUsers)
			users.GET("/:id", middleware.JWTAuth(jwtManager), uHandler.GetUser)
			users.PUT("/:id", middleware.JWTAuth(jwtManager), uHandler.Update)
			users.DELETE("/:id", middleware.JWTAuth(jwtManager), uHandler.Delete)
		}

		// -----------------------
		// Protected routes
		// -----------------------
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth(jwtManager))
		{
			protected.GET("/me", aHandler.Me)
		}

		// -----------------------
		// File routes (only if MinIO available)
		// -----------------------
		if fSvc != nil {
			files := v1.Group("/files")
			files.Use(middleware.JWTAuth(jwtManager))
			{
				files.POST("/upload", fileHandler.Upload)
				files.GET("/:filename", fileHandler.GetFile)
				files.DELETE("/:filename", fileHandler.DeleteFile)
				files.GET("/", fileHandler.GetUserFiles)
			}
		}

	}

	log.Info("Routes registered successfully under /api/v1")
}
