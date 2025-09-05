#!/usr/bin/env bash
# Comprehensive linting and formatting script for skv project

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

# Go linting and testing
echo "Running Go lint/tests..."
if ! command -v golangci-lint >/dev/null 2>&1; then
    echo "Installing golangci-lint locally..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${ROOT_DIR}/bin" latest
    export PATH="${ROOT_DIR}/bin:${PATH}"
fi

golangci-lint run
go test ./...

# Go code formatting
echo "Running gofmt (apply)..."
gofmt -s -w "$ROOT_DIR"

# Clean up trailing blank lines in Go files
echo "Trimming trailing blank lines in .go files..."
while IFS= read -r -d '' f; do
    awk 'BEGIN{n=0} {a[++n]=$0} END{while(n>0 && a[n]=="") n--; for(i=1;i<=n;i++) print a[i]; print ""}' "$f" >"$f.tmp" && mv "$f.tmp" "$f"
done < <(find "$ROOT_DIR" -type f -name '*.go' -print0)

# Shell script linting
echo "Running shellcheck..."
if command -v shellcheck >/dev/null 2>&1; then
    if [ -f "$ROOT_DIR/scripts/build.sh" ]; then
        shellcheck "$ROOT_DIR"/scripts/build.sh
    fi
    if [ -f "$ROOT_DIR/scripts/lint.sh" ]; then
        shellcheck "$ROOT_DIR"/scripts/lint.sh
    fi
fi

# Shell script formatting
echo "Running shfmt (apply)..."
if command -v shfmt >/dev/null 2>&1; then
    shfmt -i 4 -w "$ROOT_DIR"
fi

# Markdown linting and formatting
echo "Running markdownlint (apply)..."
if command -v markdownlint >/dev/null 2>&1; then
    markdownlint --fix "$ROOT_DIR"/**/*.md || true
else
    if command -v npx >/dev/null 2>&1; then
        npx --yes markdownlint-cli --fix "$ROOT_DIR"/**/*.md || true
    fi
fi

# YAML formatting
echo "Running yamlfmt (apply) if available..."
if command -v yamlfmt >/dev/null 2>&1; then
    yamlfmt "$ROOT_DIR"
fi

# GitHub Actions workflow linting
echo "Running actionlint (GitHub Actions workflows)..."
if command -v actionlint >/dev/null 2>&1; then
    actionlint
fi

echo "All linters completed."
