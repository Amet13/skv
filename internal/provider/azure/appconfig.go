// Package azure provides providers for Azure services.
// This file implements Azure App Configuration (key-value parameter store).
package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig"
	"skv/internal/provider"
)

type appConfigProvider struct{}

// NewAppConfig returns a provider for Azure App Configuration.
func NewAppConfig() provider.Provider { return &appConfigProvider{} }

// seam for testing
var appcfgGet = func(ctx context.Context, endpoint, key, label string) (string, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return "", fmt.Errorf("azure appconfig credential: %w", err)
	}
	client, err := azappconfig.NewClient(endpoint, cred, nil)
	if err != nil {
		return "", fmt.Errorf("azure appconfig client: %w", err)
	}
	opts := &azappconfig.GetSettingOptions{}
	if strings.TrimSpace(label) != "" {
		opts.Label = &label
	}
	resp, err := client.GetSetting(ctx, key, opts)
	if err != nil {
		return "", err
	}
	if resp.Value == nil {
		return "", fmt.Errorf("azure appconfig: empty value for %s", key)
	}
	return *resp.Value, nil
}

func (a *appConfigProvider) FetchSecret(ctx context.Context, spec provider.SecretSpec) (string, error) {
	endpoint := strings.TrimSpace(spec.Extras["endpoint"]) // e.g., https://<store>.azconfig.io
	if endpoint == "" {
		return "", fmt.Errorf("azure appconfig: missing extras.endpoint for %s", spec.Alias)
	}
	key := spec.Name
	label := spec.Extras["label"]
	val, err := appcfgGet(ctx, endpoint, key, label)
	if err != nil {
		// No clean not-found mapping in SDK; surface error
		return "", fmt.Errorf("azure appconfig get: %w", err)
	}
	return val, nil
}

