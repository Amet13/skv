package gcp

import (
	"context"
	"testing"

	secretspb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"skv/internal/provider"
)

func TestGCPMissingProjectError(t *testing.T) {
	p := &gcpProvider{}
	if _, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "not-qualified", Extras: map[string]string{}}); err == nil {
		t.Fatalf("expected error when project missing and name not qualified")
	}
}

func TestGCPEmptyPayloadError(t *testing.T) {
	old := gcpAccess
	defer func() { gcpAccess = old }()
	gcpAccess = func(_ context.Context, _ string, _ string) (*secretspb.AccessSecretVersionResponse, error) {
		return &secretspb.AccessSecretVersionResponse{}, nil
	}
	p := &gcpProvider{}
	if _, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "projects/p/secrets/s/versions/1"}); err == nil {
		t.Fatalf("expected error for empty payload")
	}
}

