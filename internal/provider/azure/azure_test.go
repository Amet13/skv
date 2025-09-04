package azure

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azsecrets "github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"skv/internal/provider"
)

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

