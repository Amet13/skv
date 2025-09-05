package gcp

import (
	"context"
	"errors"
	"testing"

	secretspb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"skv/internal/provider"
)

func TestGCPNotFoundMapsToErrNotFound(t *testing.T) {
	old := gcpAccess
	defer func() { gcpAccess = old }()
	gcpAccess = func(_ context.Context, _ string, _ string) (*secretspb.AccessSecretVersionResponse, error) {
		return nil, status.Error(codes.NotFound, "nope")
	}
	p := &gcpProvider{}
	_, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "projects/p/secrets/s/versions/1"})
	if !errors.Is(err, provider.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGCPSuccess(t *testing.T) {
	old := gcpAccess
	defer func() { gcpAccess = old }()
	gcpAccess = func(_ context.Context, _ string, _ string) (*secretspb.AccessSecretVersionResponse, error) {
		return &secretspb.AccessSecretVersionResponse{Payload: &secretspb.SecretPayload{Data: []byte("ok")}}, nil
	}
	p := &gcpProvider{}
	out, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "projects/p/secrets/s/versions/1"})
	if err != nil || out != "ok" {
		t.Fatalf("got %q err=%v", out, err)
	}
}

func TestGCPBuildsNameFromExtras(t *testing.T) {
	old := gcpAccess
	defer func() { gcpAccess = old }()
	gcpAccess = func(_ context.Context, _ string, _ string) (*secretspb.AccessSecretVersionResponse, error) {
		return &secretspb.AccessSecretVersionResponse{Payload: &secretspb.SecretPayload{Data: []byte("ok")}}, nil
	}
	p := &gcpProvider{}
	out, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "api_key", Extras: map[string]string{"project": "p", "version": "latest"}})
	if err != nil || out != "ok" {
		t.Fatalf("got %q err=%v", out, err)
	}
}

func TestGCPPermissionDeniedMapsToNotFound(t *testing.T) {
	old := gcpAccess
	defer func() { gcpAccess = old }()
	gcpAccess = func(_ context.Context, _ string, _ string) (*secretspb.AccessSecretVersionResponse, error) {
		return nil, status.Error(codes.PermissionDenied, "nope")
	}
	p := &gcpProvider{}
	if _, err := p.FetchSecret(context.Background(), provider.SecretSpec{Alias: "a", Name: "projects/p/secrets/s/versions/1"}); !errors.Is(err, provider.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

