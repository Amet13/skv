// Package provider defines the provider interface and registry.
package provider

import (
	"context"
	"fmt"
)

// Provider fetches a secret value by spec.
type Provider interface {
	FetchSecret(ctx context.Context, spec SecretSpec) (string, error)
}

// SecretSpec is an immutable specification for a secret fetch.
type SecretSpec struct {
	Alias    string
	Name     string
	Provider string
	EnvName  string
	Extras   map[string]string
}

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

// NotImplementedProvider returns a provider that always errors with a clear message.
type notImpl struct{ name string }

// NewNotImplementedProvider returns a stub provider that always fails.
func NewNotImplementedProvider(name string) Provider { return &notImpl{name: name} }

func (n *notImpl) FetchSecret(_ context.Context, _ SecretSpec) (string, error) {
	return "", fmt.Errorf("provider %s not implemented yet", n.name)
}

