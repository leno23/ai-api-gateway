package openapi

import (
	"strings"
	"testing"
)

func TestSpecEmbeddedNonEmpty(t *testing.T) {
	if len(Spec) < 200 {
		t.Fatalf("spec too short: %d", len(Spec))
	}
	if !strings.Contains(string(Spec), "openapi: 3.0.3") {
		t.Fatalf("missing openapi version header")
	}
}
