package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"skv/internal/config"
	"skv/internal/provider"
)

func newRunCmd() *cobra.Command {
	var (
		secretsCSV   string
		secretsList  []string
		all          bool
		allExceptCSV string
		dryRun       bool
		strict       bool
		mask         bool
		timeoutStr   string
		concurrency  int
		retries      int
		retryDelay   string
		requireEnv   []string
		requireAlias []string
	)

	strict = true
	mask = true

	c := &cobra.Command{
		Use:   "run [flags] -- <command> [args...]",
		Short: "Inject secrets into env and execute command",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sepIdx := indexOf(args, "--")
			var cmdArgs []string
			if sepIdx >= 0 {
				cmdArgs = args[sepIdx+1:]
			} else {
				cmdArgs = args
			}
			if len(cmdArgs) == 0 {
				return exitCodeError{code: 2, err: errors.New("no command provided; use skv run -- <cmd> [args]")}
			}

			cfg, err := config.Load(cfgPath)
			if err != nil {
				return exitCodeError{code: 2, err: err}
			}

			excluded := map[string]struct{}{}
			for _, a := range strings.Split(allExceptCSV, ",") {
				if t := strings.TrimSpace(a); t != "" {
					excluded[t] = struct{}{}
				}
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
				for _, a := range strings.Split(secretsCSV, ",") {
					a = strings.TrimSpace(a)
					if a != "" {
						requested[a] = struct{}{}
					}
				}
			}
			for _, a := range secretsList {
				requested[a] = struct{}{}
			}

			if len(requested) == 0 {
				return exitCodeError{code: 2, err: errors.New("no secrets selected; use --all or --secrets/-s")}
			}

			timeout := time.Duration(0)
			if timeoutStr != "" {
				d, err := time.ParseDuration(timeoutStr)
				if err != nil {
					return exitCodeError{code: 2, err: fmt.Errorf("invalid --timeout: %w", err)}
				}
				timeout = d
			}
			ctx := context.Background()
			if timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}

			envAdditions := map[string]string{}
			if concurrency <= 0 {
				concurrency = 4
			}
			aliases := make([]string, 0, len(requested))
			for a := range requested {
				aliases = append(aliases, a)
			}
			sort.Strings(aliases)
			var wg sync.WaitGroup
			sem := make(chan struct{}, concurrency)
			var mu sync.Mutex
			var firstErr error
			for _, alias := range aliases {
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
						if strict {
							mu.Lock()
							if firstErr == nil {
								firstErr = exitCodeError{code: 4, err: fmt.Errorf("alias not found: %s", alias)}
							}
							mu.Unlock()
						}
						return
					}
					spec := s.ToSpec()
					p, ok := provider.Get(spec.Provider)
					if !ok {
						if strict {
							mu.Lock()
							if firstErr == nil {
								firstErr = exitCodeError{code: 3, err: fmt.Errorf("unknown provider: %s", spec.Provider)}
							}
							mu.Unlock()
						}
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
						if strict {
							mu.Lock()
							if firstErr == nil {
								if errors.Is(err, provider.ErrNotFound) {
									firstErr = exitCodeError{code: 4, err: fmt.Errorf("%s: %w", alias, err)}
								} else {
									firstErr = exitCodeError{code: 3, err: fmt.Errorf("%s: %w", alias, err)}
								}
							}
							mu.Unlock()
						}
						return
					}
					mu.Lock()
					envAdditions[spec.EnvName] = val
					mu.Unlock()
				}()
			}
			wg.Wait()
			if firstErr != nil {
				return firstErr
			}

			// require-env check
			for _, e := range requireEnv {
				if _, ok := envAdditions[e]; !ok {
					return exitCodeError{code: 4, err: fmt.Errorf("required env missing: %s", e)}
				}
			}
			// require-alias check
			for _, a := range requireAlias {
				if _, ok := requested[a]; !ok {
					return exitCodeError{code: 4, err: fmt.Errorf("required alias not selected: %s", a)}
				}
			}

			command := cmdArgs[0]
			commandArgs := cmdArgs[1:]

			if dryRun {
				errw := cmd.ErrOrStderr()
				if _, err := fmt.Fprintln(errw, "[dry-run] would execute:"); err != nil {
					return err
				}
				if _, err := fmt.Fprintf(errw, "  %s %s\n", command, strings.Join(commandArgs, " ")); err != nil {
					return err
				}
				if _, err := fmt.Fprintln(errw, "[dry-run] with environment additions:"); err != nil {
					return err
				}
				keys := make([]string, 0, len(envAdditions))
				for k := range envAdditions {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					shown := envAdditions[k]
					if mask {
						shown = maskValue(shown)
					}
					if _, err := fmt.Fprintf(errw, "  %s=%s\n", k, shown); err != nil {
						return err
					}
				}
				return nil
			}

			// #nosec G204 — the command is intentionally user-provided
			cexec := exec.CommandContext(ctx, command, commandArgs...)
			cexec.Stdout = os.Stdout
			cexec.Stderr = os.Stderr
			cexec.Stdin = os.Stdin

			env := os.Environ()
			for k, v := range envAdditions {
				env = append(env, fmt.Sprintf("%s=%s", k, v))
			}
			cexec.Env = env

			if err := cexec.Run(); err != nil {
				var ee *exec.ExitError
				if errors.As(err, &ee) {
					if status, ok := exitStatusOf(ee); ok {
						return exitCodeError{code: status, err: fmt.Errorf("command failed: %w", err)}
					}
				}
				return exitCodeError{code: 5, err: err}
			}
			return nil
		},
	}

	c.Flags().StringVar(&secretsCSV, "secrets", "", "Comma-separated list of aliases")
	c.Flags().StringSliceVarP(&secretsList, "secret", "s", nil, "Secret alias (repeatable)")
	c.Flags().BoolVar(&all, "all", false, "Inject all configured secrets")
	c.Flags().StringVar(&allExceptCSV, "all-except", "", "Comma-separated aliases to exclude when using --all")
	c.Flags().BoolVar(&dryRun, "dry-run", false, "Print what would be executed and exit")
	c.Flags().BoolVar(&strict, "strict", true, "Fail if any requested secret cannot be fetched")
	c.Flags().BoolVar(&mask, "mask", true, "Mask secret values in logs and dry-run output")
	c.Flags().StringVar(&timeoutStr, "timeout", "", "Timeout for fetching secrets (e.g., 5s, 30s)")
	c.Flags().IntVar(&concurrency, "concurrency", 4, "Number of concurrent provider calls")
	c.Flags().IntVar(&retries, "retries", 0, "Number of retries on transient errors")
	c.Flags().StringVar(&retryDelay, "retry-delay", "500ms", "Delay between retries (e.g., 200ms, 1s)")
	c.Flags().StringSliceVar(&requireEnv, "require-env", nil, "Environment variable names that must be present after fetch")
	c.Flags().StringSliceVar(&requireAlias, "require-alias", nil, "Aliases that must be selected")
	return c
}

// Utility types and functions

type exitCodeError struct {
	code int
	err  error
}

func (e exitCodeError) Error() string { return e.err.Error() }

func (e exitCodeError) Unwrap() error { return e.err }

func indexOf(slice []string, val string) int {
	for i, s := range slice {
		if s == val {
			return i
		}
	}
	return -1
}

func maskValue(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}

// exitStatusOf tries to map an ExitError to a code; returns (code, true) if mapped.
func exitStatusOf(ee *exec.ExitError) (int, bool) {
	// Try to extract platform-specific exit status
	type interfaceWithExitStatus interface{ ExitStatus() int }
	if x, ok := any(ee.Sys()).(interfaceWithExitStatus); ok {
		return x.ExitStatus(), true
	}
	// Fallback: return generic exit code
	return 5, true
}

