package repository

import (
	"time"

	"github.com/leno23/ai-api-gateway/internal/model"
)

type TimePoint struct {
	T string `json:"t"`
	V int64  `json:"v"`
}

type ModelAgg struct {
	Model string `json:"model"`
	Cost  int64  `json:"cost"`
	Count int64  `json:"count"`
}

type HourlyAgg struct {
	Hour  string `json:"hour"`
	Count int64  `json:"count"`
	Cost  int64  `json:"cost"`
	Tokens int64 `json:"tokens"`
}

func parseRange(rangeKey string) time.Duration {
	switch rangeKey {
	case "24h", "1d":
		return 24 * time.Hour
	case "30d":
		return 30 * 24 * time.Hour
	default:
		return 7 * 24 * time.Hour
	}
}

func (r *Repos) DashboardStats(userID int64) (map[string]any, error) {
	var u model.User
	if err := r.DB.First(&u, userID).Error; err != nil {
		return nil, err
	}
	since24h := time.Now().Add(-24 * time.Hour)
	since1m := time.Now().Add(-time.Minute)

	type agg struct {
		Cnt        int64
		Tokens     int64
		Cost       int64
		AvgLatency float64
	}
	var a24 agg
	_ = r.DB.Model(&model.RequestLog{}).
		Select(`COUNT(*) AS cnt, COALESCE(SUM(total_tokens),0) AS tokens, COALESCE(SUM(cost_quota),0) AS cost, COALESCE(AVG(latency_ms),0) AS avg_latency`).
		Where("user_id = ? AND created_at >= ?", userID, since24h).
		Scan(&a24).Error

	var rpm int64
	_ = r.DB.Model(&model.RequestLog{}).Where("user_id = ? AND created_at >= ?", userID, since1m).Count(&rpm).Error

	var tpm int64
	_ = r.DB.Model(&model.RequestLog{}).
		Where("user_id = ? AND created_at >= ?", userID, since1m).
		Select("COALESCE(SUM(total_tokens),0)").Scan(&tpm).Error

	var sparkReq []TimePoint
	_ = r.DB.Raw(`
		SELECT to_char(date_trunc('hour', created_at), 'YYYY-MM-DD"T"HH24:00:00Z') AS t,
		       COUNT(*)::bigint AS v
		FROM request_logs
		WHERE user_id = ? AND created_at >= ?
		GROUP BY 1 ORDER BY 1
	`, userID, since24h).Scan(&sparkReq).Error

	var sparkCost []TimePoint
	_ = r.DB.Raw(`
		SELECT to_char(date_trunc('hour', created_at), 'YYYY-MM-DD"T"HH24:00:00Z') AS t,
		       COALESCE(SUM(cost_quota),0)::bigint AS v
		FROM request_logs
		WHERE user_id = ? AND created_at >= ?
		GROUP BY 1 ORDER BY 1
	`, userID, since24h).Scan(&sparkCost).Error

	groupSlug := "default"
	if u.TokenGroupID != nil {
		var g model.TokenGroup
		if err := r.DB.First(&g, *u.TokenGroupID).Error; err == nil {
			groupSlug = g.Slug
		}
	}

	return map[string]any{
		"account": map[string]any{
			"quota":      u.Quota,
			"used_quota": u.UsedQuota,
			"group":      groupSlug,
		},
		"usage": map[string]any{
			"request_count_24h": a24.Cnt,
			"total_tokens_24h":  a24.Tokens,
		},
		"consumption": map[string]any{
			"cost_quota_24h": a24.Cost,
		},
		"performance": map[string]any{
			"avg_latency_ms": int64(a24.AvgLatency),
			"rpm":            rpm,
			"tpm":            tpm,
		},
		"sparklines": map[string]any{
			"requests": sparkReq,
			"cost":     sparkCost,
		},
	}, nil
}

func (r *Repos) DashboardCharts(userID int64, rangeKey string) (map[string]any, error) {
	since := time.Now().Add(-parseRange(rangeKey))

	var byModel []ModelAgg
	_ = r.DB.Raw(`
		SELECT model, COALESCE(SUM(cost_quota),0)::bigint AS cost, COUNT(*)::bigint AS count
		FROM request_logs
		WHERE user_id = ? AND created_at >= ? AND model <> ''
		GROUP BY model ORDER BY cost DESC LIMIT 20
	`, userID, since).Scan(&byModel).Error

	var trend []HourlyAgg
	_ = r.DB.Raw(`
		SELECT to_char(date_trunc('hour', created_at), 'YYYY-MM-DD"T"HH24:00:00Z') AS hour,
		       COUNT(*)::bigint AS count,
		       COALESCE(SUM(cost_quota),0)::bigint AS cost,
		       COALESCE(SUM(total_tokens),0)::bigint AS tokens
		FROM request_logs
		WHERE user_id = ? AND created_at >= ?
		GROUP BY 1 ORDER BY 1
	`, userID, since).Scan(&trend).Error

	return map[string]any{
		"consumption_by_model": byModel,
		"call_trend":           trend,
		"call_distribution":    byModel,
		"ranking":              byModel,
	}, nil
}
