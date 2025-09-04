package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"skv/internal/provider"
)

type mockProvider3 struct{ val string }

func (m mockProvider3) FetchSecret(_ context.Context, _ provider.SecretSpec) (string, error) {
	return m.val, nil
}

func TestExportEnvFileAndExportLines(t *testing.T) {
	cfg := []byte("secrets:\n- alias: a\n  provider: mock3\n  name: n\n  env: A\n- alias: b\n  provider: mock3\n  name: n\n  env: B\n")
	dir := t.TempDir()
	path := filepath.Join(dir, ".skv.yaml")
	if err := os.WriteFile(path, cfg, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKV_CONFIG", path)
	provider.Register("mock3", mockProvider3{val: "x"})

	root := newRootCmd()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"export", "--all", "--env-file"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	s := out.String()
	if !strings.Contains(s, "A=x") || !strings.Contains(s, "B=x") {
		t.Fatalf("missing env vars: %q", s)
	}

	// Recreate root to reset flag defaults (env-file=false)
	root = newRootCmd()
	out = &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"export", "--all"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	s = out.String()
	if !strings.Contains(s, "export A=\"x\"") || !strings.Contains(s, "export B=\"x\"") {
		t.Fatalf("missing export lines: %q", s)
	}
}

