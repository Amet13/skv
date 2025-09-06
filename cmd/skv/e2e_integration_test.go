package main

import (
	"path/filepath"
	"testing"
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

		assertFileExists(t, configPath)
	})

	// Test skv validate command with generated config
	t.Run("validate_command", func(t *testing.T) {
		// Use the generated config from init
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

	// Test skv list command with generated config
	t.Run("list_command", func(t *testing.T) {
		// Use the generated config from init
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

// TestE2ECommandExecution tests actual command execution with exec provider
func TestE2ECommandExecution(t *testing.T) {
	skipIfShort(t)

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

	withTestConfig(t, testConfig, func(_ string) {
		// Test get command
		t.Run("get_command", func(t *testing.T) {
			cmd := newGetCmd()
			cmd.SetArgs([]string{"echo_secret"})

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("get command failed: %v", err)
			}
		})

		// Test run command with dry-run
		t.Run("run_dry_run", func(t *testing.T) {
			cmd := newRunCmd()
			cmd.SetArgs([]string{"--all", "--dry-run", "--", "echo", "test"})

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("run dry-run command failed: %v", err)
			}
		})

		// Test export command
		t.Run("export_command", func(t *testing.T) {
			cmd := newExportCmd()
			cmd.SetArgs([]string{"--all", "--format", "env"})

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("export command failed: %v", err)
			}
		})
	})
}

// TestE2EHealthCommand tests the health check functionality
func TestE2EHealthCommand(t *testing.T) {
	skipIfShort(t)

	testConfig := `secrets:
  - alias: healthy_secret
    provider: exec
    name: echo
    env: HEALTHY_SECRET
    extras:
      args: "test-value"
      trim: "true"`

	withTestConfig(t, testConfig, func(_ string) {
		cmd := newHealthCmd()
		cmd.SetArgs([]string{"--timeout", "10s"})

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("health command failed: %v", err)
		}
	})
}
