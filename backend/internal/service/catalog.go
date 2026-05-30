package service

import (
	"math"

	"github.com/leno23/ai-api-gateway/internal/model"
)

// Quota per 1M tokens = per-1K price * 1000 (internal integer quota units).
func quotaPerMillion(per1K int64) int64 {
	return per1K * 1000
}

func applyMultiplier(v int64, multiplier float64) int64 {
	if multiplier <= 0 {
		multiplier = 1
	}
	return int64(math.Round(float64(v) * multiplier))
}

type ModelPriceBreakdown struct {
	InputPerMillion       int64   `json:"input_per_million"`
	OutputPerMillion      int64   `json:"output_per_million"`
	CacheReadPerMillion   int64   `json:"cache_read_per_million"`
	CacheWritePerMillion  int64   `json:"cache_write_per_million"`
	UnitPrice             int64   `json:"unit_price"`
	Multiplier            float64 `json:"multiplier"`
}

func ModelPriceBreakdownFrom(mp model.ModelPrice, multiplier float64) ModelPriceBreakdown {
	if multiplier <= 0 {
		multiplier = 1
	}
	in := quotaPerMillion(mp.PromptPrice)
	out := quotaPerMillion(mp.CompletionPrice)
	cr := quotaPerMillion(mp.CacheReadPrice)
	cw := quotaPerMillion(mp.CacheWritePrice)
	unit := mp.UnitPrice
	if multiplier != 1 {
		in = applyMultiplier(in, multiplier)
		out = applyMultiplier(out, multiplier)
		cr = applyMultiplier(cr, multiplier)
		cw = applyMultiplier(cw, multiplier)
		unit = applyMultiplier(unit, multiplier)
	}
	return ModelPriceBreakdown{
		InputPerMillion:      in,
		OutputPerMillion:     out,
		CacheReadPerMillion:  cr,
		CacheWritePerMillion: cw,
		UnitPrice:            unit,
		Multiplier:           multiplier,
	}
}

func BillingLabel(t int16) string {
	if t == 2 {
		return "按次计费"
	}
	return "按量计费"
}
