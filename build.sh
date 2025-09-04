#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
DIST_DIR="$ROOT_DIR/dist"
PGO_PROFILE="${PGO_PROFILE:-}"
VERSION_OVERRIDE="${VERSION:-}"
mkdir -p "$DIST_DIR"

usage() {
    echo "Usage: $0 [host|all]" >&2
    echo "  host (default): build for current platform" >&2
    echo "  all: build for a small matrix (darwin/arm64, darwin/amd64, linux/amd64)" >&2
}

MODE="${1:-host}"

build_target() {
    local goos="$1" goarch="$2"
    echo "Building skv for ${goos}/${goarch}..."
    local build_flags=("-trimpath")
    local ldflags=("-s -w")
    local version
    if [[ -n "$VERSION_OVERRIDE" ]]; then
        version="$VERSION_OVERRIDE"
    else
        if version=$(git describe --tags --dirty --always 2>/dev/null); then :; else version="dev"; fi
    fi
    local commit
    commit=$(git rev-parse --short HEAD 2>/dev/null || echo "")
    local date
    date=$(date -u +%Y-%m-%dT%H:%M:%SZ)
    ldflags+=("-X skv/internal/version.Version=${version}")
    ldflags+=("-X skv/internal/version.Commit=${commit}")
    ldflags+=("-X skv/internal/version.Date=${date}")
    if [[ -n "$PGO_PROFILE" && -f "$PGO_PROFILE" ]]; then
        build_flags+=("-pgo=$PGO_PROFILE")
    fi
    (cd "$ROOT_DIR" &&
        GOOS="$goos" GOARCH="$goarch" \
            go mod tidy &&
        go build "${build_flags[@]}" -ldflags "${ldflags[*]}" -o "$DIST_DIR/skv_${goos}_${goarch}" ./cmd/skv)
}

case "$MODE" in
host)
    build_target "$(go env GOOS)" "$(go env GOARCH)"
    ;;
all)
    build_target darwin arm64
    build_target darwin amd64
    build_target linux amd64
    ;;
*)
    usage
    exit 2
    ;;
esac

echo "Artifacts in $DIST_DIR:"
if command -v find >/dev/null 2>&1; then
    # Print artifact file names (one per line) in a stable order
    find "$DIST_DIR" -maxdepth 1 -mindepth 1 -type f -exec basename {} \; | sort | cat
else
    # Fallback: print just the artifact paths
    for f in "$DIST_DIR"/*; do
        printf '%s\n' "$f"
    done | sort | cat
fi
