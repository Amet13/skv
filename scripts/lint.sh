#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

echo "Running Go lint/tests..."
if ! command -v golangci-lint >/dev/null 2>&1; then
    echo "Installing golangci-lint locally..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${ROOT_DIR}/bin" latest
    export PATH="${ROOT_DIR}/bin:${PATH}"
fi
golangci-lint run
go test ./...

echo "Running gofmt (apply)..."
gofmt -s -w "$ROOT_DIR"

echo "Trimming trailing blank lines in .go files..."
while IFS= read -r -d '' f; do
    awk 'BEGIN{n=0} {a[++n]=$0} END{while(n>0 && a[n]=="") n--; for(i=1;i<=n;i++) print a[i]; print ""}' "$f" >"$f.tmp" && mv "$f.tmp" "$f"
done < <(find "$ROOT_DIR" -type f -name '*.go' -print0)

echo "Running shellcheck..."
if command -v shellcheck >/dev/null 2>&1; then
    shellcheck "$ROOT_DIR"/build.sh
fi

echo "Running shfmt (apply)..."
if command -v shfmt >/dev/null 2>&1; then
    shfmt -i 4 -w "$ROOT_DIR"
fi

echo "Running markdownlint (apply)..."
if command -v markdownlint >/dev/null 2>&1; then
    markdownlint --fix "$ROOT_DIR"/**/*.md || true
else
    if command -v npx >/dev/null 2>&1; then
        npx --yes markdownlint-cli --fix "$ROOT_DIR"/**/*.md || true
    fi
fi

echo "Running yamllint..."
if command -v yamllint >/dev/null 2>&1; then
    yamllint -c "$ROOT_DIR/.yamllint.yaml" "$ROOT_DIR"
fi

echo "Running yamlfmt (apply) if available..."
if command -v yamlfmt >/dev/null 2>&1; then
    yamlfmt -w "$ROOT_DIR"
fi

echo "Running actionlint (GitHub Actions workflows)..."
if command -v actionlint >/dev/null 2>&1; then
    actionlint
fi

echo "All linters completed."
