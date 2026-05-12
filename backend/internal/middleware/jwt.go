package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/leno23/ai-api-gateway/internal/auth"
)

// JWTAuth validates `Authorization: Bearer <jwt>` for non-sk keys.
func JWTAuth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		raw := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
		if strings.HasPrefix(raw, "sk-") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "use user login token"})
			return
		}
		cl, err := auth.ParseJWT(secret, raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set(CtxUserID, cl.UserID)
		c.Set(CtxRole, cl.Role)
		c.Next()
	}
}

// RequireAdmin ensures JWT role is admin or super-admin.
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		r, ok := c.Get(CtxRole)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		role, ok := r.(int16)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if !auth.IsAdminRole(role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		c.Next()
	}
}
