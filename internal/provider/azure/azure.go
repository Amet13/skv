// Package azure implements Azure Key Vault provider.
package azure

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	azsecrets "github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"skv/internal/provider"
)

type azureProvider struct{}

// New returns a new Azure Key Vault provider.
func New() provider.Provider { return &azureProvider{} }

func (a *azureProvider) FetchSecret(ctx context.Context, spec provider.SecretSpec) (string, error) {
	vaultURL := spec.Extras["vault_url"]
	if vaultURL == "" {
		return "", fmt.Errorf("azure: missing metadata.vault_url for %s", spec.Alias)
	}
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return "", fmt.Errorf("azure: credential: %w", err)
	}
	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		return "", fmt.Errorf("azure: client: %w", err)
	}

	version := spec.Extras["version"]
	resp, err := client.GetSecret(ctx, spec.Name, version, nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) {
			if respErr.StatusCode == 404 || respErr.StatusCode == 403 {
				return "", provider.ErrNotFound
			}
		}
		return "", fmt.Errorf("azure: get secret: %w", err)
	}
	if resp.Value == nil {
		return "", fmt.Errorf("azure: empty value for %s", spec.Name)
	}
	return *resp.Value, nil
}

