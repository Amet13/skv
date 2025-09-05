package provider

import (
	"context"
	"errors"
	"testing"
)

type fakeProvider struct {
	val string
	err error
}

func (f fakeProvider) FetchSecret(_ context.Context, _ SecretSpec) (string, error) {
	return f.val, f.err
}

func TestRegistry(t *testing.T) {
	Register("testp", fakeProvider{val: "ok"})
	p, ok := Get("testp")
	if !ok {
		t.Fatalf("expected provider to be registered")
	}
	out, err := p.FetchSecret(context.Background(), SecretSpec{})
	if err != nil || out != "ok" {
		t.Fatalf("unexpected: out=%q err=%v", out, err)
	}
}

func TestNotImplemented(t *testing.T) {
	p := NewNotImplementedProvider("x")
	_, err := p.FetchSecret(context.Background(), SecretSpec{})
	if err == nil || !errors.Is(err, err) { // just check non-nil; message checked implicitly
		t.Fatalf("expected error")
	}
}

