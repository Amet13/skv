# Cross-Platform Development

This document outlines how to develop `skv` on different operating systems.

## Prerequisites

### All Platforms

- **Go 1.25+**: Download from [golang.org](https://golang.org/dl/)
- **Git**: For version control and dependency management
- **Make**: For running build tasks (optional, see alternatives below)

### Platform-Specific Tools

#### Linux/macOS

- **bash**: Usually pre-installed
- **curl**: For downloading tools (usually pre-installed)

#### Windows

- **PowerShell 5.1+**: Usually pre-installed on Windows 10/11
- **Git Bash**: Recommended for bash script compatibility
- **Windows Subsystem for Linux (WSL)**: Optional but recommended

## Development Environment Setup

### Linux/macOS

```bash
# Clone the repository
git clone https://github.com/Amet13/skv.git
cd skv

# Build the project
make build

# Run tests
make test

# Run linting
make lint
```

### Windows

#### Option 1: PowerShell (Recommended)

```powershell
# Clone the repository
git clone https://github.com/Amet13/skv.git
cd skv

# Build using PowerShell script
powershell -ExecutionPolicy Bypass -File ./scripts/build.ps1 host

# Run tests
go test ./...

# Build for all platforms
powershell -ExecutionPolicy Bypass -File ./scripts/build.ps1 all
```

#### Option 2: Command Prompt

```cmd
# Clone the repository
git clone https://github.com/Amet13/skv.git
cd skv

# Build using batch script
scripts\build.cmd host

# Run tests
go test ./...

# Build for all platforms
scripts\build.cmd all
```

#### Option 3: Git Bash/WSL

```bash
# Same as Linux/macOS instructions
make build
make test
make lint
```

## Build Scripts

The project includes multiple build scripts for cross-platform compatibility:

- **`scripts/build.sh`**: Bash script (Linux/macOS/WSL/Git Bash)
- **`scripts/build.ps1`**: PowerShell script (Windows)
- **`scripts/build.cmd`**: Batch script (Windows Command Prompt)

All scripts support the same arguments:

- `host`: Build for current platform (default)
- `all`: Build for multiple platforms (darwin/arm64, darwin/amd64, linux/amd64, windows/amd64)

## IDE Configuration

### Visual Studio Code

Recommended extensions:

- **Go**: Official Go extension
- **EditorConfig**: Respects .editorconfig settings
- **GitLens**: Enhanced Git capabilities

Settings are automatically configured via `.editorconfig`.

### GoLand/IntelliJ IDEA

- Import as Go module
- Enable Go modules support
- Configure code style to match .editorconfig

### Vim/Neovim

Install Go language server and configure:

```vim
" Example vim-go configuration
let g:go_fmt_command = "gofmt"
let g:go_fmt_options = "-s"
```

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

### Integration Tests

```bash
# Run integration tests (requires longer timeout)
go test -timeout 30s ./...

# Skip integration tests
go test -short ./...
```

## Linting and Formatting

### Automated (via Make)

```bash
make lint  # Run all linters
make fmt   # Format all code
```

### Manual

```bash
# Go formatting
gofmt -s -w .

# Go linting (requires golangci-lint)
golangci-lint run

# Shell script linting (requires shellcheck)
shellcheck scripts/*.sh

# Markdown linting (requires markdownlint or npx)
markdownlint **/*.md
```

## Common Issues

### Windows Path Issues

If you encounter path-related issues on Windows:

1. Use forward slashes in Go code: `filepath.Join()` handles conversion
2. Use `os.PathSeparator` for dynamic path separators
3. Test with both PowerShell and Command Prompt

### Line Ending Issues

The project uses LF line endings. Configure Git:

```bash
# Linux/macOS
git config core.autocrlf input

# Windows
git config core.autocrlf true
```

### Permission Issues

On Unix-like systems, ensure scripts are executable:

```bash
chmod +x scripts/*.sh
```

### Go Module Issues

If you encounter module-related issues:

```bash
# Clean module cache
go clean -modcache

# Verify dependencies
go mod verify

# Update dependencies
go mod tidy
```

## Environment Variables

### Build Configuration

- `VERSION`: Override version string
- `PGO_PROFILE`: Path to PGO profile for optimized builds
- `GOTOOLCHAIN`: Go toolchain selection (set to `auto` by default)

### Development

- `GOPATH`: Go workspace (not required with modules)
- `GOPROXY`: Go module proxy (default: `https://proxy.golang.org`)
- `GOSUMDB`: Go checksum database (default: `sum.golang.org`)

## Continuous Integration

The project uses GitHub Actions for CI/CD:

- **Linux**: Primary testing platform
- **macOS**: Cross-compilation verification
- **Windows**: Cross-compilation verification

Local testing should cover the same scenarios as CI.

## Troubleshooting

### Build Failures

1. Ensure Go version is 1.23+: `go version`
2. Verify module integrity: `go mod verify`
3. Clean build cache: `go clean -cache`
4. Check for platform-specific issues in build logs

### Test Failures

1. Run tests with verbose output: `go test -v ./...`
2. Check for race conditions: `go test -race ./...`
3. Verify test dependencies are available
4. Check for platform-specific test issues

### Linting Issues

1. Update linting tools to latest versions
2. Check .golangci.yml configuration
3. Verify file encoding (should be UTF-8)
4. Check line endings (should be LF)

For additional help, see the [troubleshooting guide](../troubleshooting.md) or open an issue on GitHub.
