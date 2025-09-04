package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"skv/internal/config"
	"skv/internal/provider"
)

func newExportCmd() *cobra.Command {
	var (
		all         bool
		secretsCSV  string
		allExcept   string
		envFile     bool
		format      string
		output      string
		noSort      bool
		retries     int
		retryDelay  string
		concurrency int
	)

	c := &cobra.Command{
		Use:   "export",
		Short: "Export secrets as env lines or a .env file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return exitCodeError{code: 2, err: err}
			}
			excluded := map[string]struct{}{}
			for _, a := range splitCSV(allExcept) {
				excluded[a] = struct{}{}
			}
			requested := map[string]struct{}{}
			if all {
				for _, s := range cfg.Secrets {
					if _, skip := excluded[s.Alias]; !skip {
						requested[s.Alias] = struct{}{}
					}
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
			if concurrency <= 0 {
				concurrency = 4
			}
			var wg sync.WaitGroup
			sem := make(chan struct{}, concurrency)
			var mu sync.Mutex
			var firstErr error
			for alias := range requested {
				alias := alias
				wg.Add(1)
				sem <- struct{}{}
				go func() {
					defer wg.Done()
					defer func() { <-sem }()
					if firstErr != nil {
						return
					}
					s, ok := cfg.FindByAlias(alias)
					if !ok {
						mu.Lock()
						if firstErr == nil {
							firstErr = exitCodeError{code: 4, err: fmt.Errorf("alias not found: %s", alias)}
						}
						mu.Unlock()
						return
					}
					spec := s.ToSpec()
					p, ok := provider.Get(spec.Provider)
					if !ok {
						mu.Lock()
						if firstErr == nil {
							firstErr = exitCodeError{code: 3, err: fmt.Errorf("unknown provider: %s", spec.Provider)}
						}
						mu.Unlock()
						return
					}
					d := 500 * time.Millisecond
					if retryDelay != "" {
						if dd, err := time.ParseDuration(retryDelay); err == nil {
							d = dd
						}
					}
					val, err := fetchWithRetry(ctx, p, spec, retries, d)
					if err != nil {
						mu.Lock()
						if firstErr == nil {
							firstErr = exitCodeError{code: 3, err: fmt.Errorf("%s: %w", alias, err)}
						}
						mu.Unlock()
						return
					}
					mu.Lock()
					kv[spec.EnvName] = val
					mu.Unlock()
				}()
			}
			wg.Wait()
			if firstErr != nil {
				return firstErr
			}

			out := cmd.OutOrStdout()
			if output != "" {
				// #nosec G304: output is a user-provided path by design (flag)
				f, err := os.Create(output)
				if err != nil {
					return err
				}
				defer func() { _ = f.Close() }()
				out = f
			}
			// Stable order for deterministic output
			keys := make([]string, 0, len(kv))
			for k := range kv {
				keys = append(keys, k)
			}
			if !noSort {
				sort.Strings(keys)
			}
			if envFile && format == "" {
				format = "env"
			}
			switch format {
			case "", "shell":
				for _, k := range keys {
					if _, err := fmt.Fprintf(out, "export %s=\"%s\"\n", k, kv[k]); err != nil {
						return err
					}
				}
			case "env":
				for _, k := range keys {
					if _, err := fmt.Fprintf(out, "%s=%s\n", k, kv[k]); err != nil {
						return err
					}
				}
			case "json":
				b, _ := json.MarshalIndent(kv, "", "  ")
				if _, err := out.Write(b); err != nil {
					return err
				}
				if _, err := fmt.Fprintln(out); err != nil {
					return err
				}
			case "yaml", "yml":
				b, _ := yaml.Marshal(kv)
				if _, err := out.Write(b); err != nil {
					return err
				}
			default:
				return exitCodeError{code: 2, err: fmt.Errorf("unsupported format: %s", format)}
			}
			return nil
		},
	}

	c.Flags().BoolVar(&all, "all", false, "Export all configured secrets")
	c.Flags().StringVar(&secretsCSV, "secrets", "", "Comma-separated list of aliases")
	c.Flags().StringVar(&allExcept, "all-except", "", "Comma-separated aliases to exclude when using --all")
	c.Flags().BoolVar(&envFile, "env-file", false, "Write in .env format (no quoting)")
	c.Flags().StringVar(&format, "format", "", "Output format: shell|env|json|yaml")
	c.Flags().StringVar(&output, "output", "", "Output file path (stdout if empty)")
	c.Flags().BoolVar(&noSort, "no-sort", false, "Do not sort environment variable names")
	c.Flags().IntVar(&retries, "retries", 0, "Number of retries on transient errors")
	c.Flags().StringVar(&retryDelay, "retry-delay", "500ms", "Delay between retries (e.g., 200ms, 1s)")
	c.Flags().IntVar(&concurrency, "concurrency", 4, "Number of concurrent provider calls")
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

