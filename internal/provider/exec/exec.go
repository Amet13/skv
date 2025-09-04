// Package execprovider implements an external command provider.
package execprovider

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"skv/internal/provider"
)

// Exec provider runs an external command to fetch a secret value.
// Extras:
// - cmd: absolute or resolvable command to execute (required)
// - args: optional space-separated arguments (no shell parsing of quotes)
// - trim: optional "true" to trim whitespace from stdout
// The secret name (spec.Name) is appended as the last argument.
type execProvider struct{}

// New returns a new exec-based provider.
func New() provider.Provider { return &execProvider{} }

func (e *execProvider) FetchSecret(ctx context.Context, spec provider.SecretSpec) (string, error) {
	command := strings.TrimSpace(spec.Extras["cmd"])
	if command == "" {
		return "", fmt.Errorf("exec provider: missing extras 'cmd' for %s", spec.Alias)
	}
	var args []string
	if a := strings.TrimSpace(spec.Extras["args"]); a != "" {
		// naive split by spaces; callers should avoid complex quoting
		args = append(args, strings.Fields(a)...)
	}
	args = append(args, spec.Name)

	// #nosec G204 — executing user-provided command is the purpose of this provider
	c := exec.CommandContext(ctx, command, args...)
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
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

