package openapi

import (
	"strings"
	"testing"
)

// portalPaths are tenant-portal REST routes that must stay documented in spec.yaml.
var portalPaths = []string{
	"/api/announcements",
	"/api/dashboard/stats",
	"/api/dashboard/charts",
	"/api/dashboard/nodes",
	"/api/nodes/ping",
	"/api/logs/usage",
	"/api/logs/tasks",
	"/api/wallet/summary",
	"/api/wallet/redeem",
	"/api/wallet/affiliate/transfer",
	"/api/wallet/recharge/mock",
	"/api/playground/chat",
	"/api/tokens",
	"/api/tokens/batch-delete",
	"/api/user/self",
	"/api/models",
	"/api/token-groups",
}

func TestSpecDocumentsPortalRoutes(t *testing.T) {
	body := string(Spec)
	for _, path := range portalPaths {
		needle := "\n  " + path + ":"
		if !strings.Contains(body, needle) {
			t.Errorf("spec.yaml missing path definition: %s", path)
		}
	}
}

func TestSpecHasPortalTag(t *testing.T) {
	if !strings.Contains(string(Spec), "name: Portal") {
		t.Fatal("spec.yaml missing Portal tag")
	}
}
