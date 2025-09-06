package version

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		commit    string
		date      string
		wantParts []string
	}{
		{
			name:      "dev version only",
			version:   "dev",
			commit:    "",
			date:      "",
			wantParts: []string{"dev"},
		},
		{
			name:      "version with commit and date",
			version:   "v1.2.3",
			commit:    "abc123",
			date:      "2023-01-01T00:00:00Z",
			wantParts: []string{"v1.2.3", "commit abc123", "built 2023-01-01T00:00:00Z"},
		},
		{
			name:      "version with commit only",
			version:   "v1.0.0",
			commit:    "def456",
			date:      "",
			wantParts: []string{"v1.0.0", "commit def456"},
		},
		{
			name:      "version with date only",
			version:   "v2.0.0",
			commit:    "",
			date:      "2023-12-31T23:59:59Z",
			wantParts: []string{"v2.0.0", "built 2023-12-31T23:59:59Z"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set package variables
			oldVersion, oldCommit, oldDate := Version, Commit, Date
			defer func() {
				Version, Commit, Date = oldVersion, oldCommit, oldDate
			}()

			Version = tt.version
			Commit = tt.commit
			Date = tt.date

			result := String()

			// Check that all expected parts are present
			for _, part := range tt.wantParts {
				if !strings.Contains(result, part) {
					t.Errorf("String() = %q, want to contain %q", result, part)
				}
			}

			// For dev version only case, ensure it's exactly "dev"
			if tt.name == "dev version only" && result != "dev" {
				t.Errorf("String() = %q, want exactly %q", result, "dev")
			}
		})
	}
}

func TestStringFormat(t *testing.T) {
	oldVersion, oldCommit, oldDate := Version, Commit, Date
	defer func() {
		Version, Commit, Date = oldVersion, oldCommit, oldDate
	}()

	Version = "v1.0.0"
	Commit = "abc123"
	Date = "2023-01-01T00:00:00Z"

	result := String()
	expected := "v1.0.0 (commit abc123, built 2023-01-01T00:00:00Z)"

	if result != expected {
		t.Errorf("String() = %q, want %q", result, expected)
	}
}
