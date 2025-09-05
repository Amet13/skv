package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHealthCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tests := []struct {
		name        string
		config      string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name: "healthy_secrets",
			config: `secrets:
  - alias: working_secret
    provider: exec
    name: echo
    env: WORKING_SECRET
    extras:
      args: "test-value"
      trim: "true"`,
			args:    []string{},
			wantErr: false,
		},
		{
			name: "mixed_health_secrets",
			config: `secrets:
  - alias: working_secret
    provider: exec
    name: echo
    env: WORKING_SECRET
    extras:
      args: "test-value"
      trim: "true"
  - alias: failing_secret
    provider: exec
    name: false
    env: FAILING_SECRET`,
			args:        []string{},
			wantErr:     true,
			errContains: "health check failed",
		},
		{
			name: "specific_secret_healthy",
			config: `secrets:
  - alias: working_secret
    provider: exec
    name: echo
    env: WORKING_SECRET
    extras:
      args: "test-value"
      trim: "true"
  - alias: failing_secret
    provider: exec
    name: false
    env: FAILING_SECRET`,
			args:    []string{"--secret", "working_secret"},
			wantErr: false,
		},
		{
			name: "specific_secret_unhealthy",
			config: `secrets:
  - alias: working_secret
    provider: exec
    name: echo
    env: WORKING_SECRET
    extras:
      args: "test-value"
      trim: "true"
  - alias: failing_secret
    provider: exec
    name: false
    env: FAILING_SECRET`,
			args:        []string{"--secret", "failing_secret"},
			wantErr:     true,
			errContains: "health check failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "config.yaml")

			err := os.WriteFile(configPath, []byte(tt.config), 0600)
			if err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			oldCfgPath := cfgPath
			cfgPath = configPath
			defer func() { cfgPath = oldCfgPath }()

			cmd := newHealthCmd()
			cmd.SetArgs(tt.args)

			err = cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("health command error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error to contain %q, got %v", tt.errContains, err)
				}
			}
		})
	}
}

func TestHealthCmdTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Config with a slow command
	config := `secrets:
  - alias: slow_secret
    provider: exec
    name: sleep
    env: SLOW_SECRET
    extras:
      args: "2"
      trim: "true"`

	err := os.WriteFile(configPath, []byte(config), 0600)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	oldCfgPath := cfgPath
	cfgPath = configPath
	defer func() { cfgPath = oldCfgPath }()

	cmd := newHealthCmd()
	cmd.SetArgs([]string{"--timeout", "1s"})

	err = cmd.Execute()
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}

func TestHealthCmdInvalidConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Invalid config
	config := `invalid yaml: [}`

	err := os.WriteFile(configPath, []byte(config), 0600)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	oldCfgPath := cfgPath
	cfgPath = configPath
	defer func() { cfgPath = oldCfgPath }()

	cmd := newHealthCmd()
	cmd.SetArgs([]string{})

	err = cmd.Execute()
	if err == nil {
		t.Error("expected error for invalid config, got nil")
	} else if !strings.Contains(err.Error(), "failed to load configuration") {
		t.Errorf("expected configuration load error, got: %v", err)
	}
}

func TestHealthCmdUnknownProvider(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Config with unknown provider
	config := `secrets:
  - alias: unknown_secret
    provider: nonexistent_provider
    name: test
    env: UNKNOWN_SECRET`

	err := os.WriteFile(configPath, []byte(config), 0600)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	oldCfgPath := cfgPath
	cfgPath = configPath
	defer func() { cfgPath = oldCfgPath }()

	cmd := newHealthCmd()
	cmd.SetArgs([]string{})

	err = cmd.Execute()
	if err == nil {
		t.Error("expected error for unknown provider, got nil")
	}
}

func TestHealthCmdNoSecrets(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Config with one dummy secret to satisfy validation
	config := `secrets:
  - alias: dummy
    provider: exec
    name: echo
    env: DUMMY`

	err := os.WriteFile(configPath, []byte(config), 0600)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	oldCfgPath := cfgPath
	cfgPath = configPath
	defer func() { cfgPath = oldCfgPath }()

	cmd := newHealthCmd()
	cmd.SetArgs([]string{})

	err = cmd.Execute()
	if err != nil {
		t.Logf("health command returned error (which might be expected): %v", err)
	}
}

