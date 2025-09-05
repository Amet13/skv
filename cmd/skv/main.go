// skv is a CLI for fetching secrets and injecting them into processes.
package main

import (
	"log/slog"
	"os"
	"strings"

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
	logFmt   string
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
			var level slog.Level
			switch strings.ToLower(strings.TrimSpace(logLevel)) {
			case "debug":
				level = slog.LevelDebug
			case "warn":
				level = slog.LevelWarn
			case "error":
				level = slog.LevelError
			case "info", "":
				level = slog.LevelInfo
			default:
				level = slog.LevelInfo
			}
			var handler slog.Handler
			if strings.EqualFold(logFmt, "json") {
				handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})
			} else {
				handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
			}
			slog.SetDefault(slog.New(handler))
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&cfgPath, "config", "", "Path to config file (overrides SKV_CONFIG and default $HOME/.skv.yaml)")
	cmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level: error|warn|info|debug")
	cmd.PersistentFlags().StringVar(&logFmt, "log-format", "text", "Log format: text|json")

	// Register providers
	provider.Register("aws-secrets-manager", awsprovider.New())
	// Aliases for convenience and docs consistency
	provider.Register("aws", awsprovider.New())
	// AWS SSM Parameter Store
	provider.Register("aws-ssm", awsprovider.NewSSM())
	provider.Register("ssm", awsprovider.NewSSM())
	provider.Register("aws-parameter-store", awsprovider.NewSSM())
	provider.Register("vault", vaultprovider.New())
	provider.Register("gcp-secret-manager", gcpprovider.New())
	provider.Register("gcp", gcpprovider.New())
	provider.Register("azure-key-vault", azureprovider.New())
	provider.Register("azure", azureprovider.New())
	// Azure App Configuration (parameter store)
	provider.Register("azure-appconfig", azureprovider.NewAppConfig())
	provider.Register("appconfig", azureprovider.NewAppConfig())
	provider.Register("exec", execprovider.New())

	// Future/stub providers for extensibility (return clear not-implemented errors)
	provider.Register("oci", provider.NewNotImplementedProvider("oci"))
	provider.Register("ibm", provider.NewNotImplementedProvider("ibm"))
	provider.Register("alibaba", provider.NewNotImplementedProvider("alibaba"))

	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newRunCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newExportCmd())
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newCompletionCmd())

	return cmd
}

