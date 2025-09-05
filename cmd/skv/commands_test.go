package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".skv.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func TestListFormats(t *testing.T) {
	// Ensure providers are registered
	_ = newRootCmd()
	cfg := "secrets:\n  - alias: a\n    provider: exec\n    name: /bin/echo\n    env: A\n    extras:\n      args: one\n  - alias: b\n    provider: exec\n    name: /bin/echo\n    env: B\n    extras:\n      args: two\n"
	cfgPath = writeTempConfig(t, cfg)

	// text
	cmd := newListCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("list text: %v", err)
	}
	s := out.String()
	if !strings.Contains(s, "a") || !strings.Contains(s, "b") {
		t.Fatalf("unexpected output: %q", s)
	}

	// json
	out.Reset()
	cmd = newListCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--format", "json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("list json: %v", err)
	}
	if !strings.Contains(out.String(), "\"Alias\": \"a\"") {
		t.Fatalf("unexpected json: %q", out.String())
	}
}

func TestExportEnv(t *testing.T) {
	_ = newRootCmd()
	cfg := "secrets:\n  - alias: token\n    provider: exec\n    name: /bin/echo\n    env: TOKEN\n    extras:\n      args: secret\n      trim: \"true\"\n"
	cfgPath = writeTempConfig(t, cfg)
	cmd := newExportCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--all", "--format", "env"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("export env: %v", err)
	}
	if strings.TrimSpace(out.String()) != "TOKEN=secret" {
		t.Fatalf("unexpected export: %q", out.String())
	}
}

func TestRunDryRun(t *testing.T) {
	_ = newRootCmd()
	cfg := "secrets:\n  - alias: token\n    provider: exec\n    name: /bin/echo\n    env: TOKEN\n    extras:\n      args: secret\n      trim: \"true\"\n"
	cfgPath = writeTempConfig(t, cfg)
	cmd := newRunCmd()
	var errBuf bytes.Buffer
	cmd.SetErr(&errBuf)
	cmd.SetArgs([]string{"--all", "--dry-run", "--", "/bin/echo", "hi"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("run dry-run: %v", err)
	}
	s := errBuf.String()
	if !strings.Contains(s, "[dry-run]") || !strings.Contains(s, "TOKEN=") {
		t.Fatalf("unexpected dry-run output: %q", s)
	}
}

