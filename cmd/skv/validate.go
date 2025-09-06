package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"skv/internal/config"
	"skv/internal/provider"
)

func newValidateCmd() *cobra.Command {
	var (
		checkProviders bool
		checkSecrets   bool
		verbose        bool
	)

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration file",
		Long: `Validate the configuration file for syntax errors, missing providers,
and connectivity issues. This command helps ensure your configuration
is correct before using it in production.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("configuration validation failed: %w", err)
			}

			fmt.Printf("Configuration loaded successfully from %s\n", getConfigPath())
			fmt.Printf("Found %d secret(s) configured\n", len(cfg.Secrets))

			if verbose {
				fmt.Println("\nConfiguration summary:")
				for _, secret := range cfg.Secrets {
					spec := secret.ToSpec()
					fmt.Printf("  - %s (%s) -> %s\n", secret.Alias, secret.Provider, spec.EnvName)
				}
			}

			// Check provider registration
			if checkProviders {
				fmt.Println("\nChecking provider availability...")
				providerIssues := 0
				for _, secret := range cfg.Secrets {
					if _, ok := provider.Get(secret.Provider); !ok {
						fmt.Printf("ERROR: Provider '%s' not found for secret '%s'\n", secret.Provider, secret.Alias)
						providerIssues++
					} else if verbose {
						fmt.Printf("Provider '%s' available for secret '%s'\n", secret.Provider, secret.Alias)
					}
				}
				if providerIssues > 0 {
					return fmt.Errorf("found %d provider issues", providerIssues)
				}
				fmt.Println("All providers are available")
			}

			// Test secret connectivity (dry-run fetch)
			if checkSecrets {
				fmt.Println("\nTesting secret connectivity...")
				secretIssues := 0
				for _, secret := range cfg.Secrets {
					spec := secret.ToSpec()
					p, ok := provider.Get(spec.Provider)
					if !ok {
						fmt.Printf("ERROR: Provider '%s' not available for secret '%s'\n", spec.Provider, secret.Alias)
						secretIssues++
						continue
					}

					// Test with a short timeout
					ctx := cmd.Context()
					_, err := p.FetchSecret(ctx, spec)
					if err != nil {
						if err == provider.ErrNotFound {
							fmt.Printf("WARNING: Secret '%s' not found in provider '%s'\n", secret.Alias, spec.Provider)
						} else {
							fmt.Printf("ERROR: Error fetching secret '%s': %v\n", secret.Alias, err)
						}
						secretIssues++
					} else if verbose {
						fmt.Printf("Secret '%s' accessible\n", secret.Alias)
					}
				}
				if secretIssues > 0 {
					fmt.Printf("\nWARNING: Found %d connectivity issues (this might be expected in some environments)\n", secretIssues)
				} else {
					fmt.Println("All secrets are accessible")
				}
			}

			fmt.Println("\nValidation completed successfully!")
			return nil
		},
	}

	cmd.Flags().BoolVar(&checkProviders, "check-providers", true, "Verify all providers are available")
	cmd.Flags().BoolVar(&checkSecrets, "check-secrets", false, "Test connectivity to all secrets (requires valid credentials)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed validation results")

	return cmd
}

func getConfigPath() string {
	if cfgPath != "" {
		return cfgPath
	}
	if env := os.Getenv("SKV_CONFIG"); env != "" {
		return env
	}
	home, _ := os.UserHomeDir()
	return home + "/.skv.yaml"
}
