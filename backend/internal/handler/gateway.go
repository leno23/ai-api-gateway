package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/leno23/ai-api-gateway/internal/audit"
	"github.com/leno23/ai-api-gateway/internal/circuit"
	"github.com/leno23/ai-api-gateway/internal/metrics"
	"github.com/leno23/ai-api-gateway/internal/middleware"
	"github.com/leno23/ai-api-gateway/internal/model"
	"github.com/leno23/ai-api-gateway/internal/ratelimit"
	"github.com/leno23/ai-api-gateway/internal/rebate"
	"github.com/leno23/ai-api-gateway/internal/relay"
	"github.com/leno23/ai-api-gateway/internal/repository"
	"github.com/leno23/ai-api-gateway/internal/routing"
	"github.com/leno23/ai-api-gateway/internal/service"
	"github.com/leno23/ai-api-gateway/internal/tokenizer"
)

type GatewayDeps struct {
	Repos         *repository.Repos
	DB            *gorm.DB
	RDB           *redis.Client
	UpstreamBase  string
	UpstreamKey   string
	Log           *zap.Logger
	RebateEnabled bool
	RebateBPS     int
}

type chatReqLite struct {
	Model     string          `json:"model"`
	MaxTokens *int            `json:"max_tokens"`
	Stream    bool            `json:"stream"`
	Messages  json.RawMessage `json:"messages"`
}

type openAIUsage struct {
	Usage struct {
		PromptTokens     int64 `json:"prompt_tokens"`
		CompletionTokens int64 `json:"completion_tokens"`
		TotalTokens      int64 `json:"total_tokens"`
	} `json:"usage"`
}

func ModelsList(repos *repository.Repos) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := repos.ListModelPrices()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "models"})
			return
		}
		out := make([]gin.H, 0, len(rows))
		for _, r := range rows {
			out = append(out, gin.H{
				"id":       r.Model,
				"object":   "model",
				"owned_by": "gateway",
			})
		}
		c.JSON(http.StatusOK, gin.H{"object": "list", "data": out})
	}
}

func shouldFailoverUpstream(status int) bool {
	return status >= 500 || status == http.StatusTooManyRequests
}

func ChatCompletions(deps GatewayDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "read body"})
			return
		}
		var req chatReqLite
		if err := json.Unmarshal(body, &req); err != nil || req.Model == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json or model"})
			return
		}
		if allowed, ok := c.Get(middleware.CtxAPIModels); ok {
			if list, ok := allowed.([]string); ok && len(list) > 0 {
				found := false
				for _, m := range list {
					if m == req.Model {
						found = true
						break
					}
				}
				if !found {
					c.JSON(http.StatusForbidden, gin.H{"error": "model not allowed for this key"})
					return
				}
			}
		}
		uid, _ := c.Get(middleware.CtxUserID)
		userID := uid.(int64)
		apiKeyIDVal, _ := c.Get(middleware.CtxAPIKeyID)
		apiKeyID := apiKeyIDVal.(int64)

		targets, err := routing.BuildUpstreamTargets(deps.Repos, deps.UpstreamBase, deps.UpstreamKey, req.Model)
		if err != nil {
			deps.Log.Warn("upstream targets", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "routing"})
			return
		}
		if len(targets) == 0 {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no upstream configured"})
			return
		}

		maxOut := int64(1024)
		if req.MaxTokens != nil && *req.MaxTokens > 0 {
			maxOut = int64(*req.MaxTokens)
		}
		promptTok := tokenizer.EstimatePromptTokens(req.Messages)
		if promptTok < 1 {
			promptTok = 1
		}
		estimate, err := service.EstimateMaxCost(deps.DB, req.Model, promptTok, maxOut)
		if err != nil {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "unknown model pricing", "code": "insufficient_quota"})
			return
		}
		if estimate <= 0 {
			estimate = 1
		}

		bal, err := service.PreDeduct(ctx, deps.RDB, deps.DB, userID, estimate)
		if err != nil {
			deps.Log.Warn("pre-deduct", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "quota"})
			return
		}
		if bal < 0 {
			_ = service.RefundQuota(ctx, deps.RDB, userID, estimate)
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "insufficient quota", "code": "insufficient_quota"})
			return
		}

		var resp *http.Response
		var picked routing.UpstreamTarget
		var pickedOK bool

		for _, t := range targets {
			chLabel := metrics.ChannelLabel(t.ChannelID)
			if t.ChannelID != nil {
				if !circuit.Allow(ctx, deps.RDB, *t.ChannelID) {
					metrics.UpstreamAttempts.WithLabelValues(chLabel, "circuit_open").Inc()
					continue
				}
				okRL, _ := ratelimit.AllowFixedWindow(ctx, deps.RDB, ratelimit.ChannelKey(*t.ChannelID), int64(t.RateLimit))
				if !okRL {
					metrics.UpstreamAttempts.WithLabelValues(chLabel, "ratelimit").Inc()
					continue
				}
			}

			body2, err := routing.RemapChatBody(body, req.Model, t.ModelMap)
			if err != nil {
				deps.Log.Warn("remap model", zap.Error(err))
				metrics.UpstreamAttempts.WithLabelValues(chLabel, "remap_error").Inc()
				continue
			}

			rc := relay.New(t.BaseURL, t.APIKey)
			if req.Stream {
				rc.HTTP = &http.Client{Transport: http.DefaultTransport}
			}
			upResp, err := rc.ChatCompletions(ctx, body2, req.Stream)
			if err != nil {
				deps.Log.Warn("upstream dial", zap.Error(err), zap.String("channel", chLabel))
				metrics.UpstreamAttempts.WithLabelValues(chLabel, "transport_error").Inc()
				if t.ChannelID != nil {
					circuit.RecordFailure(ctx, deps.RDB, *t.ChannelID)
				}
				continue
			}
			if shouldFailoverUpstream(upResp.StatusCode) {
				_, _ = io.Copy(io.Discard, upResp.Body)
				_ = upResp.Body.Close()
				metrics.UpstreamAttempts.WithLabelValues(chLabel, "upstream_retryable").Inc()
				if t.ChannelID != nil {
					circuit.RecordFailure(ctx, deps.RDB, *t.ChannelID)
				}
				continue
			}

			resp = upResp
			picked = t
			pickedOK = true
			metrics.UpstreamAttempts.WithLabelValues(chLabel, "selected").Inc()
			if t.ChannelID != nil {
				circuit.RecordSuccess(ctx, deps.RDB, *t.ChannelID)
			}
			break
		}

		if !pickedOK || resp == nil {
			_ = service.RefundQuota(ctx, deps.RDB, userID, estimate)
			c.JSON(http.StatusBadGateway, gin.H{"error": "upstream unavailable"})
			return
		}

		if resp.StatusCode >= 400 {
			b, _ := io.ReadAll(resp.Body)
			relay.DrainClose(resp)
			_ = service.RefundQuota(ctx, deps.RDB, userID, estimate)
			c.Status(resp.StatusCode)
			c.Header("Content-Type", "application/json")
			_, _ = c.Writer.Write(b)
			return
		}

		start := time.Now()
		if req.Stream {
			err := pipeSSE(c, resp)
			lat := int(time.Since(start).Milliseconds())
			streamLabel := "true"
			if err != nil {
				_ = service.RefundQuota(ctx, deps.RDB, userID, estimate)
				metrics.ChatCompletionsLatency.WithLabelValues(streamLabel).Observe(time.Since(start).Seconds())
				deps.Log.Warn("stream", zap.Error(err))
				return
			}
			metrics.ChatCompletionsLatency.WithLabelValues(streamLabel).Observe(time.Since(start).Seconds())
			if err := settleStream(ctx, deps, userID, apiKeyID, picked.ChannelID, req.Model, estimate, maxOut, promptTok, lat, resp.StatusCode); err != nil {
				deps.Log.Warn("settle stream", zap.Error(err))
			}
			return
		}

		respBody, err := io.ReadAll(resp.Body)
		relay.DrainClose(resp)
		if err != nil {
			_ = service.RefundQuota(ctx, deps.RDB, userID, estimate)
			c.JSON(http.StatusBadGateway, gin.H{"error": "read upstream"})
			return
		}
		lat := int(time.Since(start).Milliseconds())
		metrics.ChatCompletionsLatency.WithLabelValues("false").Observe(time.Since(start).Seconds())

		var u openAIUsage
		_ = json.Unmarshal(respBody, &u)
		inTok := u.Usage.PromptTokens
		outTok := u.Usage.CompletionTokens
		if inTok == 0 && outTok == 0 {
			outTok = maxOut
			inTok = promptTok
		}
		actual, err := service.ActualCost(deps.DB, req.Model, inTok, outTok)
		if err != nil || actual <= 0 {
			actual = estimate
		}
		if actual > estimate {
			extra := actual - estimate
			b2, err := service.PreDeduct(ctx, deps.RDB, deps.DB, userID, extra)
			if err != nil || b2 < 0 {
				_ = service.RefundQuota(ctx, deps.RDB, userID, estimate)
				c.JSON(http.StatusPaymentRequired, gin.H{"error": "insufficient quota for actual usage"})
				return
			}
		} else if diff := estimate - actual; diff > 0 {
			_ = service.RefundQuota(ctx, deps.RDB, userID, diff)
		}

		if err := persistUsage(ctx, deps, userID, apiKeyID, picked.ChannelID, req.Model, actual, int(inTok), int(outTok), lat, resp.StatusCode, ""); err != nil {
			deps.Log.Warn("persist", zap.Error(err))
		}
		c.Status(resp.StatusCode)
		c.Header("Content-Type", resp.Header.Get("Content-Type"))
		if c.Writer.Header().Get("Content-Type") == "" {
			c.Header("Content-Type", "application/json")
		}
		_, _ = c.Writer.Write(respBody)
	}
}

func pipeSSE(c *gin.Context, resp *http.Response) error {
	defer func() { _ = resp.Body.Close() }()
	c.Status(resp.StatusCode)
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		c.Header("Content-Type", ct)
	}
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		_, err := io.Copy(c.Writer, resp.Body)
		return err
	}
	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := c.Writer.Write(buf[:n]); werr != nil {
				return werr
			}
			flusher.Flush()
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func settleStream(
	ctx context.Context,
	deps GatewayDeps,
	userID, apiKeyID int64,
	channelID *int64,
	modelName string,
	estimate, maxOut, promptTok int64,
	latencyMs, status int,
) error {
	actual := estimate
	if err := deps.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.User{}).Where("id = ? AND quota >= ?", userID, actual).
			Updates(map[string]any{
				"quota":      gorm.Expr("quota - ?", actual),
				"used_quota": gorm.Expr("used_quota + ?", actual),
			}).Error; err != nil {
			return err
		}
		var q int64
		if err := tx.Model(&model.User{}).Select("quota").Where("id = ?", userID).Scan(&q).Error; err != nil {
			return err
		}
		return tx.Create(&model.QuotaLog{
			UserID:    userID,
			Delta:     -actual,
			Balance:   q,
			Type:      model.QuotaLogTypeConsume,
			Reference: "",
			Remark:    "chat stream",
		}).Error
	}); err != nil {
		return err
	}
	var u model.User
	_ = deps.DB.Select("quota").Where("id = ?", userID).First(&u).Error
	if err := service.SyncQuotaToRedis(ctx, deps.RDB, userID, u.Quota); err != nil {
		return err
	}
	sc := status
	l := latencyMs
	pt := int(promptTok)
	co := int(maxOut)
	tt := pt + co
	rec := &model.RequestLog{
		UserID:           userID,
		APIKeyID:         &apiKeyID,
		ChannelID:        channelID,
		Model:            modelName,
		RequestMethod:    http.MethodPost,
		RequestPath:      "/v1/chat/completions",
		PromptTokens:     pt,
		CompletionTokens: co,
		TotalTokens:      tt,
		CostQuota:        actual,
		LatencyMs:        &l,
		StatusCode:       &sc,
	}
	audit.SubmitRequestLog(deps.DB, deps.Log, rec)
	rebate.EnqueueConsumeRebate(deps.DB, deps.Log, userID, actual, rebate.RefChat(modelName), deps.RebateEnabled, deps.RebateBPS)
	return nil
}

func persistUsage(
	ctx context.Context,
	deps GatewayDeps,
	userID, apiKeyID int64,
	channelID *int64,
	modelName string,
	actual int64,
	inTok, outTok, latencyMs, status int,
	errMsg string,
) error {
	err := deps.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.User{}).Where("id = ? AND quota >= ?", userID, actual).
			Updates(map[string]any{
				"quota":      gorm.Expr("quota - ?", actual),
				"used_quota": gorm.Expr("used_quota + ?", actual),
			}).Error; err != nil {
			return err
		}
		var q int64
		if err := tx.Model(&model.User{}).Select("quota").Where("id = ?", userID).Scan(&q).Error; err != nil {
			return err
		}
		return tx.Create(&model.QuotaLog{
			UserID:    userID,
			Delta:     -actual,
			Balance:   q,
			Type:      model.QuotaLogTypeConsume,
			Reference: "",
			Remark:    "chat",
		}).Error
	})
	if err != nil {
		return err
	}
	var u model.User
	_ = deps.DB.Select("quota").Where("id = ?", userID).First(&u).Error
	if err := service.SyncQuotaToRedis(ctx, deps.RDB, userID, u.Quota); err != nil {
		return err
	}
	sc := status
	l := latencyMs
	tt := inTok + outTok
	rec := &model.RequestLog{
		UserID:           userID,
		APIKeyID:         &apiKeyID,
		ChannelID:        channelID,
		Model:            modelName,
		RequestMethod:    http.MethodPost,
		RequestPath:      "/v1/chat/completions",
		PromptTokens:     inTok,
		CompletionTokens: outTok,
		TotalTokens:      tt,
		CostQuota:        actual,
		LatencyMs:        &l,
		StatusCode:       &sc,
		ErrorMessage:     errMsg,
	}
	audit.SubmitRequestLog(deps.DB, deps.Log, rec)
	rebate.EnqueueConsumeRebate(deps.DB, deps.Log, userID, actual, rebate.RefChat(modelName), deps.RebateEnabled, deps.RebateBPS)
	return nil
}

func Placeholder(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": name + " not implemented"})
	}
}

func StubModelsUpstream(deps GatewayDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, deps.UpstreamBase+"/models", bytes.NewReader(nil))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "build request"})
			return
		}
		if deps.UpstreamKey != "" {
			req.Header.Set("Authorization", "Bearer "+deps.UpstreamKey)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "upstream"})
			return
		}
		defer relay.DrainClose(resp)
		b, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), b)
	}
}

// ProxyModels returns DB-backed list by default; use ?source=upstream to forward to env upstream.
func ProxyModels(deps GatewayDeps) gin.HandlerFunc {
	dbList := ModelsList(deps.Repos)
	stub := StubModelsUpstream(deps)
	return func(c *gin.Context) {
		if c.Query("source") == "upstream" && deps.UpstreamKey != "" {
			stub(c)
			return
		}
		dbList(c)
	}
}
