package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"skv/internal/provider"
)

type mockProvErr struct{}

func (mockProvErr) FetchSecret(_ context.Context, _ provider.SecretSpec) (string, error) {
	return "", errors.New("boom")
}

func TestRunStrictPartialFailureReturnsError(t *testing.T) {
	cfg := []byte("secrets:\n- alias: ok\n  provider: ok\n  name: n\n  env: A\n- alias: bad\n  provider: bad\n  name: n\n  env: B\n")
	dir := t.TempDir()
	path := filepath.Join(dir, ".skv.yaml")
	if err := os.WriteFile(path, cfg, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKV_CONFIG", path)
	provider.Register("ok", mockProvOK{val: "x"})
	provider.Register("bad", mockProvErr{})

	root := newRootCmd()
	out := &bytes.Buffer{}
	errb := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errb)
	root.SetArgs([]string{"run", "-s", "ok", "-s", "bad", "--", "/usr/bin/env"})
	if err := root.Execute(); err == nil {
		t.Fatalf("expected strict mode failure")
	}
}

type okProv struct{}

func (okProv) FetchSecret(_ context.Context, _ provider.SecretSpec) (string, error) { return "x", nil }

type nfProv2 struct{}

func (nfProv2) FetchSecret(_ context.Context, _ provider.SecretSpec) (string, error) {
	return "", provider.ErrNotFound
}

func TestRunNonStrictPartialSuccess(t *testing.T) {
	cfg := []byte("secrets:\n- alias: a\n  provider: okprov\n  name: n\n  env: A\n- alias: b\n  provider: nfprov\n  name: n\n  env: B\n")
	dir := t.TempDir()
	path := filepath.Join(dir, ".skv.yaml")
	if err := os.WriteFile(path, cfg, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKV_CONFIG", path)
	provider.Register("okprov", okProv{})
	provider.Register("nfprov", nfProv2{})

	root := newRootCmd()
	out := &bytes.Buffer{}
	errb := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errb)
	root.SetArgs([]string{"run", "--strict=false", "-s", "a", "-s", "b", "--", "/bin/echo", "hi"})
	if err := root.Execute(); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

