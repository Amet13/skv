package main

import (
	"bytes"
	"context"
	"os"
	"testing"

	"skv/internal/provider"
)

type mockProvider struct{ val string }

func (m mockProvider) FetchSecret(_ context.Context, _ provider.SecretSpec) (string, error) {
	return m.val, nil
}

func TestGetCommandPrintsValue(t *testing.T) {
	// Prepare a temporary config file
	cfg := []byte("secrets:\n- alias: a\n  provider: mock\n  name: n\n")
	dir := t.TempDir()
	path := dir + "/.skv.yaml"
	if err := os.WriteFile(path, cfg, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("SKV_CONFIG", path); err != nil {
		t.Fatal(err)
	}

	provider.Register("mock", mockProvider{val: "VALUE"})

	root := newRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"get", "a", "--newline"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	got := buf.String()
	if got != "VALUE\n" {
		t.Fatalf("got %q, want %q", got, "VALUE\n")
	}
}

