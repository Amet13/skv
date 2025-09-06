package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"skv/internal/config"
	"skv/internal/provider"
)

func newWatchCmd() *cobra.Command {
	var (
		secretsCSV   string
		secretsList  []string
		all          bool
		allExceptCSV string
		interval     string
		command      string
		onChangeOnly bool
		timeoutStr   string
	)

	c := &cobra.Command{
		Use:   "watch [flags] -- <command>",
		Short: "Watch secrets for changes and execute command",
		Long: `Watch configured secrets for changes and execute a command when they change.

The watch command periodically checks secrets and runs the specified command
when changes are detected. This is useful for restarting services or updating
configurations when secrets change.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			command = strings.Join(args, " ")
			return runWatch(secretsCSV, secretsList, all, allExceptCSV, interval, command, onChangeOnly, timeoutStr)
		},
	}

	c.Flags().StringVar(&secretsCSV, "secrets", "", "Comma-separated list of aliases to watch")
	c.Flags().StringSliceVarP(&secretsList, "secret", "s", nil, "Secret alias to watch (repeatable)")
	c.Flags().BoolVar(&all, "all", false, "Watch all configured secrets")
	c.Flags().StringVar(&allExceptCSV, "all-except", "", "Comma-separated aliases to exclude when using --all")
	c.Flags().StringVar(&interval, "interval", "30s", "Check interval (e.g., 30s, 5m, 1h)")
	c.Flags().BoolVar(&onChangeOnly, "on-change-only", false, "Only execute command when secrets change")
	c.Flags().StringVar(&timeoutStr, "timeout", "", "Timeout for watch command (e.g., 30s, 5m, 1h)")

	return c
}

func runWatch(secretsCSV string, secretsList []string, all bool, allExceptCSV, intervalStr, command string, onChangeOnly bool, timeoutStr string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return exitCodeError{code: 2, err: err}
	}

	// Parse interval
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return exitCodeError{code: 2, err: fmt.Errorf("invalid interval: %w", err)}
	}

	// Determine which secrets to watch
	excluded := map[string]struct{}{}
	for _, a := range strings.Split(allExceptCSV, ",") {
		if t := strings.TrimSpace(a); t != "" {
			excluded[t] = struct{}{}
		}
	}

	watchList := map[string]struct{}{}
	if all {
		for _, s := range cfg.Secrets {
			if _, skip := excluded[s.Alias]; !skip {
				watchList[s.Alias] = struct{}{}
			}
		}
	}
	if secretsCSV != "" {
		for _, a := range strings.Split(secretsCSV, ",") {
			a = strings.TrimSpace(a)
			if a != "" {
				watchList[a] = struct{}{}
			}
		}
	}
	for _, a := range secretsList {
		watchList[a] = struct{}{}
	}

	if len(watchList) == 0 {
		return exitCodeError{code: 2, err: fmt.Errorf("no secrets selected; use --all or --secrets/-s")}
	}

	fmt.Printf("Watching %d secret(s) with %v interval\n", len(watchList), interval)
	fmt.Printf("Command: %s\n", command)

	// Initial execution
	if !onChangeOnly {
		fmt.Println("\nInitial execution...")
		if err := executeCommand(command); err != nil {
			fmt.Printf("ERROR: Initial command failed: %v\n", err)
		}
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Set up timeout if specified
	var timeoutChan <-chan time.Time
	if timeoutStr != "" {
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return exitCodeError{code: 2, err: fmt.Errorf("invalid --timeout: %w", err)}
		}
		timeoutChan = time.After(timeout)
		fmt.Printf("Watch will timeout after %v\n", timeout)
	}

	// Track secret values
	lastValues := make(map[string]string)

	// Watch loop
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Printf("\nStarting watch loop (Ctrl+C to stop)...\n")

	for {
		select {
		case <-sigChan:
			fmt.Println("\nWatch stopped by user")
			return nil
		case <-ticker.C:
			if err := checkAndExecute(cfg, watchList, lastValues, command, onChangeOnly); err != nil {
				fmt.Printf("ERROR: Watch error: %v\n", err)
			}
		}
		// Check timeout outside select if specified
		if timeoutChan != nil {
			select {
			case <-timeoutChan:
				fmt.Println("\nWatch timed out")
				return nil
			default:
			}
		}
	}
}

func checkAndExecute(cfg *config.Config, watchList map[string]struct{}, lastValues map[string]string, command string, onChangeOnly bool) error {
	ctx := context.Background()
	changed := false

	for alias := range watchList {
		s, ok := cfg.FindByAlias(alias)
		if !ok {
			return fmt.Errorf("secret '%s' not found in configuration", alias)
		}

		spec := s.ToSpec()
		p, ok := provider.Get(spec.Provider)
		if !ok {
			return fmt.Errorf("unknown provider '%s' for secret '%s'", spec.Provider, alias)
		}

		value, err := p.FetchSecret(ctx, spec)
		if err != nil {
			return fmt.Errorf("failed to fetch secret '%s': %w", alias, err)
		}

		lastValue, exists := lastValues[alias]
		if !exists || lastValue != value {
			if exists {
				fmt.Printf("INFO: Secret '%s' changed\n", alias)
			}
			lastValues[alias] = value
			changed = true
		}
	}

	if changed || !onChangeOnly {
		fmt.Printf("Executing command...\n")
		if err := executeCommand(command); err != nil {
			return fmt.Errorf("command execution failed: %w", err)
		}
	}

	return nil
}

func executeCommand(command string) error {
	// Parse command (simple splitting, no complex shell parsing)
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	cmd := exec.Command(parts[0], parts[1:]...) // #nosec G204 - command parts come from user input, which is expected for watch command
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

