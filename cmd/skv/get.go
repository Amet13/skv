package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"skv/internal/config"
	"skv/internal/provider"
)

func newGetCmd() *cobra.Command {
	var newline bool
	var raw = true
	var timeoutStr string
	var retries int
	var retryDelayStr string

	c := &cobra.Command{
		Use:   "get <alias>",
		Short: "Fetch a single secret and print it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			alias := args[0]
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return exitCodeError{code: 2, err: err}
			}

			s, ok := cfg.FindByAlias(alias)
			if !ok {
				return exitCodeError{code: 4, err: fmt.Errorf("alias not found: %s", alias)}
			}

			spec := s.ToSpec()
			p, ok := provider.Get(spec.Provider)
			if !ok {
				return exitCodeError{code: 3, err: fmt.Errorf("unknown provider: %s", spec.Provider)}
			}

			ctx := context.Background()
			if timeoutStr != "" {
				d, err := time.ParseDuration(timeoutStr)
				if err != nil {
					return exitCodeError{code: 2, err: fmt.Errorf("invalid --timeout: %w", err)}
				}
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, d)
				defer cancel()
			}

			delay := 500 * time.Millisecond
			if retryDelayStr != "" {
				if d, err := time.ParseDuration(retryDelayStr); err == nil {
					delay = d
				}
			}

			val, err := fetchWithRetry(ctx, p, spec, retries, delay)
			if err != nil {
				if errors.Is(err, provider.ErrNotFound) {
					return exitCodeError{code: 4, err: err}
				}
				return exitCodeError{code: 3, err: err}
			}

			// Apply transformation if configured
			transformedVal, err := s.TransformValue(val)
			if err != nil {
				return exitCodeError{code: 3, err: fmt.Errorf("transform error: %w", err)}
			}
			val = transformedVal

			out := cmd.OutOrStdout()
			if raw {
				if newline {
					if _, err := fmt.Fprintln(out, val); err != nil {
						return err
					}
				} else {
					if _, err := fmt.Fprint(out, val); err != nil {
						return err
					}
				}
				return nil
			}

			if newline {
				if _, err := fmt.Fprintln(out, val); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprint(out, val); err != nil {
					return err
				}
			}
			return nil
		},
	}

	c.Flags().BoolVar(&newline, "newline", false, "Append trailing newline")
	c.Flags().BoolVar(&raw, "raw", true, "Print raw secret value")
	c.Flags().StringVar(&timeoutStr, "timeout", "", "Timeout for fetching the secret (e.g., 5s, 30s)")
	c.Flags().IntVar(&retries, "retries", 0, "Number of retries on transient errors")
	c.Flags().StringVar(&retryDelayStr, "retry-delay", "500ms", "Delay between retries (e.g., 200ms, 1s)")
	return c
}

