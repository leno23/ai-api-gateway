package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv         string
	HTTPAddr       string
	DatabaseDSN    string
	RedisAddr      string
	JWTSecret      string
	UpstreamBase   string
	UpstreamAPIKey string
	// RateLimitGlobalPerMin is a coarse safety valve for the whole gateway (0 disables).
	RateLimitGlobalPerMin int
	// RateLimitUserPerMin limits each authenticated user on /v1 (0 disables).
	RateLimitUserPerMin int
	// QuotaReconcileSec scans Redis quota:user:* and refreshes from Postgres; 0 disables.
	QuotaReconcileSec int
	// RebateEnabled turns on post-consume rebate queue (default false).
	RebateEnabled bool
	// RebateBPS is rebate as basis points of consumed quota (e.g. 100 = 1%).
	RebateBPS int
	// RebatePollSec is worker poll interval when RebateEnabled; 0 uses 10.
	RebatePollSec int
	// PortalAPINodesJSON is JSON array of {name,url,region} for dashboard node card.
	PortalAPINodesJSON string
	// RechargeEnabled allows mock online recharge API (dev only by default).
	RechargeEnabled bool
	// PortalOrigin is used to build invite URLs in wallet API.
	PortalOrigin string
	// AffiliateRechargeBPS is inviter share on invitee recharge (e.g. 1000 = 10%).
	AffiliateRechargeBPS int
}

func Load() (*Config, error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("HTTP_ADDR", ":8080")
	viper.SetDefault("DATABASE_DSN", "host=localhost user=gateway password=gateway dbname=ai_gateway port=5432 sslmode=disable TimeZone=UTC")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("JWT_SECRET", "dev-secret-change-me-please-32chars")
	viper.SetDefault("UPSTREAM_BASE_URL", "https://api.openai.com/v1")
	viper.SetDefault("RATE_LIMIT_GLOBAL_PER_MIN", 200000)
	viper.SetDefault("RATE_LIMIT_USER_PER_MIN", 120)
	viper.SetDefault("QUOTA_RECONCILE_INTERVAL_SEC", 120)
	viper.SetDefault("REBATE_ENABLED", false)
	viper.SetDefault("REBATE_BPS", 0)
	viper.SetDefault("REBATE_POLL_SEC", 10)
	viper.SetDefault("PORTAL_API_NODES_JSON", "")
	viper.SetDefault("RECHARGE_ENABLED", false)
	viper.SetDefault("PORTAL_ORIGIN", "http://localhost:3001")
	viper.SetDefault("AFFILIATE_RECHARGE_BPS", 1000)
	_ = viper.ReadInConfig()

	return &Config{
		AppEnv:                viper.GetString("APP_ENV"),
		HTTPAddr:              viper.GetString("HTTP_ADDR"),
		DatabaseDSN:           viper.GetString("DATABASE_DSN"),
		RedisAddr:             viper.GetString("REDIS_ADDR"),
		JWTSecret:             viper.GetString("JWT_SECRET"),
		UpstreamBase:          strings.TrimRight(viper.GetString("UPSTREAM_BASE_URL"), "/"),
		UpstreamAPIKey:        viper.GetString("UPSTREAM_API_KEY"),
		RateLimitGlobalPerMin: viper.GetInt("RATE_LIMIT_GLOBAL_PER_MIN"),
		RateLimitUserPerMin:   viper.GetInt("RATE_LIMIT_USER_PER_MIN"),
		QuotaReconcileSec:     viper.GetInt("QUOTA_RECONCILE_INTERVAL_SEC"),
		RebateEnabled:         viper.GetBool("REBATE_ENABLED"),
		RebateBPS:             viper.GetInt("REBATE_BPS"),
		RebatePollSec:         viper.GetInt("REBATE_POLL_SEC"),
		PortalAPINodesJSON:    viper.GetString("PORTAL_API_NODES_JSON"),
		RechargeEnabled:       viper.GetBool("RECHARGE_ENABLED"),
		PortalOrigin:          strings.TrimRight(viper.GetString("PORTAL_ORIGIN"), "/"),
		AffiliateRechargeBPS:  viper.GetInt("AFFILIATE_RECHARGE_BPS"),
	}, nil
}
