package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"skv/internal/config"
	"skv/internal/provider"
)

func newGetCmd() *cobra.Command {
	var newline bool
	var raw = true

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
			val, err := p.FetchSecret(ctx, spec)
			if err != nil {
				if errors.Is(err, provider.ErrNotFound) {
					return exitCodeError{code: 4, err: err}
				}
				return exitCodeError{code: 3, err: err}
			}

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
	return c
}

