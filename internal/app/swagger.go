package bootstrap

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go_platform_template/internal/platform/config"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// SetupSwagger initializes Swagger documentation and UI for the application.
// Only active in debug/development modes. In production mode, this function returns early
// without registering any routes.
//
// It performs the following operations:
//   - Generates initial Swagger documentation asynchronously
//   - Sets up file watching for auto-regeneration when API code changes
//   - Registers Swagger UI routes with Gin engine
//
// Parameters:
//   - r: Gin engine instance where Swagger routes will be registered
//   - cfg: Application configuration containing server settings and Gin mode
//   - logger: Logger for swagger operations
//
// Routes registered:
//   - /swagger/*any: Serves Swagger UI and JSON documentation
func SetupSwagger(r *gin.Engine, cfg *config.Config, logger *zap.SugaredLogger) {
	// Only setup swagger in debug/development modes
	if cfg.GinMode != "debug" && cfg.GinMode != "development" {
		return
	}

	// Swagger URL points to the generated JSON (docs.go imports generated swagger from docs/)
	url := ginSwagger.URL("/swagger/doc.json")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Generate docs in background
	go func() {
		if err := generateSwagger(); err != nil {
			logger.Warnf("Failed to generate Swagger docs: %v. Install 'swag' with: go install github.com/swaggo/swag/cmd/swag@latest", err)
		} else {
			logger.Info("Swagger docs generated successfully")
		}
	}()

	// Start file watcher for auto-regeneration
	go watchSwaggerBatch(time.Second, "./internal", logger)
}

// generateSwagger executes the `swag init` command to generate Swagger documentation.
// It configures the command to:
//   - Use cmd/server/main.go as the entry point (-g flag)
//   - Output generated files to the docs directory (-o flag)
//   - Suppress command output
//
// Returns:
//   - error: Any error encountered during command execution, nil on success
func generateSwagger() error {
	cmd := exec.Command("swag", "init", "-g", "cmd/server/main.go", "-o", "docs")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

// watchSwaggerBatch monitors Go files in API and model directories for changes
// and triggers Swagger documentation regeneration with debouncing.
//
// The function:
//   - Recursively watches directories ending with /api or /model
//   - Uses debouncing to batch multiple rapid file changes
//   - Only regenerates documentation for files containing Swagger annotations
//   - Dynamically adds new matching directories as they are created
//
// Parameters:
//   - debounceDuration: Minimum time between regeneration triggers
//   - baseDir: Root directory to start watching for API/model directories
//   - logger: Logger for watcher operations
func watchSwaggerBatch(debounceDuration time.Duration, baseDir string, logger *zap.SugaredLogger) {
	// Initialize file system watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Warnf("Failed to create file watcher: %v", err)
		return
	}
	defer watcher.Close()

	var mu sync.Mutex
	watchedDirs := make(map[string]struct{}) // Tracks already watched directories

	// addWatcher recursively finds and watches API/model directories
	addWatcher := func(dir string) {
		if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Continue walking on error
			}
			// Check if directory matches API/model patterns
			if info.IsDir() &&
				(strings.HasSuffix(path, "/api") || strings.HasSuffix(path, "/model") ||
					strings.Contains(path, "/api/") || strings.Contains(path, "/model/") || strings.Contains(path, "/dto/")) {
				mu.Lock()
				// Add to watcher if not already watching
				if _, ok := watchedDirs[path]; !ok {
					if err := watcher.Add(path); err == nil {
						watchedDirs[path] = struct{}{}
						logger.Infof("Watching directory: %s", path)
					}
				}
				mu.Unlock()
			}
			return nil
		}); err != nil {
			logger.Warnf("Error walking directory %s: %v", dir, err)
		}
	}

	// Start watching initial base directory
	addWatcher(baseDir)

	// Debouncing setup
	trigger := make(chan struct{}, 1) // Buffered channel to prevent blocking
	var debounceTimer *time.Timer

	// Regeneration handler - waits for trigger and regenerates after debounce period
	go func() {
		for range trigger {
			// Stop existing timer if any
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			// Create new timer - will regenerate after debounceDuration
			debounceTimer = time.AfterFunc(debounceDuration, func() {
				if err := generateSwagger(); err != nil {
					logger.Warnf("Failed to regenerate Swagger docs: %v", err)
				} else {
					logger.Info("Swagger docs regenerated")
				}
			})
		}
	}()

	// Main event loop
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return // Channel closed
			}

			// Dynamically add new directories that match our patterns
			info, err := os.Stat(event.Name)
			if err == nil && info.IsDir() {
				addWatcher(event.Name)
			}

			// Trigger regeneration for relevant Go files
			if strings.HasSuffix(event.Name, ".go") && !strings.HasSuffix(event.Name, "_test.go") {
				if containsSwaggerAnnotations(event.Name) {
					select {
					case trigger <- struct{}{}: // Schedule regeneration
					default: // Already scheduled, skip duplicate
					}
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return // Channel closed
			}
			logger.Warnf("File watcher error: %v", err)
		}
	}
}

// containsSwaggerAnnotations scans a Go source file for Swagger annotation comments.
// Swagger annotations are identified by lines starting with "// @".
//
// Parameters:
//   - filePath: Path to the Go file to scan
//
// Returns:
//   - bool: true if file contains at least one Swagger annotation, false otherwise
func containsSwaggerAnnotations(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "// @") {
			return true
		}
	}
	return false
}
