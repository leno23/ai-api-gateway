package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/leno23/ai-api-gateway/internal/repository"
)

const (
	CtxUserID    = "userID"
	CtxAPIKeyID  = "apiKeyID"
	CtxRole      = "role"
	CtxAPIModels = "apiKeyModels"
)

// GatewayAPIKey validates `Authorization: Bearer sk-...` against api_keys.key_hash.
func GatewayAPIKey(repos *repository.Repos) gin.HandlerFunc {
	prefixLen := 16
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		raw := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
		if !strings.HasPrefix(raw, "sk-") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key format"})
			return
		}
		userID, keyID, models, err := repos.AuthenticateGatewayAPIKey(raw, prefixLen)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrInvalidAPIKey):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			case errors.Is(err, repository.ErrUserInactive):
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user inactive"})
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "key lookup failed"})
			}
			return
		}
		c.Set(CtxUserID, userID)
		c.Set(CtxAPIKeyID, keyID)
		c.Set(CtxAPIModels, models)
		c.Next()
	}
}
