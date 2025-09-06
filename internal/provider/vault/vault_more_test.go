package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"skv/internal/provider"
)

// Minimal integration-style test against a fake HTTP server exercising logical read fallback path.
func TestVaultLogicalReadFallback(t *testing.T) {
	// Fake Vault logical read returning data without KVv2 nesting
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]interface{}{"value": "ok"}})
	}))
	defer srv.Close()

	p := &vaultProvider{}
	spec := provider.SecretSpec{
		Alias:  "a",
		Name:   srv.URL, // direct URL causes client.Logical().ReadWithContext to call this endpoint
		Extras: map[string]string{"address": srv.URL},
	}
	out, err := p.FetchSecret(context.Background(), spec)
	if err != nil || out == "" {
		t.Fatalf("unexpected err=%v out=%q", err, out)
	}
}

func TestVaultKV2EmptyPayloadReturnsJSON(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/kv/data/empty", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"data":     map[string]any{},
				"metadata": map[string]any{"version": 1},
			},
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	p := &vaultProvider{}
	spec := provider.SecretSpec{
		Alias:  "e",
		Name:   "kv/data/empty",
		Extras: map[string]string{"address": srv.URL},
	}
	out, err := p.FetchSecret(context.Background(), spec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "{}" {
		t.Fatalf("expected empty JSON object, got %q", out)
	}
}

