package adapter

import (
	"context"
	"net/http"
)

// ProviderAdapter abstracts an OpenAI-compatible upstream HTTP surface.
type ProviderAdapter interface {
	Name() string
	ChatCompletions(ctx context.Context, body []byte, stream bool) (*http.Response, error)
	Models(ctx context.Context) (*http.Response, error)
}
