package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"skv/internal/config"
	"skv/internal/provider"
)

func newExportCmd() *cobra.Command {
	var (
		all        bool
		secretsCSV string
		envFile    bool
	)

	c := &cobra.Command{
		Use:   "export",
		Short: "Export secrets as env lines or a .env file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return exitCodeError{code: 2, err: err}
			}
			requested := map[string]struct{}{}
			if all {
				for _, s := range cfg.Secrets {
					requested[s.Alias] = struct{}{}
				}
			}
			if secretsCSV != "" {
				for _, a := range splitCSV(secretsCSV) {
					requested[a] = struct{}{}
				}
			}
			if len(requested) == 0 {
				return exitCodeError{code: 2, err: fmt.Errorf("no secrets selected; use --all or --secrets")}
			}

			ctx := context.Background()
			kv := map[string]string{}
			for alias := range requested {
				s, ok := cfg.FindByAlias(alias)
				if !ok {
					return exitCodeError{code: 4, err: fmt.Errorf("alias not found: %s", alias)}
				}
				spec := s.ToSpec()
				p, ok := provider.Get(spec.Provider)
				if !ok {
					return exitCodeError{code: 3, err: fmt.Errorf("unknown provider: %s", spec.Provider)}
				}
				val, err := p.FetchSecret(ctx, spec)
				if err != nil {
					return exitCodeError{code: 3, err: fmt.Errorf("%s: %w", alias, err)}
				}
				kv[spec.EnvName] = val
			}
			out := cmd.OutOrStdout()
			// Stable order for deterministic output
			keys := make([]string, 0, len(kv))
			for k := range kv {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				if envFile {
					if _, err := fmt.Fprintf(out, "%s=%s\n", k, kv[k]); err != nil {
						return err
					}
				} else {
					if _, err := fmt.Fprintf(out, "export %s=%q\n", k, kv[k]); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}

	c.Flags().BoolVar(&all, "all", false, "Export all configured secrets")
	c.Flags().StringVar(&secretsCSV, "secrets", "", "Comma-separated list of aliases")
	c.Flags().BoolVar(&envFile, "env-file", false, "Write in .env format (no quoting)")
	return c
}

func splitCSV(s string) []string {
	var out []string
	field := ""
	for _, r := range s {
		if r == ',' {
			if field != "" {
				out = append(out, trim(field))
				field = ""
			}
			continue
		}
		field += string(r)
	}
	if t := trim(field); t != "" {
		out = append(out, t)
	}
	return out
}

func trim(s string) string {
	i, j := 0, len(s)
	for i < j && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	for j > i && (s[j-1] == ' ' || s[j-1] == '\t') {
		j--
	}
	return s[i:j]
}

