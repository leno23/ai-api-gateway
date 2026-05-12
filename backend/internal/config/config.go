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
	}, nil
}
