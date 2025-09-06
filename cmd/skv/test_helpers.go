package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// Common test configurations
const (
	validExecConfig = `secrets:
  - alias: test_secret
    provider: exec
    name: echo "test"
    env: TEST_SECRET`

	validExecConfigWithArgs = `secrets:
  - alias: test_secret
    provider: exec
    name: echo
    env: TEST_SECRET
    extras:
      args: "test-value"
      trim: "true"`

	// #nosec G101: This is test configuration, not real credentials
	multiSecretConfig = `secrets:
  - alias: test_secret
    provider: exec
    name: echo "test"
    env: TEST_SECRET
  - alias: another_secret
    provider: exec
    name: echo "another"
    env: ANOTHER_SECRET`

	invalidMissingAlias = `secrets:
  - provider: exec
    name: echo "test"
    env: TEST_SECRET`

	invalidMissingProvider = `secrets:
  - alias: test
    name: echo "test"
    env: TEST`

	invalidDuplicateAlias = `secrets:
  - alias: duplicate
    provider: exec
    name: echo "test1"
    env: TEST1
  - alias: duplicate
    provider: exec
    name: echo "test2"
    env: TEST2`
)

// TestConfig represents a test configuration
type TestConfig struct {
	Name        string
	Config      string
	Args        []string
	WantErr     bool
	ErrContains string
	Contains    []string // For content validation
}

// writeTestConfig creates a temporary config file with the given content
func writeTestConfig(t *testing.T, content string) string {
	t.Helper()
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test-config.yaml")

	err := os.WriteFile(configPath, []byte(content), 0600)
	if err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	return configPath
}

// withTestConfig runs a test function with a temporary config file
func withTestConfig(t *testing.T, content string, fn func(configPath string)) {
	t.Helper()
	configPath := writeTestConfig(t, content)

	// Save and restore original config path
	oldCfgPath := cfgPath
	cfgPath = configPath
	defer func() { cfgPath = oldCfgPath }()

	fn(configPath)
}

// runCommandTest runs a command with the given config and args, checking for expected results
func runCommandTest(t *testing.T, tc TestConfig, cmdFactory func() *cobra.Command) {
	t.Helper()

	withTestConfig(t, tc.Config, func(_ string) {
		cmd := cmdFactory()
		cmd.SetArgs(tc.Args)

		err := cmd.Execute()

		if tc.WantErr {
			if err == nil {
				t.Errorf("expected error but got none")
				return
			}
			if tc.ErrContains != "" && !strings.Contains(err.Error(), tc.ErrContains) {
				t.Errorf("expected error to contain %q, got %q", tc.ErrContains, err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}
	})
}

// runTableTests runs a table of test cases using the provided command factory
func runTableTests(t *testing.T, tests []TestConfig, cmdFactory func() *cobra.Command) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			runCommandTest(t, tt, cmdFactory)
		})
	}
}

// Common validation test cases
func getValidationTestCases() []TestConfig {
	return []TestConfig{
		{
			Name:    "valid_config",
			Config:  validExecConfig,
			Args:    []string{},
			WantErr: false,
		},
		{
			Name:    "valid_config_verbose",
			Config:  multiSecretConfig,
			Args:    []string{"--verbose"},
			WantErr: false,
		},
		{
			Name:        "invalid_config_missing_alias",
			Config:      invalidMissingAlias,
			Args:        []string{},
			WantErr:     true,
			ErrContains: "alias is required",
		},
		{
			Name:        "invalid_config_missing_provider",
			Config:      invalidMissingProvider,
			Args:        []string{},
			WantErr:     true,
			ErrContains: "provider is required",
		},
		{
			Name:        "invalid_config_duplicate_alias",
			Config:      invalidDuplicateAlias,
			Args:        []string{},
			WantErr:     true,
			ErrContains: "duplicate alias",
		},
	}
}

// Common health test cases
func getHealthTestCases() []TestConfig {
	return []TestConfig{
		{
			Name:    "healthy_secrets",
			Config:  validExecConfigWithArgs,
			Args:    []string{},
			WantErr: false,
		},
		{
			Name: "mixed_health_secrets",
			Config: `secrets:
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
			Args:        []string{},
			WantErr:     true,
			ErrContains: "health check failed",
		},
		{
			Name: "specific_secret_healthy",
			Config: `secrets:
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
			Args:    []string{"--secret", "working_secret"},
			WantErr: false,
		},
	}
}

// skipIfShort skips the test if running in short mode
func skipIfShort(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
}

// assertFileExists checks that a file exists at the given path
func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("file does not exist: %s", path)
	}
}

// assertFileContains checks that a file contains the expected content
func assertFileContains(t *testing.T, path string, expected []string) {
	t.Helper()
	// #nosec G304: path is controlled by test and is safe
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}

	contentStr := string(content)
	for _, exp := range expected {
		if !strings.Contains(contentStr, exp) {
			t.Errorf("file %s does not contain expected content %q", path, exp)
		}
	}
}

// assertStringContains checks that a string contains expected substrings
func assertStringContains(t *testing.T, str string, expected []string) {
	t.Helper()
	for _, exp := range expected {
		if !strings.Contains(str, exp) {
			t.Errorf("string does not contain expected content %q", exp)
		}
	}
}

