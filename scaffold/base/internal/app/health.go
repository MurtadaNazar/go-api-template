package bootstrap

import (
	"net/http"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler returns a simple DB ping health check
func HealthCheckHandler(db *gorm.DB, log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := c.Get("RequestID")
		sqlDB, _ := db.DB()
		if err := sqlDB.PingContext(c.Request.Context()); err != nil {
			log.Warnw("Health check failed", "error", err, "request_id", requestID)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "db down", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
