package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/leno23/ai-api-gateway/internal/config"
	"github.com/leno23/ai-api-gateway/internal/handler"
	"github.com/leno23/ai-api-gateway/internal/middleware"
	"github.com/leno23/ai-api-gateway/internal/openapi"
	"github.com/leno23/ai-api-gateway/internal/repository"
)

func New(cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	if cfg.AppEnv != "production" {
		r.Use(middleware.CORS())
	}

	repos := repository.New(db)

	authDeps := handler.AuthDeps{DB: db, JWTSecret: []byte(cfg.JWTSecret)}
	gwDeps := handler.GatewayDeps{
		Repos:         repos,
		DB:            db,
		RDB:           rdb,
		UpstreamBase:  cfg.UpstreamBase,
		UpstreamKey:   cfg.UpstreamAPIKey,
		Log:           log,
		RebateEnabled: cfg.RebateEnabled,
		RebateBPS:     cfg.RebateBPS,
	}
	keyDeps := handler.APIKeyDeps{DB: db}
	redeemDeps := handler.RedeemDeps{DB: db}
	adminDeps := handler.AdminDeps{DB: db}
	adminCh := handler.AdminChannelDeps{Repos: repos}
	adminUser := handler.AdminUserDeps{Repos: repos}
	adminRedeem := handler.AdminRedeemListDeps{Repos: repos}

	r.GET("/health", handler.Health())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/openapi.yaml", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/yaml; charset=utf-8", openapi.Spec)
	})

	r.POST("/auth/register", handler.Register(authDeps))
	r.POST("/auth/login", handler.Login(authDeps))

	u := r.Group("/user")
	u.Use(middleware.JWTAuth([]byte(cfg.JWTSecret)))
	u.POST("/api-keys", handler.CreateAPIKey(keyDeps))
	u.POST("/redeem", handler.RedeemUser(redeemDeps))

	a := r.Group("/admin")
	a.Use(middleware.JWTAuth([]byte(cfg.JWTSecret)), middleware.RequireAdmin())
	a.POST("/redeem/batch", handler.BatchRedeem(adminDeps))
	a.GET("/redeem/codes", handler.ListRedeemCodes(adminRedeem))
	a.GET("/redeem/stats", handler.RedeemStats(adminRedeem))
	a.GET("/channels", handler.ListChannels(adminCh))
	a.POST("/channels", handler.CreateChannel(adminCh))
	a.PUT("/channels/:id", handler.UpdateChannel(adminCh))
	a.DELETE("/channels/:id", handler.DeleteChannel(adminCh))
	a.PATCH("/users/:id/status", handler.PatchUserStatus(adminUser))

	v1 := r.Group("/v1")
	v1.Use(middleware.GatewayAPIKey(repos), middleware.RateLimitGateway(cfg, rdb))
	v1.GET("/models", handler.ProxyModels(gwDeps))
	v1.POST("/chat/completions", handler.ChatCompletions(gwDeps))
	v1.POST("/embeddings", handler.Placeholder("embeddings"))
	v1.POST("/images/generations", handler.Placeholder("images/generations"))

	return r
}
