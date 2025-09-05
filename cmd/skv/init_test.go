package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitCmd(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains []string
	}{
		{
			name:     "default_template",
			args:     []string{},
			wantErr:  false,
			contains: []string{"AWS Secrets Manager", "Google Cloud Secret Manager", "Azure Key Vault", "HashiCorp Vault", "exec"},
		},
		{
			name:     "aws_template",
			args:     []string{"--provider", "aws"},
			wantErr:  false,
			contains: []string{"AWS Secrets Manager", "AWS SSM Parameter Store"},
		},
		{
			name:     "gcp_template",
			args:     []string{"--provider", "gcp"},
			wantErr:  false,
			contains: []string{"Google Secret Manager", "projects/my-project"},
		},
		{
			name:     "azure_template",
			args:     []string{"--provider", "azure"},
			wantErr:  false,
			contains: []string{"Azure Key Vault", "Azure App Configuration"},
		},
		{
			name:     "vault_template",
			args:     []string{"--provider", "vault"},
			wantErr:  false,
			contains: []string{"HashiCorp Vault", "kv/data/", "AppRole"},
		},
		{
			name:     "exec_template",
			args:     []string{"--provider", "exec"},
			wantErr:  false,
			contains: []string{"exec", "scripts/", "./scripts/"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			outputPath := filepath.Join(tempDir, "test-config.yaml")

			args := append(tt.args, "--output", outputPath)

			cmd := newInitCmd()
			cmd.SetArgs(args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("init command error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was created
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					t.Fatalf("config file was not created at %s", outputPath)
				}

				// Read and verify content
				// #nosec G304: outputPath is controlled by test and is safe
				content, err := os.ReadFile(outputPath)
				if err != nil {
					t.Fatalf("failed to read generated config: %v", err)
				}

				contentStr := string(content)
				for _, expectedContent := range tt.contains {
					if !strings.Contains(contentStr, expectedContent) {
						t.Errorf("generated config does not contain expected content %q", expectedContent)
					}
				}

				// Verify it's valid YAML by checking basic structure
				if !strings.Contains(contentStr, "secrets:") {
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
	originalContent := "existing content"
	err := os.WriteFile(outputPath, []byte(originalContent), 0600)
	if err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Try to init with force flag - should succeed
	cmd := newInitCmd()
	cmd.SetArgs([]string{"--output", outputPath, "--force", "--provider", "exec"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("init command with --force failed: %v", err)
	}

	// Verify file was overwritten
	// #nosec G304: outputPath is controlled by test and is safe
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read overwritten file: %v", err)
	}

	if string(content) == originalContent {
		t.Error("file was not overwritten despite --force flag")
	}

	if !strings.Contains(string(content), "secrets:") {
		t.Error("overwritten file does not contain valid config structure")
	}
}

func TestInitCmdDefaultOutput(t *testing.T) {
	// This test verifies that default output path logic works
	// We can't easily test the actual home directory path, but we can test the logic

	cmd := newInitCmd()
	// Don't execute, just verify the command was created properly
	if cmd == nil {
		t.Fatal("newInitCmd() returned nil")
	}

	if cmd.Use != "init" {
		t.Errorf("expected command name 'init', got %q", cmd.Use)
	}

	if !strings.Contains(cmd.Short, "configuration template") {
		t.Errorf("expected short description to mention 'configuration template', got %q", cmd.Short)
	}
}

