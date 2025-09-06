// Package mock implements a simple testing provider.
// Not registered in production; used by tests to simulate providers.
package mock

import (
	"context"
	"fmt"
	"strings"

	"skv/internal/provider"
)

type mockProvider struct{}

// New returns a mock provider for testing.
func New() provider.Provider { return &mockProvider{} }

// Behavior controlled by extras:
// - value: returned as secret value
// - not_found: "true" to return provider.ErrNotFound
// - error: non-empty string to return an error with that message
// If none provided, returns spec.Name as value.
func (m *mockProvider) FetchSecret(_ context.Context, spec provider.SecretSpec) (string, error) {
	if strings.EqualFold(spec.Extras["not_found"], "true") {
		return "", provider.ErrNotFound
	}
	if msg := strings.TrimSpace(spec.Extras["error"]); msg != "" {
		return "", fmt.Errorf("mock: %s", msg)
	}
	if v := strings.TrimSpace(spec.Extras["value"]); v != "" {
		return v, nil
	}
	return spec.Name, nil
}

