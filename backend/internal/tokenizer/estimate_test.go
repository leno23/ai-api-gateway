package tokenizer

import (
	"encoding/json"
	"testing"
)

func TestEstimatePromptTokens(t *testing.T) {
	raw := json.RawMessage(`[{"role":"user","content":"hello"}]`)
	n := EstimatePromptTokens(raw)
	if n < 1 {
		t.Fatalf("expected >= 1, got %d", n)
	}
	if EstimatePromptTokens(nil) != 1 {
		t.Fatalf("empty messages should be 1")
	}
}
