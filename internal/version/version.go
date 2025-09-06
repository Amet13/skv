// Package version provides build-time version metadata.
package version

import "fmt"

var (
	// Version is the semantic version (set at build time).
	Version = "dev"
	// Commit is the git commit hash (set at build time).
	Commit = ""
	// Date is the build timestamp in RFC3339 (set at build time).
	Date = ""
)

// String formats version information for CLI output.
func String() string {
	if Commit == "" && Date == "" {
		return Version
	}
	return fmt.Sprintf("%s (commit %s, built %s)", Version, Commit, Date)
}

