package relay

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	HTTP    *http.Client
	BaseURL string
	APIKey  string
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		HTTP: &http.Client{
			Timeout: 120 * time.Second,
		},
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  apiKey,
	}
}

// ChatCompletions posts the raw JSON body to upstream OpenAI-compatible path.
func (c *Client) ChatCompletions(ctx context.Context, body []byte, stream bool) (*http.Response, error) {
	url := c.BaseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}
	if stream {
		req.Header.Set("Accept", "text/event-stream")
	}
	return c.HTTP.Do(req)
}

// DrainClose reads and discards remaining body then closes.
func DrainClose(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}
