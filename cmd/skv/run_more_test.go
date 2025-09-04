package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"skv/internal/provider"
)

type mockProvOK struct{ val string }

func (m mockProvOK) FetchSecret(_ context.Context, _ provider.SecretSpec) (string, error) {
	return m.val, nil
}

func TestRunInvalidTimeoutReturnsCode2(t *testing.T) {
	cfg := []byte("secrets:\n- alias: a\n  provider: ok\n  name: n\n  env: A\n")
	dir := t.TempDir()
	path := filepath.Join(dir, ".skv.yaml")
	if err := os.WriteFile(path, cfg, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKV_CONFIG", path)
	provider.Register("ok", mockProvOK{val: "x"})

	root := newRootCmd()
	out := &bytes.Buffer{}
	errb := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errb)
	root.SetArgs([]string{"run", "--timeout", "bad", "-s", "a", "--", "/bin/echo", "hi"})
	if err := root.Execute(); err == nil {
		t.Fatalf("expected error for invalid timeout")
	}
}

func TestRunDryRunConcurrency(t *testing.T) {
	cfg := []byte("secrets:\n- alias: a\n  provider: ok\n  name: n\n  env: A\n- alias: b\n  provider: ok\n  name: n\n  env: B\n")
	dir := t.TempDir()
	path := filepath.Join(dir, ".skv.yaml")
	if err := os.WriteFile(path, cfg, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKV_CONFIG", path)
	provider.Register("ok", mockProvOK{val: "x"})

	root := newRootCmd()
	out := &bytes.Buffer{}
	errb := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(errb)
	root.SetArgs([]string{"run", "--dry-run", "--concurrency", "8", "-s", "a", "-s", "b", "--", "/bin/echo", "hi"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	s := errb.String()
	if s == "" || !bytes.Contains([]byte(s), []byte("A=")) || !bytes.Contains([]byte(s), []byte("B=")) {
		t.Fatalf("expected dry-run env additions, got %q", s)
	}
}

