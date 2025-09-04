package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListPrintsAliases(t *testing.T) {
	cfg := []byte("secrets:\n- alias: a\n  provider: mock\n  name: n\n- alias: b\n  provider: mock\n  name: n\n")
	dir := t.TempDir()
	path := filepath.Join(dir, ".skv.yaml")
	if err := os.WriteFile(path, cfg, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKV_CONFIG", path)

	root := newRootCmd()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(out)
	root.SetArgs([]string{"list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	s := out.String()
	if !strings.Contains(s, "a") || !strings.Contains(s, "b") {
		t.Fatalf("expected aliases a and b in output: %q", s)
	}
}

