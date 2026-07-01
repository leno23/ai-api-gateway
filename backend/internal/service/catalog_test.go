package service

import (
	"testing"

	"github.com/leno23/ai-api-gateway/internal/model"
)

func TestModelPriceBreakdownFrom_multiplier(t *testing.T) {
	mp := model.ModelPrice{
		PromptPrice:     9000,
		CompletionPrice: 27000,
		CacheReadPrice:  4500,
	}
	base := ModelPriceBreakdownFrom(mp, 1)
	if base.InputPerMillion != 9_000_000 {
		t.Fatalf("input base: got %d", base.InputPerMillion)
	}
	applied := ModelPriceBreakdownFrom(mp, 1.5)
	if applied.InputPerMillion != 13_500_000 {
		t.Fatalf("input 1.5x: got %d", applied.InputPerMillion)
	}
}
