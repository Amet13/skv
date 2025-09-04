package provider

import (
	"context"
	"testing"
)

type testProv struct{ out string }

func (t testProv) FetchSecret(_ context.Context, _ SecretSpec) (string, error) { return t.out, nil }

func TestRegistryRegisterAndGet(t *testing.T) {
	Register("_test", testProv{out: "ok"})
	p, ok := Get("_test")
	if !ok {
		t.Fatalf("expected provider to be registered")
	}
	got, err := p.FetchSecret(context.Background(), SecretSpec{Name: "n"})
	if err != nil || got != "ok" {
		t.Fatalf("got %q, err=%v", got, err)
	}
}

func TestRegistryGetUnknown(t *testing.T) {
	if _, ok := Get("__does_not_exist__"); ok {
		t.Fatalf("expected unknown provider to not be found")
	}
}

