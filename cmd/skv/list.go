package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"skv/internal/config"
)

func newListCmd() *cobra.Command {
	var verbose bool
	var format string

	c := &cobra.Command{
		Use:   "list",
		Short: "List configured secret aliases",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return exitCodeError{code: 2, err: err}
			}
			out := cmd.OutOrStdout()
			switch format {
			case "", "text":
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
			case "json":
				type item struct{ Alias, Provider, Env string }
				arr := make([]item, 0, len(cfg.Secrets))
				for _, s := range cfg.Secrets {
					arr = append(arr, item{Alias: s.Alias, Provider: s.Provider, Env: s.Env})
				}
				b, _ := json.MarshalIndent(arr, "", "  ")
				if _, err := out.Write(b); err != nil {
					return err
				}
				if _, err := fmt.Fprintln(out); err != nil {
					return err
				}
			case "yaml", "yml":
				arr := make([]map[string]string, 0, len(cfg.Secrets))
				for _, s := range cfg.Secrets {
					arr = append(arr, map[string]string{"alias": s.Alias, "provider": s.Provider, "env": s.Env})
				}
				b, _ := yaml.Marshal(arr)
				if _, err := out.Write(b); err != nil {
					return err
				}
			default:
				return exitCodeError{code: 2, err: fmt.Errorf("unsupported format: %s", format)}
			}
			return nil
		},
	}

	c.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show provider and env mapping")
	c.Flags().StringVar(&format, "format", "", "Output format: text|json|yaml")
	return c
}
