package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeriveEnvName(t *testing.T) {
	cases := map[string]string{
		"db-password":     "DB_PASSWORD",
		"My.API.Token":    "MY_API_TOKEN",
		"__weird__name__": "WEIRD_NAME",
		"":                "SECRET",
	}
	for in, want := range cases {
		got := deriveEnvName(in)
		if got != want {
			t.Fatalf("deriveEnvName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLoadInterpolatesEnv(t *testing.T) {
	t.Setenv("FOO", "bar")
	dir := t.TempDir()
	path := filepath.Join(dir, ".skv.yaml")
	data := []byte("secrets:\n  - alias: test\n    provider: aws-secrets-manager\n    name: \"{{ FOO }}\"\n")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Setenv("SKV_CONFIG", path); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if len(cfg.Secrets) != 1 || cfg.Secrets[0].Name != "bar" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadMissingEnvFails(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".skv.yaml")
	data := []byte("secrets:\n  - alias: test\n    provider: aws-secrets-manager\n    name: \"{{ MISSING_VAR }}\"\n")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SKV_CONFIG", path)
	if _, err := Load(""); err == nil {
		t.Fatalf("expected error due to missing env var")
	}
}
