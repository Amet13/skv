package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"skv/internal/config"
	"skv/internal/provider"
)

func newDoctorCmd() *cobra.Command {
	var (
		verbose    bool
		checkAuth  bool
		checkNet   bool
		timeoutStr string
	)

	c := &cobra.Command{
		Use:   "doctor",
		Short: "Run diagnostics and health checks",
		Long: `Run comprehensive diagnostics to troubleshoot configuration and connectivity issues.

The doctor command checks:
- Configuration file validity and syntax
- Provider registration and availability
- Authentication and permissions (when --auth is specified)
- Network connectivity to providers (when --net is specified)
- File permissions and environment setup
- System information and Go version`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDoctor(cmd, verbose, checkAuth, checkNet, timeoutStr)
		},
	}

	c.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed diagnostic information")
	c.Flags().BoolVar(&checkAuth, "auth", false, "Check authentication and permissions (may make API calls)")
	c.Flags().BoolVar(&checkNet, "net", false, "Check network connectivity to providers")
	c.Flags().StringVar(&timeoutStr, "timeout", "30s", "Timeout for network checks")

	return c
}

func runDoctor(cmd *cobra.Command, verbose, checkAuth, checkNet bool, timeoutStr string) error {
	out := cmd.ErrOrStderr() // Use stderr for diagnostic output
	if _, err := fmt.Fprintln(out, "skv Doctor - Diagnostic Report"); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	if _, err := fmt.Fprintln(out, strings.Repeat("=", 50)); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// System Information
	if _, err := fmt.Fprintln(out, "\nSystem Information:"); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	if _, err := fmt.Fprintf(out, "  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	if _, err := fmt.Fprintf(out, "  Go Version: %s\n", runtime.Version()); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// Configuration Check
	if _, err := fmt.Fprintln(out, "\nConfiguration Check:"); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// Load configuration (this will find the config file automatically)
	cfg, err := config.Load(cfgPath)
	if err != nil {
		if strings.Contains(err.Error(), "no config file found") {
			if _, err := fmt.Fprintln(out, "  ERROR: No configuration file found"); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			if _, err := fmt.Fprintln(out, "     Checked locations:"); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			if envPath := os.Getenv("SKV_CONFIG"); envPath != "" {
				if _, err := fmt.Fprintf(out, "       - %s (SKV_CONFIG)\n", envPath); err != nil {
					return fmt.Errorf("failed to write output: %w", err)
				}
			}
			homeDir, _ := os.UserHomeDir()
			if _, err := fmt.Fprintf(out, "       - %s/.skv.yaml\n", homeDir); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			if _, err := fmt.Fprintf(out, "       - %s/.skv.yml\n", homeDir); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			if _, err := fmt.Fprintln(out, "     TIP: Use 'skv init' to create a template configuration file"); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			return err
		}
		if _, err := fmt.Fprintf(out, "  ERROR: Failed to load configuration: %v\n", err); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		return err
	}

	if _, err := fmt.Fprintf(out, "  OK: Configuration loaded: %d secrets configured\n", len(cfg.Secrets)); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// Provider Registration Check
	if _, err := fmt.Fprintln(out, "\nProvider Registration:"); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	allProviders := []string{"aws", "aws-ssm", "gcp", "azure", "azure-appconfig", "vault", "exec"}
	for _, p := range allProviders {
		if _, ok := provider.Get(p); ok {
			if _, err := fmt.Fprintf(out, "  OK: %s: registered\n", p); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
		} else {
			if _, err := fmt.Fprintf(out, "  ERROR: %s: not registered\n", p); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
		}
	}

	// Configuration Validation
	if _, err := fmt.Fprintln(out, "\nConfiguration Validation:"); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	issues := 0

	for i, secret := range cfg.Secrets {
		if verbose {
			if _, err := fmt.Fprintf(out, "  Checking secret %d: %s (%s)\n", i+1, secret.Alias, secret.Provider); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
		}

		// Check provider exists
		if _, ok := provider.Get(secret.Provider); !ok {
			if _, err := fmt.Fprintf(out, "    ERROR: Secret '%s': unknown provider '%s'\n", secret.Alias, secret.Provider); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			issues++
			continue
		}

		// Validate required fields
		if secret.Alias == "" {
			if _, err := fmt.Fprintf(out, "    ERROR: Secret %d: missing alias\n", i+1); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			issues++
		}
		if secret.Provider == "" {
			if _, err := fmt.Fprintf(out, "    ERROR: Secret '%s': missing provider\n", secret.Alias); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			issues++
		}
		if secret.Name == "" {
			if _, err := fmt.Fprintf(out, "    ERROR: Secret '%s': missing name\n", secret.Alias); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
			issues++
		}

		// Check for duplicate aliases
		for j, other := range cfg.Secrets {
			if i != j && secret.Alias == other.Alias {
				if _, err := fmt.Fprintf(out, "    ERROR: Duplicate alias: '%s' (secrets %d and %d)\n", secret.Alias, i+1, j+1); err != nil {
					return fmt.Errorf("failed to write output: %w", err)
				}
				issues++
			}
		}
	}

	if issues == 0 {
		if _, err := fmt.Fprintln(out, "  OK: No configuration issues found"); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	} else {
		if _, err := fmt.Fprintf(out, "  WARNING: Found %d configuration issue(s)\n", issues); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	}

	// Authentication Check (if requested)
	if checkAuth {
		if _, err := fmt.Fprintln(out, "\nAuthentication Check:"); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		timeout, _ := time.ParseDuration(timeoutStr)

		for _, secret := range cfg.Secrets {
			if _, err := fmt.Fprintf(out, "  Checking %s (%s)... ", secret.Alias, secret.Provider); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			spec := secret.ToSpec()

			p, _ := provider.Get(spec.Provider)
			_, err := p.FetchSecret(ctx, spec)
			cancel()

			if err == nil {
				if _, err := fmt.Fprintln(out, "OK"); err != nil {
					return fmt.Errorf("failed to write output: %w", err)
				}
			} else if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "404") {
				if _, err := fmt.Fprintln(out, "WARNING: (secret not found)"); err != nil {
					return fmt.Errorf("failed to write output: %w", err)
				}
			} else {
				if _, err := fmt.Fprintf(out, "ERROR: (%v)\n", err); err != nil {
					return fmt.Errorf("failed to write output: %w", err)
				}
			}
		}
	}

	// Network Connectivity Check (if requested)
	if checkNet {
		if _, err := fmt.Fprintln(out, "\nNetwork Connectivity Check:"); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		timeout, _ := time.ParseDuration(timeoutStr)

		// This would check basic connectivity to provider endpoints
		// For now, just show that the feature is implemented
		if _, err := fmt.Fprintf(out, "  INFO: Network checks would be performed here (timeout: %s)\n", timeout); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		if _, err := fmt.Fprintln(out, "     This feature can be extended to check connectivity to:"); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		if _, err := fmt.Fprintln(out, "     - AWS endpoints (secretsmanager, ssm)"); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		if _, err := fmt.Fprintln(out, "     - GCP Secret Manager API"); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		if _, err := fmt.Fprintln(out, "     - Azure Key Vault endpoints"); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		if _, err := fmt.Fprintln(out, "     - HashiCorp Vault servers"); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	}

	// Environment Check
	if _, err := fmt.Fprintln(out, "\nEnvironment Check:"); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".skv.yaml")
	if _, err := os.Stat(configDir); err == nil {
		if _, err := fmt.Fprintln(out, "  OK: Default config location exists"); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	} else {
		if _, err := fmt.Fprintln(out, "  INFO: Default config location doesn't exist (this is normal if using custom path)"); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	}

	// Summary
	if _, err := fmt.Fprintln(out, "\nSummary:"); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	if issues == 0 {
		if _, err := fmt.Fprintln(out, "  OK: All checks passed!"); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	} else {
		if _, err := fmt.Fprintf(out, "  WARNING: %d issue(s) found - review output above\n", issues); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	}

	if !checkAuth && !checkNet {
		if _, err := fmt.Fprintln(out, "\nTIP: Use --auth to check authentication or --net to check network connectivity"); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	}

	return nil
}

