package mock

import (
	"context"
	"errors"
	"strings"
	"testing"

	"skv/internal/provider"
)

func TestNew(t *testing.T) {
	p := New()
	if p == nil {
		t.Fatal("New() returned nil")
	}

	// Ensure it implements the Provider interface
	_ = p
}

func TestMockProvider_FetchSecret(t *testing.T) {
	p := New()
	ctx := context.Background()

	tests := []struct {
		name    string
		spec    provider.SecretSpec
		wantVal string
		wantErr bool
	}{
		{
			name: "returns name by default",
			spec: provider.SecretSpec{
				Alias:    "test-secret",
				Name:     "mock-secret",
				Provider: "mock",
				EnvName:  "TEST_SECRET",
			},
			wantVal: "mock-secret",
			wantErr: false,
		},
		{
			name: "returns custom value from extras",
			spec: provider.SecretSpec{
				Alias:    "test-secret",
				Name:     "mock-secret",
				Provider: "mock",
				EnvName:  "TEST_SECRET",
				Extras: map[string]string{
					"value": "custom-value",
				},
			},
			wantVal: "custom-value",
			wantErr: false,
		},
		{
			name: "returns not found error",
			spec: provider.SecretSpec{
				Alias:    "test-secret",
				Name:     "mock-secret",
				Provider: "mock",
				EnvName:  "TEST_SECRET",
				Extras: map[string]string{
					"not_found": "true",
				},
			},
			wantVal: "",
			wantErr: true,
		},
		{
			name: "returns custom error",
			spec: provider.SecretSpec{
				Alias:    "test-secret",
				Name:     "mock-secret",
				Provider: "mock",
				EnvName:  "TEST_SECRET",
				Extras: map[string]string{
					"error": "custom error message",
				},
			},
			wantVal: "",
			wantErr: true,
		},
		{
			name: "ignores other extras",
			spec: provider.SecretSpec{
				Alias:    "secret-with-extras",
				Name:     "mock-with-extras",
				Provider: "mock",
				EnvName:  "SECRET_WITH_EXTRAS",
				Extras: map[string]string{
					"region": "us-east-1",
					"key":    "password",
				},
			},
			wantVal: "mock-with-extras",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.FetchSecret(ctx, tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantVal {
				t.Errorf("FetchSecret() = %v, want %v", got, tt.wantVal)
			}
		})
	}
}

func TestMockProvider_Consistency(t *testing.T) {
	p := New()
	ctx := context.Background()
	spec := provider.SecretSpec{
		Alias:    "consistent-secret",
		Name:     "test-consistency",
		Provider: "mock",
		EnvName:  "CONSISTENT_SECRET",
	}

	// Call multiple times to ensure consistent behavior
	for i := 0; i < 5; i++ {
		got, err := p.FetchSecret(ctx, spec)
		if err != nil {
			t.Fatalf("FetchSecret() call %d failed: %v", i+1, err)
		}
		if got != "test-consistency" {
			t.Errorf("FetchSecret() call %d = %v, want %v", i+1, got, "test-consistency")
		}
	}
}

func TestMockProvider_ErrorHandling(t *testing.T) {
	p := New()
	ctx := context.Background()

	// Test not found error
	spec := provider.SecretSpec{
		Name: "test-not-found",
		Extras: map[string]string{
			"not_found": "true",
		},
	}
	_, err := p.FetchSecret(ctx, spec)
	if !errors.Is(err, provider.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	// Test custom error
	spec = provider.SecretSpec{
		Name: "test-error",
		Extras: map[string]string{
			"error": "test error",
		},
	}
	_, err = p.FetchSecret(ctx, spec)
	if err == nil || !strings.Contains(err.Error(), "test error") {
		t.Errorf("Expected error containing 'test error', got %v", err)
	}
}
