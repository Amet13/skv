package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"skv/internal/config"
)

func newListCmd() *cobra.Command {
	var verbose bool

	c := &cobra.Command{
		Use:   "list",
		Short: "List configured secret aliases",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return exitCodeError{code: 2, err: err}
			}
			out := cmd.OutOrStdout()
			for _, s := range cfg.Secrets {
				if verbose {
					if _, err := fmt.Fprintf(out, "%s\t%s\t%s\n", s.Alias, s.Provider, s.Env); err != nil {
						return err
					}
				} else {
					if _, err := fmt.Fprintln(out, s.Alias); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}

	c.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show provider and env mapping")
	return c
}

