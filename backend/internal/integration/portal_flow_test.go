package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/leno23/ai-api-gateway/internal/model"
)

// TestIntegration_PortalCoreFlow covers register → login → portal token → playground chat.
func TestIntegration_PortalCoreFlow(t *testing.T) {
	env := newIntegrationEnv(t)
	r := env.Router

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	email := "p" + suffix + "@example.com"
	username := "pu" + suffix
	if len(username) > 64 {
		username = username[:64]
	}

	regBody := map[string]any{
		"email": email, "username": username, "password": "testpass12",
	}
	b, _ := json.Marshal(regBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(b)))
	if w.Code != http.StatusCreated {
		t.Fatalf("register: %d %s", w.Code, w.Body.String())
	}

	loginBody := map[string]any{"username": username, "password": "testpass12"}
	lb, _ := json.Marshal(loginBody)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(lb)))
	if w.Code != http.StatusOK {
		t.Fatalf("login: %d %s", w.Code, w.Body.String())
	}
	var loginResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil {
		t.Fatal(err)
	}
	if loginResp.AccessToken == "" {
		t.Fatal("empty access_token")
	}
	jwt := loginResp.AccessToken

	if err := env.DB.Model(&model.User{}).Where("email = ?", email).
		Update("quota", int64(50_000_000)).Error; err != nil {
		t.Fatalf("seed quota: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/user/self", nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("user self: %d %s", w.Code, w.Body.String())
	}

	tokenBody := map[string]any{"name": "portal-e2e", "token_group": "default"}
	tb, _ := json.Marshal(tokenBody)
	req = httptest.NewRequest(http.MethodPost, "/api/tokens", bytes.NewReader(tb))
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create token: %d %s", w.Code, w.Body.String())
	}
	var tokenResp struct {
		Success bool `json:"success"`
		Data    struct {
			Token struct {
				ID int64 `json:"id"`
			} `json:"token"`
			APIKey string `json:"api_key"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &tokenResp); err != nil {
		t.Fatal(err)
	}
	if tokenResp.Data.APIKey == "" || tokenResp.Data.Token.ID == 0 {
		t.Fatalf("unexpected token response: %s", w.Body.String())
	}

	playBody := map[string]any{
		"model":       "gpt-4o-mini",
		"stream":      true,
		"max_tokens":  16,
		"api_key_id":  tokenResp.Data.Token.ID,
		"token_group": "default",
		"messages": []map[string]string{
			{"role": "user", "content": "Say hi in one word."},
		},
	}
	pb, _ := json.Marshal(playBody)
	req = httptest.NewRequest(http.MethodPost, "/api/playground/chat", bytes.NewReader(pb))
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	switch w.Code {
	case http.StatusOK:
		ct := w.Header().Get("Content-Type")
		if ct == "" {
			t.Fatalf("playground: missing content-type")
		}
	case http.StatusPaymentRequired, http.StatusBadGateway, http.StatusServiceUnavailable:
		t.Logf("playground upstream/quota path: %d %s", w.Code, w.Body.String())
	default:
		t.Fatalf("playground: %d %s", w.Code, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/dashboard/stats", nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("dashboard stats: %d %s", w.Code, w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/api/nodes/ping", nil)
	req.Header.Set("Authorization", "Bearer "+jwt)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("nodes ping: %d %s", w.Code, w.Body.String())
	}
}
