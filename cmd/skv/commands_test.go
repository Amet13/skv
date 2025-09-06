package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestListFormats(t *testing.T) {
	// Ensure providers are registered
	_ = newRootCmd()

	cfg := `secrets:
  - alias: a
    provider: exec
    name: /bin/echo
    env: A
    extras:
      args: one
  - alias: b
    provider: exec
    name: /bin/echo
    env: B
    extras:
      args: two`

	withTestConfig(t, cfg, func(_ string) {
		// Test text format
		cmd := newListCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetArgs([]string{})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("list text: %v", err)
		}

		output := out.String()
		assertStringContains(t, output, []string{"a", "b"})

		// Test JSON format
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
	})
}

func TestExportEnv(t *testing.T) {
	_ = newRootCmd()

	cfg := `secrets:
  - alias: token
    provider: exec
    name: /bin/echo
    env: TOKEN
    extras:
      args: secret
      trim: "true"`

	withTestConfig(t, cfg, func(_ string) {
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
	})
}

func TestRunDryRun(t *testing.T) {
	_ = newRootCmd()

	cfg := `secrets:
  - alias: token
    provider: exec
    name: /bin/echo
    env: TOKEN
    extras:
      args: secret
      trim: "true"`

	withTestConfig(t, cfg, func(_ string) {
		cmd := newRunCmd()
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetArgs([]string{"--all", "--dry-run", "--", "echo", "test"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("run dry-run: %v", err)
		}

		output := out.String()
		assertStringContains(t, output, []string{"TOKEN"})
	})
}
