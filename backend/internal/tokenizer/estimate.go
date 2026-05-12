package tokenizer

import "encoding/json"

// EstimatePromptTokens returns a cheap heuristic from raw `messages` JSON (bytes length / 4).
// Replace with a real tokenizer (e.g. tiktoken) when 6.5 is fully implemented.
func EstimatePromptTokens(messages json.RawMessage) int64 {
	n := len(messages)
	if n == 0 {
		return 1
	}
	return int64(n)/4 + 1
}
