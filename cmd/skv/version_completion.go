package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	skvversion "skv/internal/version"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, _ []string) {
			out := cmd.OutOrStdout()
			_, _ = fmt.Fprintln(out, skvversion.String())
		},
	}
}

func newCompletionCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := args[0]
			out := cmd.OutOrStdout()
			switch shell {
			case "bash":
				return cmd.Root().GenBashCompletionV2(out, true)
			case "zsh":
				return cmd.Root().GenZshCompletion(out)
			case "fish":
				return cmd.Root().GenFishCompletion(out, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(out)
			default:
				return fmt.Errorf("unsupported shell: %s", shell)
			}
		},
	}

	// Add install subcommand
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install shell completions automatically",
		Long: `Install shell completions for the current shell.

This command automatically detects your shell and installs completions
in the appropriate location. It supports bash, zsh, and fish shells.`,
		RunE: installCompletions,
	}
	c.AddCommand(installCmd)

	return c
}

func installCompletions(cmd *cobra.Command, _ []string) error {
	// Detect current shell
	shell := detectShell()
	if shell == "" {
		return fmt.Errorf("unable to detect shell; please specify shell manually: skv completion [bash|zsh|fish|powershell]")
	}

	fmt.Printf("Detected shell: %s\n", shell)

	var completionPath string
	var err error

	switch shell {
	case "bash":
		completionPath, err = installBashCompletion(cmd)
	case "zsh":
		completionPath, err = installZshCompletion(cmd)
	case "fish":
		completionPath, err = installFishCompletion(cmd)
	default:
		return fmt.Errorf("shell '%s' is not supported for automatic installation", shell)
	}

	if err != nil {
		return fmt.Errorf("failed to install completions: %w", err)
	}

	fmt.Printf("OK: Completions installed successfully!\n")
	fmt.Printf("Location: %s\n", completionPath)
	fmt.Println("INFO: Restart your shell or run 'source' on your profile to activate completions")
	fmt.Println("TIP: You can test with: skv <TAB>")

	return nil
}

func detectShell() string {
	// Check SHELL environment variable
	if shell := os.Getenv("SHELL"); shell != "" {
		switch {
		case strings.Contains(shell, "bash"):
			return "bash"
		case strings.Contains(shell, "zsh"):
			return "zsh"
		case strings.Contains(shell, "fish"):
			return "fish"
		}
	}

	// Check if we're running in a known shell
	if ps := os.Getenv("PS1"); ps != "" {
		// Likely bash
		return "bash"
	}

	return ""
}

func installBashCompletion(cmd *cobra.Command) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Create completion directory if it doesn't exist
	completionDir := filepath.Join(homeDir, ".local", "share", "bash-completion", "completions")
	if err := os.MkdirAll(completionDir, 0750); err != nil {
		return "", err
	}

	completionPath := filepath.Join(completionDir, "skv")

	// Generate and write completion script
	file, err := os.Create(completionPath) // #nosec G304 - completionPath is derived from user home directory, safe for this use case
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close() // Ignore error in defer
	}()

	if err := cmd.Root().GenBashCompletionV2(file, true); err != nil {
		return "", err
	}

	return completionPath, nil
}

func installZshCompletion(cmd *cobra.Command) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Create zsh functions directory
	zfuncDir := filepath.Join(homeDir, ".zfunc")
	if err := os.MkdirAll(zfuncDir, 0750); err != nil {
		return "", err
	}

	completionPath := filepath.Join(zfuncDir, "_skv")

	// Generate and write completion script
	file, err := os.Create(completionPath) // #nosec G304 - completionPath is derived from user home directory, safe for this use case
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close() // Ignore error in defer
	}()

	if err := cmd.Root().GenZshCompletion(file); err != nil {
		return "", err
	}

	return completionPath, nil
}

func installFishCompletion(cmd *cobra.Command) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Create fish completions directory
	completionDir := filepath.Join(homeDir, ".config", "fish", "completions")
	if err := os.MkdirAll(completionDir, 0750); err != nil {
		return "", err
	}

	completionPath := filepath.Join(completionDir, "skv.fish")

	// Generate and write completion script
	file, err := os.Create(completionPath) // #nosec G304 - completionPath is derived from user home directory, safe for this use case
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close() // Ignore error in defer
	}()

	if err := cmd.Root().GenFishCompletion(file, true); err != nil {
		return "", err
	}

	return completionPath, nil
}

