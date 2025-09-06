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

type nfProv struct{}

func (nfProv) FetchSecret(_ context.Context, _ provider.SecretSpec) (string, error) {
	return "", provider.ErrNotFound
}

func TestGetNotFoundExit4(t *testing.T) {
	cfg := []byte("secrets:\n- alias: a\n  provider: nf\n  name: n\n")
	dir := t.TempDir()
	path := filepath.Join(dir, ".skv.yaml")
	if err := os.WriteFile(path, cfg, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKV_CONFIG", path)
	provider.Register("nf", nfProv{})

	root := newRootCmd()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"get", "a"})
	err := root.Execute()
	var ee exitCodeError
	if err == nil || !errors.As(err, &ee) || ee.code != 4 {
		t.Fatalf("expected exit code 4, got %v", err)
	}
}
