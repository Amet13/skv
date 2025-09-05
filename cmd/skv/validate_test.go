package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateCmd(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid_config",
			config: `secrets:
  - alias: test_secret
    provider: exec
    name: echo "test"
    env: TEST_SECRET`,
			args:    []string{},
			wantErr: false,
		},
		{
			name: "valid_config_verbose",
			config: `secrets:
  - alias: test_secret
    provider: exec
    name: echo "test"
    env: TEST_SECRET
  - alias: another_secret
    provider: exec
    name: echo "another"
    env: ANOTHER_SECRET`,
			args:    []string{"--verbose"},
			wantErr: false,
		},
		{
			name: "invalid_config_missing_alias",
			config: `secrets:
  - provider: exec
    name: echo "test"
    env: TEST_SECRET`,
			args:        []string{},
			wantErr:     true,
			errContains: "alias is required",
		},
		{
			name: "invalid_config_duplicate_alias",
			config: `secrets:
  - alias: duplicate
    provider: exec
    name: echo "test1"
    env: TEST1
  - alias: duplicate
    provider: exec
    name: echo "test2"
    env: TEST2`,
			args:        []string{},
			wantErr:     true,
			errContains: "duplicate alias",
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

			// Set the config path
			oldCfgPath := cfgPath
			cfgPath = configPath
			defer func() { cfgPath = oldCfgPath }()

			cmd := newValidateCmd()
			cmd.SetArgs(tt.args)

			err = cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("validate command error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error to contain %q, got %v", tt.errContains, err)
				}
			}
		})
	}
}

func TestValidateCmdProviderCheck(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Config with unknown provider
	config := `secrets:
  - alias: test_secret
    provider: nonexistent_provider
    name: test
    env: TEST_SECRET`

	err := os.WriteFile(configPath, []byte(config), 0600)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	oldCfgPath := cfgPath
	cfgPath = configPath
	defer func() { cfgPath = oldCfgPath }()

	cmd := newValidateCmd()
	cmd.SetArgs([]string{"--check-providers"})

	err = cmd.Execute()
	if err == nil {
		t.Error("expected error for unknown provider, got nil")
	} else if !contains(err.Error(), "provider issues") {
		t.Errorf("expected error about provider issues, got: %v", err)
	}
}

func TestValidateCmdSecretCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Config with exec provider that should work
	config := `secrets:
  - alias: working_secret
    provider: exec
    name: echo
    env: WORKING_SECRET
    extras:
      args: "test-value"
      trim: "true"`

	err := os.WriteFile(configPath, []byte(config), 0600)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	oldCfgPath := cfgPath
	cfgPath = configPath
	defer func() { cfgPath = oldCfgPath }()

	cmd := newValidateCmd()
	cmd.SetArgs([]string{"--check-secrets", "--verbose"})

	err = cmd.Execute()
	if err != nil {
		t.Errorf("validate command with working secret failed: %v", err)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

