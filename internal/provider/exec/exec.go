// Package execprovider implements an external command provider.
package execprovider

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"skv/internal/provider"
)

// Exec provider runs an external command to fetch a secret value.
// Extras:
// - cmd: absolute or resolvable command to execute (required)
// - args: optional space-separated arguments (no shell parsing of quotes)
// - cwd: optional working directory
// - env: optional CSV of k=v entries to extend env
// - trim: optional "true" to trim whitespace from stdout
// The secret name (spec.Name) is appended as the last argument.
type execProvider struct{}

// New returns a new exec-based provider.
func New() provider.Provider { return &execProvider{} }

func (e *execProvider) FetchSecret(ctx context.Context, spec provider.SecretSpec) (string, error) {
	command := strings.TrimSpace(spec.Extras["cmd"])
	if command == "" {
		// Default to using spec.Name as the command when cmd is not provided
		command = strings.TrimSpace(spec.Name)
	}
	if command == "" {
		return "", fmt.Errorf("exec provider: missing command; set extras.cmd or name for %s", spec.Alias)
	}
	var args []string
	if a := strings.TrimSpace(spec.Extras["args"]); a != "" {
		// naive split by spaces; callers should avoid complex quoting
		args = append(args, strings.Fields(a)...)
	}
	// Append spec.Name as the last argument only when it is not the command itself
	if spec.Name != "" && command != spec.Name {
		args = append(args, spec.Name)
	}

	// #nosec G204 — executing user-provided command is the purpose of this provider
	c := exec.CommandContext(ctx, command, args...)
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	if cwd := strings.TrimSpace(spec.Extras["cwd"]); cwd != "" {
		c.Dir = cwd
	}
	if ev := strings.TrimSpace(spec.Extras["env"]); ev != "" {
		env := os.Environ()
		for _, part := range strings.Split(ev, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				env = append(env, kv[0]+"="+kv[1])
			}
		}
		c.Env = env
	}
	if err := c.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("exec provider: %v: %s", err, strings.TrimSpace(stderr.String()))
		}
		return "", fmt.Errorf("exec provider: %w", err)
	}
	out := stdout.String()
	if strings.EqualFold(spec.Extras["trim"], "true") {
		out = strings.TrimSpace(out)
	}
	return out, nil
}

