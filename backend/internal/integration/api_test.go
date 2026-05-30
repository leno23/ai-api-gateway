package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/leno23/ai-api-gateway/internal/config"
	applog "github.com/leno23/ai-api-gateway/internal/log"
	"github.com/leno23/ai-api-gateway/internal/router"
	"github.com/leno23/ai-api-gateway/internal/store"
)

func skipWithoutIntegrationDSN(t *testing.T) {
	t.Helper()
	if os.Getenv("INTEGRATION_DATABASE_DSN") == "" {
		t.Skip("set INTEGRATION_DATABASE_DSN (and optionally INTEGRATION_REDIS_ADDR) to run integration tests")
	}
}

type integrationEnv struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func newIntegrationEnv(t *testing.T) integrationEnv {
	t.Helper()
	skipWithoutIntegrationDSN(t)
	t.Setenv("DATABASE_DSN", os.Getenv("INTEGRATION_DATABASE_DSN"))
	if addr := os.Getenv("INTEGRATION_REDIS_ADDR"); addr != "" {
		t.Setenv("REDIS_ADDR", addr)
	}
	t.Setenv("JWT_SECRET", "integration-test-jwt-secret-min-32-chars!!")
	t.Setenv("RATE_LIMIT_GLOBAL_PER_MIN", "0")
	t.Setenv("RATE_LIMIT_USER_PER_MIN", "0")
	t.Setenv("QUOTA_RECONCILE_INTERVAL_SEC", "0")

	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	log, err := applog.New("development")
	if err != nil {
		t.Fatal(err)
	}
	db, err := store.OpenPostgres(cfg.DatabaseDSN, log)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	})
	rdb, err := store.OpenRedis(cfg.RedisAddr)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = rdb.Close() })

	gin.SetMode(gin.TestMode)
	return integrationEnv{Router: router.New(cfg, db, rdb, log), DB: db}
}

func newIntegrationRouter(t *testing.T) *gin.Engine {
	t.Helper()
	return newIntegrationEnv(t).Router
}

func TestIntegration_Health(t *testing.T) {
	r := newIntegrationRouter(t)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("health: %d %s", w.Code, w.Body.String())
	}
}

func TestIntegration_OpenAPIYAML(t *testing.T) {
	r := newIntegrationRouter(t)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("openapi: %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "openapi: 3.0.3") {
		t.Fatalf("unexpected openapi body")
	}
}

func TestIntegration_RegisterLoginModelsFlow(t *testing.T) {
	r := newIntegrationRouter(t)
	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	email := "i" + suffix + "@example.com"
	username := "u" + suffix
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

	loginBody := map[string]any{"email": email, "password": "testpass12"}
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

	keyBody := map[string]any{"name": "itest-key"}
	kb, _ := json.Marshal(keyBody)
	req := httptest.NewRequest(http.MethodPost, "/user/api-keys", bytes.NewReader(kb))
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("api-keys: %d %s", w.Code, w.Body.String())
	}
	var keyResp struct {
		APIKey string `json:"api_key"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &keyResp); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(keyResp.APIKey, "sk-") {
		t.Fatalf("bad api key: %q", keyResp.APIKey)
	}

	req = httptest.NewRequest(http.MethodGet, "/v1/models", nil)
	req.Header.Set("Authorization", "Bearer "+keyResp.APIKey)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("models: %d %s", w.Code, w.Body.String())
	}
}

func TestIntegration_AdminForbiddenForNormalUser(t *testing.T) {
	r := newIntegrationRouter(t)
	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	email := "n" + suffix + "@example.com"
	username := "nu" + suffix
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

	lb, _ := json.Marshal(map[string]any{"email": email, "password": "testpass12"})
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(lb)))
	var loginResp struct {
		AccessToken string `json:"access_token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &loginResp)

	req := httptest.NewRequest(http.MethodGet, "/admin/redeem/codes", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d %s", w.Code, w.Body.String())
	}
}
