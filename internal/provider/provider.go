// Package provider defines the provider interface and registry.
package provider

import (
	"context"
)

// Provider fetches a secret value by spec.
type Provider interface {
	FetchSecret(ctx context.Context, spec SecretSpec) (string, error)
}

// SecretSpec is an immutable specification for a secret fetch.
type SecretSpec struct {
	Alias    string            // Human-readable alias for the secret
	Name     string            // Provider-specific secret name/path
	Provider string            // Provider type (aws, gcp, azure, etc.)
	EnvName  string            // Environment variable name
	Extras   map[string]string // Provider-specific configuration options
}

// Global registry of available providers
var registry = map[string]Provider{}

// Register adds a provider implementation under a name.
func Register(name string, p Provider) {
	registry[name] = p
}

// Get returns a provider by name.
func Get(name string) (Provider, bool) {
	p, ok := registry[name]
	return p, ok
}

