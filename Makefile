export GOTOOLCHAIN=auto
.PHONY: build build-all fmt lint test cover ci all clean release snapshot tools

build:
	./scripts/build.sh host

build-all:
	./scripts/build.sh all

fmt:
	./scripts/lint.sh || true # ensure tools present
	@echo "Formatting Go files..."
	gofmt -s -w .
	@echo "Formatting shell scripts..."
	@if command -v shfmt >/dev/null 2>&1; then shfmt -i 4 -w .; fi
	@echo "Formatting markdown..."
	@if command -v markdownlint >/dev/null 2>&1; then markdownlint --fix **/*.md || true; else if command -v npx >/dev/null 2>&1; then npx --yes markdownlint-cli --fix **/*.md || true; fi; fi

lint:
	./scripts/lint.sh

test:
	go test ./...

cover:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out | awk -v threshold=60.0 '/total:/ { gsub("%","",$$3); if ($$3+0 < threshold) { printf "Coverage %.2f%% below %.2f%%\n", $$3, threshold; exit 1 } }'

ci: lint cover build

all: build lint test

clean:
	rm -rf dist

release:
	goreleaser release --clean --skip=publish --snapshot

snapshot:
	GORELEASER_CURRENT_TAG=v0.0.0-dev goreleaser release --clean --skip=publish --snapshot
