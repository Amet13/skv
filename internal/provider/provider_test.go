package provider

import (
	"context"
	"testing"
)

// Mock provider for testing
type testProvider struct{}

func (t *testProvider) FetchSecret(_ context.Context, _ SecretSpec) (string, error) {
	return "test_value", nil
}

func TestRegister(t *testing.T) {
	// Test registering a provider
	mockProv := &testProvider{}
	Register("test_provider", mockProv)

	// Verify it can be retrieved
	retrieved, ok := Get("test_provider")
	if !ok {
		t.Error("Expected provider to be registered")
	}
	if retrieved != mockProv {
		t.Error("Retrieved provider doesn't match registered provider")
	}
}

func TestGet(t *testing.T) {
	// Register test providers
	Register("aws", &testProvider{})
	Register("mock", &testProvider{})

	tests := []struct {
		name       string
		provider   string
		shouldFind bool
	}{
		{
			name:       "existing provider",
			provider:   "aws",
			shouldFind: true,
		},
		{
			name:       "non-existing provider",
			provider:   "nonexistent",
			shouldFind: false,
		},
		{
			name:       "mock provider",
			provider:   "mock",
			shouldFind: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := Get(tt.provider)
			if ok != tt.shouldFind {
				t.Errorf("Get(%q) found = %v, expected %v", tt.provider, ok, tt.shouldFind)
			}
		})
	}
}

func BenchmarkProviderRegistry(b *testing.B) {
	// Benchmark provider registry operations
	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = Get("aws")
		}
	})
}

func BenchmarkSecretSpecCreation(b *testing.B) {
	b.Run("NewSecretSpec", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			spec := SecretSpec{
				Alias:    "test",
				Name:     "test-secret",
				Provider: "aws",
				EnvName:  "TEST_SECRET",
				Extras: map[string]string{
					"region": "us-east-1",
				},
			}
			_ = spec
		}
	})
}

