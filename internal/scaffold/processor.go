package scaffold

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
)

type ProcessCompleteMsg struct {
	Message string
	Err     error
}

// scaffoldFS will be set by init in main package
var scaffoldFS fs.FS

// SetScaffoldFS allows the main package to inject the embedded scaffold FS
func SetScaffoldFS(fs fs.FS) {
	scaffoldFS = fs
}

func (m *Model) processScaffold() tea.Cmd {
	return func() tea.Msg {
		selectedFeatures := make(map[string]bool)
		for _, feat := range m.features {
			selectedFeatures[feat.Name] = feat.Selected
		}

		if err := createProject(m.projectName, m.moduleName, m.projectPath, selectedFeatures); err != nil {
			return ProcessCompleteMsg{Err: err}
		}
		return ProcessCompleteMsg{
			Message: fmt.Sprintf("Project '%s' created successfully", m.projectName),
		}
	}
}

func createProject(projectName, moduleName, projectPath string, selectedFeatures map[string]bool) error {
	// Resolve project path
	var basePath string
	if projectPath == "." {
		basePath, _ = os.Getwd()
	} else {
		var err error
		basePath, err = filepath.Abs(projectPath)
		if err != nil {
			return fmt.Errorf("invalid project path: %w", err)
		}
	}

	projectDir := filepath.Join(basePath, projectName)

	// Check if directory exists
	if _, err := os.Stat(projectDir); err == nil {
		return fmt.Errorf("directory '%s' already exists", projectName)
	}

	// Create project directory
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Copy base files first (from embedded FS)
	if err := copyBaseScaffoldFromEmbed(projectDir); err != nil {
		os.RemoveAll(projectDir)
		return fmt.Errorf("failed to copy base files: %w", err)
	}

	// Copy selected features (from embedded FS)
	if err := copySelectedFeaturesFromEmbed(projectDir, selectedFeatures); err != nil {
		os.RemoveAll(projectDir)
		return fmt.Errorf("failed to copy features: %w", err)
	}

	// Generate main.go from template
	if err := generateMainGo(projectDir, moduleName, selectedFeatures); err != nil {
		os.RemoveAll(projectDir)
		return fmt.Errorf("failed to generate main.go: %w", err)
	}

	// Generate routes.go from template
	if err := generateRoutesGo(projectDir, moduleName, selectedFeatures); err != nil {
		os.RemoveAll(projectDir)
		return fmt.Errorf("failed to generate routes.go: %w", err)
	}

	// Replace placeholders
	if err := replaceModuleNames(projectDir, projectName, moduleName); err != nil {
		os.RemoveAll(projectDir)
		return fmt.Errorf("failed to update module names: %w", err)
	}

	// Initialize git
	if err := initializeGit(projectDir); err != nil {
		os.RemoveAll(projectDir)
		return fmt.Errorf("failed to initialize git: %w", err)
	}

	return nil
}

func copyBaseScaffoldFromEmbed(projectDir string) error {
	// Copy base files from embedded FS (scaffold/base/)
	baseDir := "scaffold/base"

	return copyDirFromEmbed(baseDir, projectDir)
}

func copySelectedFeaturesFromEmbed(projectDir string, selectedFeatures map[string]bool) error {
	scaffoldDir := "scaffold/features"

	// Feature ID to name mapping
	featureMap := map[string]string{
		"Authentication (JWT)": "auth",
		"User Management":      "user-management",
		"Database":             "database",
		"File Storage":         "file-storage",
		"API Docs":             "api-docs",
		"Docker":               "docker",
	}

	for featureName, isSelected := range selectedFeatures {
		if !isSelected {
			continue
		}

		featureID, ok := featureMap[featureName]
		if !ok {
			continue
		}

		featureDir := filepath.Join(scaffoldDir, featureID)

		// Read feature definition from embedded FS
		featureFile := filepath.Join(featureDir, "feature.json")
		content, err := fs.ReadFile(scaffoldFS, featureFile)
		if err != nil {
			// Feature not set up, skip
			continue
		}

		var feature struct {
			Directories       []string `json:"directories"`
			Files             []string `json:"files"`
			DirectoriesToCopy []string `json:"directories_to_copy"`
		}

		if err := parseJSON(content, &feature); err != nil {
			continue
		}

		// Copy directories for this feature
		for _, dir := range feature.DirectoriesToCopy {
			srcPath := filepath.Join("scaffold", dir)
			dstPath := filepath.Join(projectDir, dir)

			if _, err := fs.Stat(scaffoldFS, srcPath); err == nil {
				if err := copyDirFromEmbed(srcPath, dstPath); err != nil {
					continue
				}
			}
		}

		// Copy files for this feature
		for _, file := range feature.Files {
			srcPath := filepath.Join("scaffold", file)
			dstPath := filepath.Join(projectDir, file)

			if _, err := fs.Stat(scaffoldFS, srcPath); err == nil {
				if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
					continue
				}
				content, err := fs.ReadFile(scaffoldFS, srcPath)
				if err != nil {
					continue
				}
				if err := os.WriteFile(dstPath, content, 0600); err != nil {
					continue
				}
			}
		}
	}

	return nil
}

func copyDirFromEmbed(srcPath, dstPath string) error {
	return fs.WalkDir(scaffoldFS, srcPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path from source
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}

		// Skip the source directory itself
		if relPath == "." {
			return nil
		}

		targetPath := filepath.Join(dstPath, relPath)

		if entry.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Create parent directory
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// Read and write file
		content, err := fs.ReadFile(scaffoldFS, path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, content, 0600)
	})
}

func copyBaseScaffold(templateDir, projectDir string) error {
	// Copy base files from scaffold/base/
	baseDir := filepath.Join(templateDir, "scaffold", "base")

	// Check if base directory exists
	if _, err := os.Stat(baseDir); err != nil {
		return fmt.Errorf("scaffold/base directory not found: %w", err)
	}

	// Copy entire base directory structure to project
	if err := copyDir(baseDir, projectDir); err != nil {
		return fmt.Errorf("failed to copy scaffold base: %w", err)
	}

	return nil
}

func copySelectedFeatures(templateDir, projectDir string, selectedFeatures map[string]bool) error {
	scaffoldDir := filepath.Join(templateDir, "scaffold", "features")

	// Feature ID to name mapping
	featureMap := map[string]string{
		"Authentication (JWT)": "auth",
		"User Management":      "user-management",
		"Database":             "database",
		"File Storage":         "file-storage",
		"API Docs":             "api-docs",
		"Docker":               "docker",
	}

	for featureName, isSelected := range selectedFeatures {
		if !isSelected {
			continue
		}

		featureID, ok := featureMap[featureName]
		if !ok {
			continue
		}

		featureDir := filepath.Join(scaffoldDir, featureID)
		featureFile := filepath.Join(featureDir, "feature.json")

		// Check if feature definition exists
		if _, err := os.Stat(featureFile); err != nil {
			// Feature not yet set up, copy from template
			continue
		}

		// Read feature definition
		content, err := os.ReadFile(featureFile)
		if err != nil {
			continue
		}

		var feature struct {
			Directories       []string `json:"directories"`
			Files             []string `json:"files"`
			DirectoriesToCopy []string `json:"directories_to_copy"`
		}

		if err := parseJSON(content, &feature); err != nil {
			continue
		}

		// Copy directories for this feature
		for _, dir := range feature.DirectoriesToCopy {
			srcPath := filepath.Join(templateDir, dir)
			dstPath := filepath.Join(projectDir, dir)

			if _, err := os.Stat(srcPath); err == nil {
				if err := copyDir(srcPath, dstPath); err != nil {
					// Log but continue
					continue
				}
			}
		}

		// Copy files for this feature
		for _, file := range feature.Files {
			srcPath := filepath.Join(templateDir, file)
			dstPath := filepath.Join(projectDir, file)

			if _, err := os.Stat(srcPath); err == nil {
				if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
					continue
				}
				if err := copyFile(srcPath, dstPath); err != nil {
					// Log but continue
					continue
				}
			}
		}
	}

	return nil
}

func parseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func generateMainGo(projectDir, moduleName string, selectedFeatures map[string]bool) error {
	mainGoTemplate := `package main

import (
{{if .HasDocs}}	_ "{{.Module}}/docs" // Important: import the generated docs
{{end}}	bootstrap "{{.Module}}/internal/app"
	"{{.Module}}/internal/platform/config"
	"{{.Module}}/internal/platform/logger"

	"github.com/gin-gonic/gin"
)

// @title           Go Platform Template API
// @version         1.0
// @description     Go Platform Template - Production-ready Go API platform
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
	logr.Sugar.Infof("Starting server on %s", cfg.ServerAddr)

{{if .HasDatabase}}	// Init DB
	db := bootstrap.InitDB(cfg, logr.Sugar)
{{end}}
	// Init Gin
	r := gin.New()
	bootstrap.SetupMiddleware(r, logr.Sugar)

	// Register domain routes
{{if .HasDatabase}}	bootstrap.RegisterRoutes(r, db, cfg, logr.Sugar)
{{else}}	// No database features configured
{{end}}
{{if .HasDocs}}	// Setup Swagger
	bootstrap.SetupSwagger(r, cfg, logr.Sugar)
{{end}}
	// Health check
{{if .HasDatabase}}	r.GET("/health", bootstrap.HealthCheckHandler(db, logr.Sugar))
{{else}}	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
{{end}}
	// Start server
{{if .HasDatabase}}	bootstrap.StartServer(r, cfg.ServerAddr, db, logr.Sugar)
{{else}}	bootstrap.StartServer(r, cfg.ServerAddr, nil, logr.Sugar)
{{end}}
}
`

	data := struct {
		Module      string
		HasAuth     bool
		HasUser     bool
		HasDatabase bool
		HasFile     bool
		HasDocs     bool
		HasDocker   bool
	}{
		Module:      moduleName,
		HasAuth:     selectedFeatures["Authentication (JWT)"],
		HasUser:     selectedFeatures["User Management"],
		HasDatabase: selectedFeatures["Database"],
		HasFile:     selectedFeatures["File Storage"],
		HasDocs:     selectedFeatures["API Docs"],
		HasDocker:   selectedFeatures["Docker"],
	}

	tmpl, err := template.New("main.go").Parse(mainGoTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse main.go template: %w", err)
	}

	mainGoPath := filepath.Join(projectDir, "cmd", "server", "main.go")

	// Create directory structure first
	if err := os.MkdirAll(filepath.Dir(mainGoPath), 0755); err != nil {
		return fmt.Errorf("failed to create cmd/server directory: %w", err)
	}

	f, err := os.Create(mainGoPath)
	if err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute main.go template: %w", err)
	}

	return nil
}

func generateRoutesGo(projectDir, moduleName string, selectedFeatures map[string]bool) error {
	routesGoTemplate := `package bootstrap

import (
	"time"

	"{{.Module}}/internal/platform/config"
	"{{.Module}}/internal/platform/http/middleware"
{{if .HasAuth}}
	authApi "{{.Module}}/internal/domain/auth/api"
	authRepo "{{.Module}}/internal/domain/auth/repo"
	authService "{{.Module}}/internal/domain/auth/service"
{{end}}
{{if .HasUser}}
	userApi "{{.Module}}/internal/domain/user/api"
	userRepo "{{.Module}}/internal/domain/user/repo"
	userService "{{.Module}}/internal/domain/user/service"
{{end}}
{{if .HasFile}}
	fileApi "{{.Module}}/internal/domain/file/api"
	fileRepo "{{.Module}}/internal/domain/file/repo"
	fileService "{{.Module}}/internal/domain/file/service"
{{end}}
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config, log *zap.SugaredLogger) {
{{if .HasAuth}}	// -----------------------
	// JWT & Auth setup
	// -----------------------
	jwtManager := authService.NewJWTManager(
		cfg.JWT.SigningKey,
		cfg.JWT.RefreshKey,
		cfg.JWT.AccessExpiresIn,
		cfg.JWT.RefreshExpiresIn,
	)
{{end}}
{{if .HasUser}}	uRepo := userRepo.NewUserRepo(db)
	uService := userService.NewUserService(uRepo, log)
	uHandler := userApi.NewUserHandler(uService, log)
{{end}}
{{if .HasAuth}}	tRepo := authRepo.NewTokenRepo(db)
	tStore := authService.NewTokenStore(tRepo, log)
	aService := authService.NewAuthService(uRepo, jwtManager, tStore, log)
	aHandler := authApi.NewAuthHandler(aService, log)

	// Start background job to clean up expired tokens every 24 hours
	go authService.StartTokenCleanupJob(tStore, 24*time.Hour)
{{end}}
{{if .HasFile}}	fRepo := fileRepo.NewFileRepo(db)
	var fileHandler *fileApi.FileHandler
	fSvc, err := fileService.NewFileService(fRepo, cfg, log)
	if err != nil {
		log.Warnf("FileService initialization failed (MinIO unavailable): %v", err)
		log.Warn("File upload/download endpoints will be unavailable")
	} else {
		fileHandler = fileApi.NewFileHandler(fSvc, log)
	}
{{end}}
	// -----------------------
	// API Versioning: v1
	// -----------------------
	v1 := r.Group("/api/v1")
	{
{{if .HasAuth}}		// -----------------------
		// Auth routes
		// -----------------------
		auth := v1.Group("/")
		{
			auth.POST("/login", aHandler.Login)
			auth.POST("/refresh", aHandler.Refresh)
			auth.POST("/logout", aHandler.Logout)
		}
{{end}}
{{if .HasUser}}		// -----------------------
		// User routes
		// -----------------------
		users := v1.Group("/users")
		{
			users.POST("/", uHandler.Register)
{{if .HasAuth}}			users.GET("/", middleware.JWTAuth(jwtManager), uHandler.ListUsers)
			users.GET("/:id", middleware.JWTAuth(jwtManager), uHandler.GetUser)
			users.PUT("/:id", middleware.JWTAuth(jwtManager), uHandler.Update)
			users.DELETE("/:id", middleware.JWTAuth(jwtManager), uHandler.Delete)
{{else}}			users.GET("/", uHandler.ListUsers)
			users.GET("/:id", uHandler.GetUser)
			users.PUT("/:id", uHandler.Update)
			users.DELETE("/:id", uHandler.Delete)
{{end}}		}
{{end}}
{{if .HasAuth}}		// -----------------------
		// Protected routes
		// -----------------------
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuth(jwtManager))
		{
			protected.GET("/me", aHandler.Me)
		}
{{end}}
{{if .HasFile}}		// -----------------------
		// File routes (only if MinIO available)
		// -----------------------
		if fSvc != nil {
			files := v1.Group("/files")
{{if .HasAuth}}			files.Use(middleware.JWTAuth(jwtManager))
{{end}}			{
				files.POST("/upload", fileHandler.Upload)
				files.GET("/:filename", fileHandler.GetFile)
				files.DELETE("/:filename", fileHandler.DeleteFile)
				files.GET("/", fileHandler.GetUserFiles)
			}
		}
{{end}}	}

	log.Info("Routes registered successfully under /api/v1")
}
`

	data := struct {
		Module      string
		HasAuth     bool
		HasUser     bool
		HasDatabase bool
		HasFile     bool
		HasDocs     bool
		HasDocker   bool
	}{
		Module:      moduleName,
		HasAuth:     selectedFeatures["Authentication (JWT)"],
		HasUser:     selectedFeatures["User Management"],
		HasDatabase: selectedFeatures["Database"],
		HasFile:     selectedFeatures["File Storage"],
		HasDocs:     selectedFeatures["API Docs"],
		HasDocker:   selectedFeatures["Docker"],
	}

	tmpl, err := template.New("routes.go").Parse(routesGoTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse routes.go template: %w", err)
	}

	routesGoPath := filepath.Join(projectDir, "internal", "app", "routes.go")

	// Create directory structure first
	if err := os.MkdirAll(filepath.Dir(routesGoPath), 0755); err != nil {
		return fmt.Errorf("failed to create internal/app directory: %w", err)
	}

	f, err := os.Create(routesGoPath)
	if err != nil {
		return fmt.Errorf("failed to create routes.go: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute routes.go template: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, content, 0600)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, content, info.Mode())
	})
}

func replaceModuleNames(projectDir, projectName, moduleName string) error {
	templateModule := "go_platform_template"
	templateName := "go-platform-template"

	// Walk through Go files
	if err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Process Go files
		if strings.HasSuffix(path, ".go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			content = []byte(strings.ReplaceAll(string(content), templateModule, moduleName))
			content = []byte(strings.ReplaceAll(string(content), templateName, projectName))

			return os.WriteFile(path, content, info.Mode())
		}

		// Process config files
		if isConfigFile(path) {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			content = []byte(strings.ReplaceAll(string(content), templateName, projectName))
			return os.WriteFile(path, content, info.Mode())
		}

		return nil
	}); err != nil {
		return err
	}

	// Update go.mod
	goModPath := filepath.Join(projectDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err == nil {
		lines := strings.Split(string(content), "\n")
		if len(lines) > 0 {
			lines[0] = fmt.Sprintf("module %s", moduleName)
		}
		_ = os.WriteFile(goModPath, []byte(strings.Join(lines, "\n")), 0600)
	}

	return nil
}

func isConfigFile(path string) bool {
	configExts := []string{".yaml", ".yml", ".json", ".toml"}
	configNames := []string{"Makefile", "Dockerfile"}

	for _, ext := range configExts {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	name := filepath.Base(path)
	for _, configName := range configNames {
		if name == configName {
			return true
		}
	}

	return false
}

func initializeGit(projectDir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = projectDir
	if err := cmd.Run(); err != nil {
		return err
	}

	cmds := [][]string{
		{"git", "config", "user.email", "dev@example.com"},
		{"git", "config", "user.name", "Developer"},
		{"git", "add", "."},
		{"git", "commit", "-m", "Initial commit: created from go-platform-template"},
	}

	for _, args := range cmds {
		//nolint:gosec
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = projectDir
		cmd.Stdout = nil
		cmd.Stderr = nil
		_ = cmd.Run() // Ignore errors
	}

	return nil
}
