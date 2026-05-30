package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/leno23/ai-api-gateway/internal/middleware"
	"github.com/leno23/ai-api-gateway/internal/model"
)

const ctxLogTokenNameOverride = "log_token_name_override"

type playgroundChatReq struct {
	Model              string          `json:"model" binding:"required"`
	Messages           json.RawMessage `json:"messages" binding:"required"`
	Stream             *bool           `json:"stream"`
	MaxTokens          *int            `json:"max_tokens"`
	Temperature        *float64        `json:"temperature"`
	TopP               *float64        `json:"top_p"`
	FrequencyPenalty   *float64        `json:"frequency_penalty"`
	PresencePenalty    *float64        `json:"presence_penalty"`
	APIKeyID           *int64          `json:"api_key_id"`
	TokenGroup         string          `json:"token_group"`
	CustomBody         bool            `json:"custom_body"`
	CustomBodyRaw      json.RawMessage `json:"custom_body_raw"`
}

func PlaygroundChat(deps GatewayDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var pg playgroundChatReq
		if err := c.ShouldBindJSON(&pg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		apiKeyID, tokenOverride, allowedModels, err := resolvePlaygroundAPIKey(deps.DB, userID, pg.APIKeyID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token"})
			return
		}

		upstreamBody, err := buildPlaygroundUpstreamBody(pg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Set(middleware.CtxUserID, userID)
		c.Set(middleware.CtxAPIKeyID, apiKeyID)
		if tokenOverride != "" {
			c.Set(ctxLogTokenNameOverride, tokenOverride)
		}
		if len(allowedModels) > 0 {
			c.Set(middleware.CtxAPIModels, allowedModels)
		}

		c.Request.Body = io.NopCloser(bytes.NewReader(upstreamBody))
		c.Request.ContentLength = int64(len(upstreamBody))
		c.Request.Header.Set("Content-Type", "application/json")

		ChatCompletions(deps)(c)
	}
}

func resolvePlaygroundAPIKey(db *gorm.DB, userID int64, requestedID *int64) (apiKeyID int64, tokenOverride string, allowedModels []string, err error) {
	base := db.Where("user_id = ? AND status = ?", userID, model.APIKeyStatusActive)
	if requestedID != nil && *requestedID > 0 {
		var key model.APIKey
		if err := base.Where("id = ?", *requestedID).First(&key).Error; err != nil {
			return 0, "", nil, err
		}
		return key.ID, "", []string(key.Models), nil
	}
	var key model.APIKey
	if err := base.Order("id ASC").First(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, "playground-default", nil, nil
		}
		return 0, "", nil, err
	}
	return key.ID, "", []string(key.Models), nil
}

func buildPlaygroundUpstreamBody(pg playgroundChatReq) ([]byte, error) {
	if pg.CustomBody && len(pg.CustomBodyRaw) > 0 {
		var probe struct {
			Model string `json:"model"`
		}
		if err := json.Unmarshal(pg.CustomBodyRaw, &probe); err != nil || probe.Model == "" {
			return nil, errors.New("custom body must include model")
		}
		return pg.CustomBodyRaw, nil
	}

	stream := true
	if pg.Stream != nil {
		stream = *pg.Stream
	}
	payload := map[string]any{
		"model":    pg.Model,
		"messages": json.RawMessage(pg.Messages),
		"stream":   stream,
	}
	if pg.MaxTokens != nil {
		payload["max_tokens"] = *pg.MaxTokens
	}
	if pg.Temperature != nil {
		payload["temperature"] = *pg.Temperature
	}
	if pg.TopP != nil {
		payload["top_p"] = *pg.TopP
	}
	if pg.FrequencyPenalty != nil {
		payload["frequency_penalty"] = *pg.FrequencyPenalty
	}
	if pg.PresencePenalty != nil {
		payload["presence_penalty"] = *pg.PresencePenalty
	}
	return json.Marshal(payload)
}

func logTokenNameOverride(c *gin.Context) string {
	v, ok := c.Get(ctxLogTokenNameOverride)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func apiKeyIDPtr(id int64) *int64 {
	if id <= 0 {
		return nil
	}
	v := id
	return &v
}
