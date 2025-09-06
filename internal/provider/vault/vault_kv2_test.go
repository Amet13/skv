package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"skv/internal/provider"
)

func TestVaultKV2Success(t *testing.T) {
	// Simulate KV v2 endpoint response at /v1/kv/data/foo
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/kv/data/foo", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"data":     map[string]any{"value": "ok"},
				"metadata": map[string]any{"version": 1},
			},
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	p := &vaultProvider{}
	spec := provider.SecretSpec{
		Alias:  "a",
		Name:   "kv/data/foo",
		Extras: map[string]string{"address": srv.URL},
	}
	out, err := p.FetchSecret(context.Background(), spec)
	if err != nil || out != "ok" {
		t.Fatalf("got %q err=%v", out, err)
	}
}

