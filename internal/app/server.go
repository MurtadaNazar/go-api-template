package bootstrap

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StartServer runs the Gin server and handles graceful shutdown
func StartServer(r *gin.Engine, addr string, db *gorm.DB, log *zap.SugaredLogger) {
	go func() {
		if err := r.Run(addr); err != nil {
			log.Errorf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	sqlDB, _ := db.DB()
	_ = sqlDB.Close()
	log.Info("Server stopped gracefully")
}
