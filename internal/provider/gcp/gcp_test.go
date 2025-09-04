package gcp

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGCPStatusNotFound(t *testing.T) {
	st := status.Error(codes.NotFound, "not found")
	if s, ok := status.FromError(st); !ok || s.Code() != codes.NotFound {
		t.Fatalf("expected NotFound status")
	}
}

