package adapter

import (
	"context"
	"net/http"

	"github.com/leno23/ai-api-gateway/internal/relay"
)

// OpenAICompat wraps relay.Client for OpenAI-compatible HTTPS upstreams.
type OpenAICompat struct {
	client *relay.Client
}

func NewOpenAICompat(baseURL, apiKey string) *OpenAICompat {
	return &OpenAICompat{client: relay.New(baseURL, apiKey)}
}

func (o *OpenAICompat) Name() string { return "openai-compat" }

func (o *OpenAICompat) ChatCompletions(ctx context.Context, body []byte, stream bool) (*http.Response, error) {
	return o.client.ChatCompletions(ctx, body, stream)
}

func (o *OpenAICompat) Models(ctx context.Context) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, o.client.BaseURL+"/models", nil)
	if err != nil {
		return nil, err
	}
	if o.client.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+o.client.APIKey)
	}
	return o.client.HTTP.Do(req)
}
