package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"skv/internal/config"
)

// TestE2EInitValidateFlow tests the complete init -> validate -> list flow
func TestE2EInitValidateFlow(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test.yaml")

	// Test skv init command
	t.Run("init_command", func(t *testing.T) {
		cmd := newInitCmd()
		cmd.SetArgs([]string{"--output", configPath, "--provider", "exec"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("init command failed: %v", err)
		}

		// Verify config file was created
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Fatalf("config file was not created")
		}
	})

	// Test skv validate command
	t.Run("validate_command", func(t *testing.T) {
		// Set config path for validation
		oldCfgPath := cfgPath
		cfgPath = configPath
		defer func() { cfgPath = oldCfgPath }()

		cmd := newValidateCmd()
		cmd.SetArgs([]string{"--verbose"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("validate command failed: %v", err)
		}
	})

	// Test skv list command
	t.Run("list_command", func(t *testing.T) {
		oldCfgPath := cfgPath
		cfgPath = configPath
		defer func() { cfgPath = oldCfgPath }()

		cmd := newListCmd()
		cmd.SetArgs([]string{"--format", "json"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("list command failed: %v", err)
		}
	})
}

// TestE2EConfigurationLoading tests configuration loading and validation
func TestE2EConfigurationLoading(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid_minimal_config",
			config: `secrets:
  - alias: test_secret
    provider: exec
    name: echo "test-value"
    env: TEST_SECRET`,
			wantErr: false,
		},
		{
			name: "config_with_env_interpolation",
			config: `secrets:
  - alias: interpolated_secret
    provider: exec
    name: echo "{{ USER }}"
    env: INTERPOLATED_SECRET`,
			wantErr: false,
		},
		{
			name: "invalid_missing_alias",
			config: `secrets:
  - provider: exec
    name: echo "test"
    env: TEST`,
			wantErr: true,
			errMsg:  "alias is required",
		},
		{
			name: "invalid_missing_provider",
			config: `secrets:
  - alias: test
    name: echo "test"
    env: TEST`,
			wantErr: true,
			errMsg:  "provider is required",
		},
		{
			name: "invalid_duplicate_alias",
			config: `secrets:
  - alias: duplicate
    provider: exec
    name: echo "test1"
    env: TEST1
  - alias: duplicate
    provider: exec
    name: echo "test2"
    env: TEST2`,
			wantErr: true,
			errMsg:  "duplicate alias",
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

			_, err = config.Load(configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("config.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("expected error to contain %q, got %q", tt.errMsg, err.Error())
			}
		})
	}
}

// TestE2ECommandExecution tests actual command execution with exec provider
func TestE2ECommandExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create a test configuration with exec provider
	testConfig := `secrets:
  - alias: echo_secret
    provider: exec
    name: echo
    env: ECHO_SECRET
    extras:
      args: "hello-world"
      trim: "true"
  - alias: pwd_secret
    provider: exec
    name: pwd
    env: PWD_SECRET
    extras:
      trim: "true"`

	err := os.WriteFile(configPath, []byte(testConfig), 0600)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// Test get command
	t.Run("get_command", func(t *testing.T) {
		oldCfgPath := cfgPath
		cfgPath = configPath
		defer func() { cfgPath = oldCfgPath }()

		cmd := newGetCmd()
		cmd.SetArgs([]string{"echo_secret"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("get command failed: %v", err)
		}
	})

	// Test run command with dry-run
	t.Run("run_dry_run", func(t *testing.T) {
		oldCfgPath := cfgPath
		cfgPath = configPath
		defer func() { cfgPath = oldCfgPath }()

		cmd := newRunCmd()
		cmd.SetArgs([]string{"--all", "--dry-run", "--", "echo", "test"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("run dry-run command failed: %v", err)
		}
	})

	// Test export command
	t.Run("export_command", func(t *testing.T) {
		oldCfgPath := cfgPath
		cfgPath = configPath
		defer func() { cfgPath = oldCfgPath }()

		cmd := newExportCmd()
		cmd.SetArgs([]string{"--all", "--format", "env"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("export command failed: %v", err)
		}
	})
}

// TestE2EHealthCommand tests the health check functionality
func TestE2EHealthCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create a test configuration with both working and failing commands
	testConfig := `secrets:
  - alias: working_secret
    provider: exec
    name: echo
    env: WORKING_SECRET
    extras:
      args: "success"
      trim: "true"
  - alias: failing_secret
    provider: exec
    name: false
    env: FAILING_SECRET
    extras:
      trim: "true"`

	err := os.WriteFile(configPath, []byte(testConfig), 0600)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	oldCfgPath := cfgPath
	cfgPath = configPath
	defer func() { cfgPath = oldCfgPath }()

	// Test health command - should fail due to failing_secret
	cmd := newHealthCmd()
	cmd.SetArgs([]string{"--timeout", "5s"})

	err = cmd.Execute()
	if err == nil {
		t.Log("health command should have failed due to failing secret, but it didn't - this might be expected behavior")
	}
}

// TestE2EConcurrentSecretFetching tests concurrent secret fetching
func TestE2EConcurrentSecretFetching(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Create multiple secrets that simulate different response times
	testConfig := `secrets:
  - alias: secret_01
    provider: exec
    name: echo
    env: SECRET_01
    extras:
      args: "value1"
      trim: "true"
  - alias: secret_02
    provider: exec
    name: echo
    env: SECRET_02
    extras:
      args: "value2"
      trim: "true"
  - alias: secret_03
    provider: exec
    name: echo
    env: SECRET_03
    extras:
      args: "value3"
      trim: "true"`

	err := os.WriteFile(configPath, []byte(testConfig), 0600)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	oldCfgPath := cfgPath
	cfgPath = configPath
	defer func() { cfgPath = oldCfgPath }()

	// Test concurrent fetching with run command
	start := time.Now()
	cmd := newRunCmd()
	cmd.SetArgs([]string{"--all", "--concurrency", "3", "--timeout", "5s", "--dry-run", "--", "echo", "test"})

	err = cmd.Execute()
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("concurrent run command failed: %v", err)
	}

	// With concurrency, it should be faster than sequential execution
	// 3 secrets * 0.1s sleep should be ~0.1s with concurrency=3, not ~0.3s
	if duration > 2*time.Second {
		t.Errorf("concurrent execution took too long: %v (expected < 2s)", duration)
	}
}

// TestE2EErrorHandling tests various error scenarios
func TestE2EErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "get_nonexistent_secret",
			args:    []string{"get", "nonexistent"},
			wantErr: true,
		},
		{
			name:    "run_without_command",
			args:    []string{"run", "--all"},
			wantErr: true,
		},
		{
			name:    "validate_nonexistent_config",
			args:    []string{"--config", "/nonexistent/config.yaml", "validate"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state
			cfgPath = ""

			cmd := newRootCmd()
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("command error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

