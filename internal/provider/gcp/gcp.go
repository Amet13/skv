// Package gcp implements Google Secret Manager provider.
package gcp

import (
	"context"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretspb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"skv/internal/provider"
)

type gcpProvider struct{}

// New returns a new GCP Secret Manager provider.
func New() provider.Provider { return &gcpProvider{} }

// seam for testing access
var gcpAccess = func(ctx context.Context, name string, credsFile string) (*secretspb.AccessSecretVersionResponse, error) {
	var client *secretmanager.Client
	var err error
	if strings.TrimSpace(credsFile) != "" {
		client, err = secretmanager.NewClient(ctx, option.WithCredentialsFile(credsFile))
	} else {
		client, err = secretmanager.NewClient(ctx)
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = client.Close() }()
	return client.AccessSecretVersion(ctx, &secretspb.AccessSecretVersionRequest{Name: name})
}

func (g *gcpProvider) FetchSecret(ctx context.Context, spec provider.SecretSpec) (string, error) {
	name := spec.Name
	if !strings.HasPrefix(name, "projects/") {
		project := spec.Extras["project"]
		if project == "" {
			return "", fmt.Errorf("gcp: missing metadata.project for %s", spec.Alias)
		}
		version := spec.Extras["version"]
		if version == "" {
			version = "latest"
		}
		name = fmt.Sprintf("projects/%s/secrets/%s/versions/%s", project, name, version)
	}

	credsFile := strings.TrimSpace(spec.Extras["credentials_file"])
	res, err := gcpAccess(ctx, name, credsFile)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && (st.Code() == codes.NotFound || st.Code() == codes.PermissionDenied) {
			return "", provider.ErrNotFound
		}
		return "", fmt.Errorf("gcp: access secret: %w", err)
	}
	if res.Payload == nil || res.Payload.Data == nil {
		return "", fmt.Errorf("gcp: empty payload for %s", name)
	}
	return string(res.Payload.Data), nil
}
