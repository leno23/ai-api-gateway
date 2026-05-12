package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/leno23/ai-api-gateway/internal/config"
	"github.com/leno23/ai-api-gateway/internal/ratelimit"
)

// RateLimitGateway applies fixed-window Redis limits after API key auth (global + per user).
func RateLimitGateway(cfg *config.Config, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if ok, err := ratelimit.AllowFixedWindow(ctx, rdb, ratelimit.GlobalKey(), int64(cfg.RateLimitGlobalPerMin)); err == nil && !ok {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "global rate limit"})
			return
		}
		uid, ok := c.Get(CtxUserID)
		if ok {
			userID := uid.(int64)
			if ok2, err := ratelimit.AllowFixedWindow(ctx, rdb, ratelimit.UserKey(userID), int64(cfg.RateLimitUserPerMin)); err == nil && !ok2 {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "user rate limit"})
				return
			}
		}
		c.Next()
	}
}
