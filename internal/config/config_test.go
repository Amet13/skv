package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config",
			config: `
secrets:
  - alias: db_password
    provider: aws
    name: test/db_password
    env: DB_PASSWORD
defaults:
  region: us-east-1
`,
			wantErr: false,
		},
		{
			name: "invalid yaml",
			config: `
secrets:
  - alias: db_password
    provider: aws
    name: test/db_password
    env: DB_PASSWORD
  invalid yaml here
`,
			wantErr:     true,
			errContains: "parse config",
		},
		{
			name: "missing required fields",
			config: `
secrets:
  - alias: db_password
    provider: aws
    name: test
`,
			wantErr: false, // Should still load, validation happens elsewhere
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configFile := filepath.Join(tempDir, "test-config.yaml")

			err := os.WriteFile(configFile, []byte(tt.config), 0600)
			if err != nil {
				t.Fatal(err)
			}

			_, err = Load(configFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Load() error = %v, expected to contain %v", err, tt.errContains)
			}
		})
	}
}

func TestLocateConfigPath(t *testing.T) {
	// Test with override path
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test-config.yaml")
	err := os.WriteFile(configFile, []byte("secrets: []"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	path := locateConfigPath(configFile)
	if path != configFile {
		t.Errorf("locateConfigPath() = %v, want %v", path, configFile)
	}

	// Test with no path (should find default locations)
	path = locateConfigPath("")
	if path != "" {
		// This will vary by system, just ensure it's not empty when a config exists
		homeDir, _ := os.UserHomeDir()
		expectedPaths := []string{
			filepath.Join(homeDir, ".skv.yaml"),
			filepath.Join(homeDir, ".skv.yml"),
		}

		found := false
		for _, expected := range expectedPaths {
			if path == expected {
				found = true
				break
			}
		}
		if !found {
			t.Logf("locateConfigPath() returned unexpected path: %v", path)
		}
	}
}

func TestInterpolateEnv(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		envKey   string
		envValue string
		expected string
	}{
		{
			name:     "simple interpolation",
			input:    "{{ TEST_VAR }}",
			envKey:   "TEST_VAR",
			envValue: "test_value",
			expected: "test_value",
		},
		{
			name:     "no interpolation needed",
			input:    "no_vars_here",
			expected: "no_vars_here",
		},
		{
			name:     "missing env var",
			input:    "{{ MISSING_VAR }}",
			expected: "__MISSING_ENV_MISSING_VAR__",
		},
		{
			name:     "multiple interpolations",
			input:    "{{ VAR1 }}_{{ VAR2 }}",
			envKey:   "VAR1",
			envValue: "value1",
			expected: "value1___MISSING_ENV_VAR2__", // Only VAR1 is set, VAR2 is missing
		},
		{
			name:     "whitespace handling",
			input:    "{{  TEST_VAR  }}",
			envKey:   "TEST_VAR",
			envValue: "test_value",
			expected: "test_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.envKey != "" {
				oldValue := os.Getenv(tt.envKey)
				_ = os.Setenv(tt.envKey, tt.envValue)
				defer func() {
					if oldValue == "" {
						_ = os.Unsetenv(tt.envKey)
					} else {
						_ = os.Setenv(tt.envKey, oldValue)
					}
				}()
			}

			result := interpolateEnv(tt.input)
			if result != tt.expected {
				t.Errorf("interpolateEnv(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTransformValue(t *testing.T) {
	tests := []struct {
		name     string
		secret   Secret
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "no transform",
			secret:   Secret{},
			input:    "test_value",
			expected: "test_value",
			wantErr:  false,
		},
		{
			name: "template transform",
			secret: Secret{
				Transform: &Transform{
					Type:     "template",
					Template: "prefix_{{ .value }}_suffix",
				},
			},
			input:    "test",
			expected: "prefix_test_suffix",
			wantErr:  false,
		},
		{
			name: "mask transform",
			secret: Secret{
				Transform: &Transform{
					Type: "mask",
				},
			},
			input:    "verylongpassword",
			expected: "ve************rd",
			wantErr:  false,
		},
		{
			name: "prefix transform",
			secret: Secret{
				Transform: &Transform{
					Type:   "prefix",
					Prefix: "prefix_",
				},
			},
			input:    "value",
			expected: "prefix_value",
			wantErr:  false,
		},
		{
			name: "suffix transform",
			secret: Secret{
				Transform: &Transform{
					Type:   "suffix",
					Suffix: "_suffix",
				},
			},
			input:    "value",
			expected: "value_suffix",
			wantErr:  false,
		},
		{
			name: "unknown transform type",
			secret: Secret{
				Transform: &Transform{
					Type: "unknown",
				},
			},
			input:   "value",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.secret.TransformValue(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("TransformValue() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestDeriveEnvName(t *testing.T) {
	tests := []struct {
		name     string
		alias    string
		expected string
	}{
		{
			name:     "simple alias",
			alias:    "db_password",
			expected: "DB_PASSWORD",
		},
		{
			name:     "with underscores",
			alias:    "api_key",
			expected: "API_KEY",
		},
		{
			name:     "camelCase",
			alias:    "databaseUrl",
			expected: "DATABASE_URL",
		},
		{
			name:     "with numbers",
			alias:    "db2_password",
			expected: "DB2_PASSWORD",
		},
		{
			name:     "empty alias",
			alias:    "",
			expected: "SECRET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deriveEnvName(tt.alias)
			if result != tt.expected {
				t.Errorf("deriveEnvName(%q) = %q, want %q", tt.alias, result, tt.expected)
			}
		})
	}
}

func TestContainsMissingEnvToken(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected bool
	}{
		{
			name:     "no missing tokens",
			values:   []string{"value1", "value2"},
			expected: false,
		},
		{
			name:     "has missing token",
			values:   []string{"value1", "__MISSING_ENV_VAR__", "value2"},
			expected: true,
		},
		{
			name:     "empty slice",
			values:   []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsMissingEnvToken(tt.values...)
			if result != tt.expected {
				t.Errorf("containsMissingEnvToken(%v) = %v, want %v", tt.values, result, tt.expected)
			}
		})
	}
}

func BenchmarkConfigLoad(b *testing.B) {
	// Create a test config file
	testConfig := `
secrets:
  - alias: db_password
    provider: aws
    name: test/db_password
    env: DB_PASSWORD
  - alias: api_key
    provider: gcp
    name: projects/test/secrets/api_key/versions/latest
    env: API_KEY
defaults:
  region: us-east-1
`

	// Write test config to a temporary file for benchmarking
	tempFile := "benchmark-config.yaml"
	err := os.WriteFile(tempFile, []byte(testConfig), 0600)
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = os.Remove(tempFile) }()

	b.Run("Load", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := Load(tempFile)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

