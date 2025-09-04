package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"skv/internal/provider"
)

type mockProvider2 struct{ val string }

func (m mockProvider2) FetchSecret(_ context.Context, _ provider.SecretSpec) (string, error) {
	return m.val, nil
}

func TestRunDryRunMasksValues(t *testing.T) {
	cfg := []byte("secrets:\n- alias: token\n  provider: mock2\n  name: n\n  env: APP_TOKEN\n")
	dir := t.TempDir()
	path := filepath.Join(dir, ".skv.yaml")
	if err := os.WriteFile(path, cfg, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKV_CONFIG", path)

	provider.Register("mock2", mockProvider2{val: "supersecretvalue"})

	root := newRootCmd()
	out := &bytes.Buffer{}
	errb := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errb)
	root.SetArgs([]string{"run", "--dry-run", "--secret", "token", "--", "echo", "hi"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if got := errb.String(); got == "" {
		t.Fatalf("expected dry-run output, got empty")
	}
}

