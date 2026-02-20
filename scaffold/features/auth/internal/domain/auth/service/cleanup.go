package service

import (
	"context"
	"time"
)

// StartTokenCleanupJob runs a background job to cleanup expired tokens
func StartTokenCleanupJob(tokenStore *TokenStore, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			tokenStore.logger.Info("Starting expired token cleanup...")
			// Use background context for cleanup operations
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			if err := tokenStore.CleanupExpiredTokens(ctx); err != nil {
				tokenStore.logger.Errorf("Token cleanup failed: %v", err)
			}
			cancel()
		}
	}()

	tokenStore.logger.Infof("Token cleanup job started (interval: %v)", interval)
}
