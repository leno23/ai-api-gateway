package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"

	"github.com/leno23/ai-api-gateway/internal/model"
	"github.com/leno23/ai-api-gateway/internal/repository"
)

type AdminChannelDeps struct {
	Repos *repository.Repos
}

type channelPayload struct {
	Name         string            `json:"name" binding:"required"`
	Provider     string            `json:"provider" binding:"required"`
	BaseURL      string            `json:"base_url" binding:"required"`
	APIKey       string            `json:"api_key" binding:"required"`
	Models       []string          `json:"models"`
	ModelMapping map[string]string `json:"model_mapping"`
	Priority     int               `json:"priority"`
	Weight       int               `json:"weight"`
	Status       *int16            `json:"status"`
	RateLimit    *int              `json:"rate_limit"`
}

type channelListItem struct {
	ID            int64             `json:"id"`
	Name          string            `json:"name"`
	Provider      string            `json:"provider"`
	BaseURL       string            `json:"base_url"`
	APIKeyPreview string            `json:"api_key_preview"`
	Models        []string          `json:"models"`
	ModelMapping  map[string]string `json:"model_mapping,omitempty"`
	Priority      int               `json:"priority"`
	Weight        int               `json:"weight"`
	Status        int16             `json:"status"`
	RateLimit     int               `json:"rate_limit"`
}

func maskKey(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return "****" + s[len(s)-4:]
}

func ListChannels(deps AdminChannelDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := deps.Repos.ListAllChannels()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "list channels"})
			return
		}
		out := make([]channelListItem, 0, len(rows))
		for _, ch := range rows {
			var mapping map[string]string
			if len(ch.ModelMapping) > 0 {
				_ = json.Unmarshal(ch.ModelMapping, &mapping)
			}
			out = append(out, channelListItem{
				ID:            ch.ID,
				Name:          ch.Name,
				Provider:      ch.Provider,
				BaseURL:       ch.BaseURL,
				APIKeyPreview: maskKey(ch.APIKey),
				Models:        []string(ch.Models),
				ModelMapping:  mapping,
				Priority:      ch.Priority,
				Weight:        ch.Weight,
				Status:        ch.Status,
				RateLimit:     ch.RateLimit,
			})
		}
		c.JSON(http.StatusOK, gin.H{"items": out})
	}
}

func CreateChannel(deps AdminChannelDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req channelPayload
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var mm []byte
		if len(req.ModelMapping) > 0 {
			b, err := json.Marshal(req.ModelMapping)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "model_mapping"})
				return
			}
			mm = b
		}
		st := model.ChannelStatusActive
		if req.Status != nil {
			st = *req.Status
		}
		rl := 1000
		if req.RateLimit != nil {
			rl = *req.RateLimit
		}
		ch := model.Channel{
			Name:         req.Name,
			Provider:     req.Provider,
			BaseURL:      req.BaseURL,
			APIKey:       req.APIKey,
			Models:       pq.StringArray(req.Models),
			ModelMapping: mm,
			Priority:     req.Priority,
			Weight:       req.Weight,
			Status:       st,
			RateLimit:    rl,
		}
		if err := deps.Repos.CreateChannel(&ch); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "create channel"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": ch.ID})
	}
}

func UpdateChannel(deps AdminChannelDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id"})
			return
		}
		ch, err := deps.Repos.GetChannel(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		var req channelPayload
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var mm []byte
		if len(req.ModelMapping) > 0 {
			b, err := json.Marshal(req.ModelMapping)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "model_mapping"})
				return
			}
			mm = b
		}
		ch.Name = req.Name
		ch.Provider = req.Provider
		ch.BaseURL = req.BaseURL
		ch.APIKey = req.APIKey
		ch.Models = pq.StringArray(req.Models)
		ch.ModelMapping = mm
		ch.Priority = req.Priority
		ch.Weight = req.Weight
		if req.Status != nil {
			ch.Status = *req.Status
		}
		if req.RateLimit != nil {
			ch.RateLimit = *req.RateLimit
		}
		if err := deps.Repos.SaveChannel(ch); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "save"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

func DeleteChannel(deps AdminChannelDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id"})
			return
		}
		if err := deps.Repos.DeleteChannel(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}
