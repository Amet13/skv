// skv is a CLI for fetching secrets and injecting them into processes.
package main

import (
	"os"

	"github.com/spf13/cobra"

	"skv/internal/provider"
	awsprovider "skv/internal/provider/aws"
	azureprovider "skv/internal/provider/azure"
	execprovider "skv/internal/provider/exec"
	gcpprovider "skv/internal/provider/gcp"
	vaultprovider "skv/internal/provider/vault"
	skvversion "skv/internal/version"
)

var (
	cfgPath  string
	logLevel string
)

func main() {
	root := newRootCmd()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "skv",
		Short:   "Secure Key/Value Manager",
		Version: skvversion.String(),
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			_ = logLevel
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&cfgPath, "config", "", "Path to config file (overrides SKV_CONFIG and default $HOME/.skv.yaml)")
	cmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level: error|warn|info|debug")

	// Register providers
	provider.Register("aws-secrets-manager", awsprovider.New())
	provider.Register("vault", vaultprovider.New())
	provider.Register("gcp-secret-manager", gcpprovider.New())
	provider.Register("azure-key-vault", azureprovider.New())
	provider.Register("exec", execprovider.New())

	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newRunCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newExportCmd())
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newCompletionCmd())

	return cmd
}

