package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/leno23/ai-api-gateway/internal/middleware"
	"github.com/leno23/ai-api-gateway/internal/model"
	"github.com/leno23/ai-api-gateway/internal/repository"
	"github.com/leno23/ai-api-gateway/internal/service"
)

type TokenDeps struct {
	DB    *gorm.DB
	Repos *repository.Repos
}

type tokenCreateReq struct {
	Name          string   `json:"name" binding:"required,min=1,max=128"`
	TokenGroup    string   `json:"token_group"`
	QuotaLimit    *int64   `json:"quota_limit"`
	Models        []string `json:"models"`
	IPWhitelist   []string `json:"ip_whitelist"`
}

type tokenUpdateReq struct {
	Name             *string  `json:"name"`
	Status           *int16   `json:"status"`
	TokenGroup       *string  `json:"token_group"`
	QuotaLimit       *int64   `json:"quota_limit"`
	ClearQuotaLimit  bool     `json:"clear_quota_limit"`
	Models           []string `json:"models"`
	IPWhitelist      []string `json:"ip_whitelist"`
}

type batchDeleteReq struct {
	IDs []int64 `json:"ids" binding:"required,min=1"`
}

func maskAPIKey(prefix string) string {
	if len(prefix) <= 10 {
		return prefix + "******"
	}
	return prefix[:10] + "******"
}

func tokenGroupSlug(db *gorm.DB, groupID *int64) string {
	if groupID == nil {
		return "default"
	}
	var g model.TokenGroup
	if err := db.First(&g, *groupID).Error; err != nil {
		return "default"
	}
	return g.Slug
}

func resolveTokenGroupID(db *gorm.DB, slug string) (*int64, error) {
	if slug == "" {
		slug = "default"
	}
	var g model.TokenGroup
	if err := db.Where("slug = ?", slug).First(&g).Error; err != nil {
		return nil, err
	}
	id := g.ID
	return &id, nil
}

func serializeToken(db *gorm.DB, k model.APIKey) gin.H {
	return gin.H{
		"id":           k.ID,
		"name":         k.Name,
		"status":       k.Status,
		"enabled":      k.Status == model.APIKeyStatusActive,
		"key_prefix":   k.KeyPrefix,
		"key_masked":   maskAPIKey(k.KeyPrefix),
		"token_group":  tokenGroupSlug(db, k.TokenGroupID),
		"quota_limit":  k.QuotaLimit,
		"used_quota":   k.UsedQuota,
		"models":       []string(k.Models),
		"ip_whitelist": []string(k.IPWhitelist),
		"rate_limit":   k.RateLimit,
		"created_at":   k.CreatedAt,
	}
}

func ctxUserID(c *gin.Context) (int64, bool) {
	uid, ok := c.Get(middleware.CtxUserID)
	if !ok {
		return 0, false
	}
	id, ok := uid.(int64)
	return id, ok
}

func ListTokens(deps TokenDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		rows, total, err := deps.Repos.ListAPIKeysByUser(userID, page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "list"})
			return
		}
		items := make([]gin.H, 0, len(rows))
		for _, k := range rows {
			items = append(items, serializeToken(deps.DB, k))
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"items":     items,
				"total":     total,
				"page":      page,
				"page_size": pageSize,
			},
		})
	}
}

func CreateToken(deps TokenDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var req tokenCreateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		groupID, err := resolveTokenGroupID(deps.DB, req.TokenGroup)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token group"})
			return
		}
		raw, err := service.RandomRedeemCode()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rng"})
			return
		}
		plain := "sk-" + raw
		hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "hash"})
			return
		}
		prefix := plain
		if len(prefix) > 16 {
			prefix = prefix[:16]
		}
		k := model.APIKey{
			UserID:       userID,
			Name:         req.Name,
			KeyHash:      string(hash),
			KeyPrefix:    prefix,
			Status:       model.APIKeyStatusActive,
			TokenGroupID: groupID,
			QuotaLimit:   req.QuotaLimit,
			Models:       pq.StringArray(req.Models),
			IPWhitelist:  pq.StringArray(req.IPWhitelist),
		}
		if err := deps.DB.Create(&k).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "persist"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data": gin.H{
				"token":   serializeToken(deps.DB, k),
				"api_key": plain,
			},
		})
	}
}

func UpdateToken(deps TokenDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		keyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		k, err := deps.Repos.GetAPIKeyForUser(userID, keyID)
		if err != nil {
			if err == repository.ErrAPIKeyNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "lookup"})
			return
		}
		var req tokenUpdateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updates := map[string]any{}
		if req.Name != nil {
			updates["name"] = *req.Name
		}
		if req.Status != nil {
			updates["status"] = *req.Status
		}
		if req.TokenGroup != nil {
			gid, err := resolveTokenGroupID(deps.DB, *req.TokenGroup)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token group"})
				return
			}
			updates["token_group_id"] = gid
		}
		if req.ClearQuotaLimit {
			updates["quota_limit"] = nil
		} else if req.QuotaLimit != nil {
			updates["quota_limit"] = *req.QuotaLimit
		}
		if req.Models != nil {
			updates["models"] = pq.StringArray(req.Models)
		}
		if req.IPWhitelist != nil {
			updates["ip_whitelist"] = pq.StringArray(req.IPWhitelist)
		}
		if len(updates) > 0 {
			if err := deps.DB.Model(k).Updates(updates).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "update"})
				return
			}
		}
		var fresh model.APIKey
		_ = deps.DB.First(&fresh, k.ID)
		c.JSON(http.StatusOK, gin.H{"success": true, "data": serializeToken(deps.DB, fresh)})
	}
}

func DeleteToken(deps TokenDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		keyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if _, err := deps.Repos.GetAPIKeyForUser(userID, keyID); err != nil {
			if err == repository.ErrAPIKeyNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "lookup"})
			return
		}
		if err := deps.DB.Delete(&model.APIKey{}, keyID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func BatchDeleteTokens(deps TokenDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var req batchDeleteReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		n, err := deps.Repos.DeleteAPIKeysForUser(userID, req.IDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "deleted": n})
	}
}
