// internal/shared/middleware/rate_limiter.go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginmiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimitMiddleware() gin.HandlerFunc {
	rate, _ := limiter.NewRateFromFormatted("100-S") // 100 req/sec
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	return ginmiddleware.NewMiddleware(instance)
}
