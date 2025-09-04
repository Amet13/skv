package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"skv/internal/provider"
)

type errProvider struct{}

func (e errProvider) FetchSecret(_ context.Context, _ provider.SecretSpec) (string, error) {
	return "", errors.New("boom")
}

func TestRunStrictMissingAliasFails(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skip on windows")
	}
	// empty config
	dir := t.TempDir()
	path := filepath.Join(dir, ".skv.yaml")
	if err := os.WriteFile(path, []byte("secrets: []\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKV_CONFIG", path)

	root := newRootCmd()
	out := &bytes.Buffer{}
	errb := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errb)
	root.SetArgs([]string{"run", "-s", "nope", "--", "/bin/echo", "hi"})
	if err := root.Execute(); err == nil {
		t.Fatalf("expected error for missing alias with strict=true")
	}
}

func TestRunNonStrictProviderErrorContinues(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skip on windows")
	}
	cfg := []byte("secrets:\n- alias: a\n  provider: errprov\n  name: n\n  env: A\n")
	dir := t.TempDir()
	path := filepath.Join(dir, ".skv.yaml")
	if err := os.WriteFile(path, cfg, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKV_CONFIG", path)
	provider.Register("errprov", errProvider{})

	root := newRootCmd()
	out := &bytes.Buffer{}
	errb := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errb)
	root.SetArgs([]string{"run", "--strict=false", "-s", "a", "--", "/bin/echo", "hi"})
	if err := root.Execute(); err != nil {
		t.Fatalf("expected success in non-strict mode, got %v", err)
	}
}

