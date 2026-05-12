package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ChatCompletionsLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gateway_chat_completions_duration_seconds",
			Help:    "Latency of successful /v1/chat/completions upstream round-trips.",
			Buckets: prometheus.ExponentialBuckets(0.05, 2, 12),
		},
		[]string{"stream"},
	)
	UpstreamAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_upstream_attempts_total",
			Help: "Upstream relay attempts per channel label (fallback uses env).",
		},
		[]string{"channel", "result"},
	)
)

func ChannelLabel(id *int64) string {
	if id == nil {
		return "env"
	}
	return strconv.FormatInt(*id, 10)
}
