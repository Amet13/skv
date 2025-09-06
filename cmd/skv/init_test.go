package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitCmd(t *testing.T) {
	tests := []TestConfig{
		{
			Name:     "default_template",
			Args:     []string{},
			WantErr:  false,
			Contains: []string{"AWS Secrets Manager", "Google Cloud Secret Manager", "Azure Key Vault", "HashiCorp Vault", "exec"},
		},
		{
			Name:     "aws_template",
			Args:     []string{"--provider", "aws"},
			WantErr:  false,
			Contains: []string{"AWS Secrets Manager", "AWS SSM Parameter Store"},
		},
		{
			Name:     "gcp_template",
			Args:     []string{"--provider", "gcp"},
			WantErr:  false,
			Contains: []string{"Google Secret Manager", "projects/my-project"},
		},
		{
			Name:     "azure_template",
			Args:     []string{"--provider", "azure"},
			WantErr:  false,
			Contains: []string{"Azure Key Vault", "Azure App Configuration"},
		},
		{
			Name:     "vault_template",
			Args:     []string{"--provider", "vault"},
			WantErr:  false,
			Contains: []string{"HashiCorp Vault", "kv/data/", "AppRole"},
		},
		{
			Name:     "exec_template",
			Args:     []string{"--provider", "exec"},
			WantErr:  false,
			Contains: []string{"Exec provider", "provider: exec", "args:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tempDir := t.TempDir()
			outputPath := filepath.Join(tempDir, "test-config.yaml")

			args := append(tt.Args, "--output", outputPath)

			cmd := newInitCmd()
			cmd.SetArgs(args)

			err := cmd.Execute()
			if (err != nil) != tt.WantErr {
				t.Errorf("init command error = %v, wantErr %v", err, tt.WantErr)
				return
			}

			if !tt.WantErr {
				assertFileExists(t, outputPath)
				assertFileContains(t, outputPath, tt.Contains)

				// Verify it's valid YAML by checking basic structure
				// #nosec G304: outputPath is controlled by test and is safe
				content, err := os.ReadFile(outputPath)
				if err != nil {
					t.Fatalf("failed to read generated config: %v", err)
				}

				if !strings.Contains(string(content), "secrets:") {
					t.Error("generated config does not contain 'secrets:' section")
				}
			}
		})
	}
}

func TestInitCmdFileExists(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "existing-config.yaml")

	// Create an existing file
	err := os.WriteFile(outputPath, []byte("existing content"), 0600)
	if err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Try to init without force flag - should fail
	cmd := newInitCmd()
	cmd.SetArgs([]string{"--output", outputPath})

	err = cmd.Execute()
	if err == nil {
		t.Error("expected error when file exists without --force, got nil")
	} else if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected error about file existing, got: %v", err)
	}
}

func TestInitCmdForceOverwrite(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "existing-config.yaml")

	// Create an existing file
	err := os.WriteFile(outputPath, []byte("existing content"), 0600)
	if err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Try to init with force flag - should succeed
	cmd := newInitCmd()
	cmd.SetArgs([]string{"--output", outputPath, "--force"})

	err = cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error with --force: %v", err)
	}

	assertFileExists(t, outputPath)

	// Verify content was overwritten
	// #nosec G304: outputPath is controlled by test and is safe
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read overwritten file: %v", err)
	}

	if string(content) == "existing content" {
		t.Error("file was not overwritten")
	}
}
