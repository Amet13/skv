package azure

import (
	"context"
	"errors"
	"testing"

	"skv/internal/provider"
)

func TestAppConfigSuccess(t *testing.T) {
	old := appcfgGet
	defer func() { appcfgGet = old }()
	appcfgGet = func(_ context.Context, _ string, _ string, _ string) (string, error) { return "ok", nil }
	p := NewAppConfig()
	out, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "key", Extras: map[string]string{"endpoint": "https://e.azconfig.io"}})
	if err != nil || out != "ok" {
		t.Fatalf("got %q err=%v", out, err)
	}
}

func TestAppConfigMissingEndpoint(t *testing.T) {
	p := NewAppConfig()
	if _, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "k", Extras: map[string]string{}}); err == nil {
		t.Fatalf("expected error for missing endpoint")
	}
}

func TestAppConfigError(t *testing.T) {
	old := appcfgGet
	defer func() { appcfgGet = old }()
	appcfgGet = func(_ context.Context, _ string, _ string, _ string) (string, error) { return "", errors.New("boom") }
	p := NewAppConfig()
	if _, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "k", Extras: map[string]string{"endpoint": "https://e.azconfig.io"}}); err == nil {
		t.Fatalf("expected error")
	}
}

