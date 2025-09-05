package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"skv/internal/config"
	"skv/internal/provider"
)

func newHealthCmd() *cobra.Command {
	var (
		timeout    string
		secretName string
	)

	cmd := &cobra.Command{
		Use:   "health",
		Short: "Health check for secret providers",
		Long: `Perform health checks on configured secret providers to ensure
they are accessible and responding correctly. This is useful for
monitoring and alerting in production environments.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			// Parse timeout
			var timeoutDuration time.Duration
			if timeout != "" {
				timeoutDuration, err = time.ParseDuration(timeout)
				if err != nil {
					return fmt.Errorf("invalid timeout duration: %w", err)
				}
			} else {
				timeoutDuration = 10 * time.Second
			}

			ctx := context.Background()
			if timeoutDuration > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, timeoutDuration)
				defer cancel()
			}

			fmt.Printf("🏥 Health check starting (timeout: %v)\n", timeoutDuration)
			fmt.Printf("📊 Checking %d secret(s)\n\n", len(cfg.Secrets))

			var healthyCount, totalCount int
			var firstError error

			for _, secret := range cfg.Secrets {
				// Skip if specific secret requested and this isn't it
				if secretName != "" && secret.Alias != secretName {
					continue
				}

				totalCount++
				fmt.Printf("🔍 Checking %s (%s)... ", secret.Alias, secret.Provider)

				spec := secret.ToSpec()
				p, ok := provider.Get(spec.Provider)
				if !ok {
					fmt.Printf("❌ Provider not found\n")
					if firstError == nil {
						firstError = fmt.Errorf("provider %s not found", spec.Provider)
					}
					continue
				}

				start := time.Now()
				_, err := p.FetchSecret(ctx, spec)
				duration := time.Since(start)

				if err != nil {
					if err == provider.ErrNotFound {
						fmt.Printf("⚠️  Not found (%.2fs)\n", duration.Seconds())
					} else {
						fmt.Printf("❌ Error: %v (%.2fs)\n", err, duration.Seconds())
						if firstError == nil {
							firstError = err
						}
					}
				} else {
					fmt.Printf("OK (%.2fs)\n", duration.Seconds())
					healthyCount++
				}
			}

			fmt.Printf("\n📈 Health check summary:\n")
			fmt.Printf("  Healthy: %d/%d (%.1f%%)\n", healthyCount, totalCount, float64(healthyCount)/float64(totalCount)*100)

			if healthyCount == totalCount {
				fmt.Printf("🎉 All secrets are healthy!\n")
				return nil
			}

			fmt.Printf("⚠️  %d secret(s) have issues\n", totalCount-healthyCount)
			if firstError != nil {
				return fmt.Errorf("health check failed: %w", firstError)
			}
			return fmt.Errorf("health check failed: %d/%d secrets unhealthy", totalCount-healthyCount, totalCount)
		},
	}

	cmd.Flags().StringVar(&timeout, "timeout", "10s", "Timeout for each health check")
	cmd.Flags().StringVarP(&secretName, "secret", "s", "", "Check specific secret only")

	return cmd
}

