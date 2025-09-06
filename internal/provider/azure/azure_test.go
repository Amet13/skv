package azure

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azsecrets "github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
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

func TestAzureNotFoundMapsToErrNotFound(t *testing.T) {
	old := azureGetSecret
	defer func() { azureGetSecret = old }()
	azureGetSecret = func(_ context.Context, _ string, _ string, _ string) (*azsecrets.GetSecretResponse, error) {
		return nil, &azcore.ResponseError{StatusCode: 404}
	}
	p := &azureProvider{}
	_, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "n", Extras: map[string]string{"vault_url": "https://v"}})
	if !errors.Is(err, provider.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestAzureSuccess(t *testing.T) {
	old := azureGetSecret
	defer func() { azureGetSecret = old }()
	val := "ok"
	azureGetSecret = func(_ context.Context, _ string, _ string, _ string) (*azsecrets.GetSecretResponse, error) {
		resp := azsecrets.GetSecretResponse{}
		resp.Value = &val
		return &resp, nil
	}
	p := &azureProvider{}
	out, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "n", Extras: map[string]string{"vault_url": "https://v"}})
	if err != nil || out != "ok" {
		t.Fatalf("got %q err=%v", out, err)
	}
}

func TestAzureMissingVaultURL(t *testing.T) {
	p := &azureProvider{}
	if _, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "n", Extras: map[string]string{}}); err == nil {
		t.Fatalf("expected error for missing vault_url")
	}
}

func TestAzureUnknownError(t *testing.T) {
	old := azureGetSecret
	defer func() { azureGetSecret = old }()
	azureGetSecret = func(_ context.Context, _ string, _ string, _ string) (*azsecrets.GetSecretResponse, error) {
		return nil, errors.New("boom")
	}
	p := &azureProvider{}
	if _, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "n", Extras: map[string]string{"vault_url": "https://v"}}); err == nil {
		t.Fatalf("expected error")
	}
}

func TestAzureEmptyValue(t *testing.T) {
	old := azureGetSecret
	defer func() { azureGetSecret = old }()
	azureGetSecret = func(_ context.Context, _ string, _ string, _ string) (*azsecrets.GetSecretResponse, error) {
		r := azsecrets.GetSecretResponse{}
		return &r, nil
	}
	p := &azureProvider{}
	if _, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "n", Extras: map[string]string{"vault_url": "https://v"}}); err == nil {
		t.Fatalf("expected error for empty value")
	}
}

