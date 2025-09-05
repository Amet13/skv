#!/usr/bin/env bash
# Build script for skv binary - supports local and cross-platform builds

set -euo pipefail

# Setup build environment
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DIST_DIR="$REPO_ROOT/dist"
PGO_PROFILE="${PGO_PROFILE:-}"  # Profile-Guided Optimization profile path
VERSION_OVERRIDE="${VERSION:-}" # Override version (useful for CI)
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

    # Build configuration
    local build_flags=("-trimpath") # Remove file system paths from binary
    local ldflags=("-s -w")         # Strip debug info and symbol table

    # Version information
    local version
    if [[ -n "$VERSION_OVERRIDE" ]]; then
        version="$VERSION_OVERRIDE"
    else
        if version=$(git describe --tags --dirty --always 2>/dev/null); then :; else version="dev"; fi
    fi

    # Git commit and build timestamp
    local commit
    commit=$(git rev-parse --short HEAD 2>/dev/null || echo "")
    local date
    date=$(date -u +%Y-%m-%dT%H:%M:%SZ)

    # Inject build-time variables
    ldflags+=("-X skv/internal/version.Version=${version}")
    ldflags+=("-X skv/internal/version.Commit=${commit}")
    ldflags+=("-X skv/internal/version.Date=${date}")

    # Enable PGO if profile is available
    if [[ -n "$PGO_PROFILE" && -f "$PGO_PROFILE" ]]; then
        build_flags+=("-pgo=$PGO_PROFILE")
    fi

    # Execute build
    (cd "$REPO_ROOT" &&
        GOOS="$goos" GOARCH="$goarch" \
            go build "${build_flags[@]}" -ldflags "${ldflags[*]}" -o "$DIST_DIR/skv_${goos}_${goarch}" ./cmd/skv)
}

# Build execution based on mode
case "$MODE" in
host)
    # Build for current platform
    build_target "$(go env GOOS)" "$(go env GOARCH)"
    ;;
all)
    # Cross-platform build matrix
    build_target darwin arm64
    build_target darwin amd64
    build_target linux amd64
    ;;
*)
    usage
    exit 2
    ;;
esac

# Display build artifacts
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
