# Contributing to skv

Thank you for your interest in contributing to `skv`! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Process](#contributing-process)
- [Code Guidelines](#code-guidelines)
- [Testing](#testing)
- [Documentation](#documentation)
- [Security](#security)

## Code of Conduct

This project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

### Types of Contributions

We welcome several types of contributions:

- **Bug Reports**: Found a bug? Please report it!
- **Feature Requests**: Have an idea for a new feature?
- **Code Contributions**: Bug fixes, new features, or improvements
- **Documentation**: Improvements to docs, examples, or tutorials
- **Provider Development**: New secret provider implementations
- **Testing**: Additional test cases or test infrastructure improvements

### Before You Start

1. Check existing [issues](https://github.com/Amet13/skv/issues) and [pull requests](https://github.com/Amet13/skv/pulls)
2. For large changes, open an issue first to discuss the approach
3. Read the [architecture documentation](docs/architecture.md) to understand the codebase

## Development Setup

### Prerequisites

- Go 1.25+
- Git
- Make
- Docker (optional, for container testing)

### Local Setup

```bash
# Clone the repository
git clone https://github.com/Amet13/skv.git
cd skv

# Install dependencies
go mod download

# Build the project
make build

# Run tests
make test

# Run linters
make lint

# Run full CI pipeline
make ci
```

### Development Commands

```bash
# Format code
make fmt

# Run tests with coverage
make cover

# Build for all platforms
make build-all

# Generate shell completions
mkdir -p completions
./dist/skv_linux_amd64 completion bash > completions/skv.bash
./dist/skv_linux_amd64 completion zsh > completions/_skv
./dist/skv_linux_amd64 completion fish > completions/skv.fish
```

## Contributing Process

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/skv.git
cd skv

# Add upstream remote
git remote add upstream https://github.com/Amet13/skv.git
```

### 2. Create a Branch

```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Or a bugfix branch
git checkout -b fix/issue-number-description
```

### 3. Make Changes

- Write clear, concise commit messages
- Follow the [code guidelines](#code-guidelines)
- Add tests for new functionality
- Update documentation as needed

### 4. Test Your Changes

```bash
# Run the full test suite
make ci

# Test specific functionality
go test ./internal/provider/yourprovider -v

# Test with real providers (if applicable)
# Create a test config and run:
./dist/skv validate --check-secrets
./dist/skv run --dry-run --all -- echo "test"
```

### 5. Submit a Pull Request

1. Push your branch to your fork
2. Open a pull request against the main branch
3. Fill out the pull request template completely
4. Wait for review and address feedback

### Pull Request Guidelines

- **Title**: Use a descriptive title that summarizes the change
- **Description**: Explain what the PR does and why
- **Testing**: Describe how you tested the changes
- **Breaking Changes**: Clearly mark any breaking changes
- **Documentation**: Update docs if needed

## Code Guidelines

### Go Code Standards

- Follow standard Go conventions and idioms
- Use `gofmt` and `golangci-lint` (run `make lint`)
- Write clear, self-documenting code
- Add comments for exported functions and complex logic
- Handle errors explicitly and meaningfully

### Code Structure

- Keep functions focused and small
- Use dependency injection for testability
- Follow the existing package structure
- Maintain backward compatibility when possible

### Provider Development

When adding a new provider:

1. Create a new package under `internal/provider/yourprovider/`
2. Implement the `Provider` interface
3. Add comprehensive tests
4. Update documentation in `docs/providers.md`
5. Add example configuration to templates
6. Register the provider in `cmd/skv/main.go`

See [developing-providers.md](docs/dev/developing-providers.md) for detailed guidance.

### Example Provider Structure

```go
// internal/provider/yourprovider/yourprovider.go
package yourprovider

import (
    "context"
    "skv/internal/provider"
)

type yourProvider struct{}

func New() provider.Provider {
    return &yourProvider{}
}

func (p *yourProvider) FetchSecret(ctx context.Context, spec provider.SecretSpec) (string, error) {
    // Implementation
}
```

## Testing

### Test Requirements

- All new code must have tests
- Maintain or improve test coverage (aim for >70%)
- Use table-driven tests where appropriate
- Mock external dependencies
- Test both success and error paths

### Test Structure

```go
func TestYourFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
        wantErr  bool
    }{
        {
            name:     "success case",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        {
            name:    "error case",
            input:   invalidInput,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := YourFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("YourFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && result != tt.expected {
                t.Errorf("YourFunction() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
make cover

# Run tests for specific package
go test ./internal/provider/aws -v

# Run specific test
go test ./internal/provider/aws -run TestAWSProvider_FetchSecret
```

## Documentation

### Documentation Requirements

- Update relevant documentation for any changes
- Add examples for new features
- Update CLI help text if commands change
- Keep README.md current

### Documentation Structure

- `README.md`: Main project documentation
- `docs/`: Detailed documentation
- Code comments: Explain complex logic
- Examples: Show real-world usage

### Writing Guidelines

- Use clear, concise language
- Provide examples where helpful
- Keep documentation up-to-date with code changes
- Use proper markdown formatting

## Security

### Security Considerations

- Never log or expose secret values
- Follow secure coding practices
- Validate all inputs
- Use secure defaults
- Consider timing attacks for sensitive operations

### Security Review

All security-related changes undergo additional review:

- Provider implementations
- Authentication mechanisms
- Configuration handling
- Error messages that might leak information

### Reporting Security Issues

Please report security vulnerabilities privately:

- Use GitHub Security Advisories
- Email maintainers directly
- Do not open public issues for security problems

## Release Process

### Version Numbering

We use [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

1. Update `CHANGELOG.md`
2. Update version in documentation
3. Ensure all tests pass
4. Create and push tag
5. GitHub Actions handles the release

## Getting Help

### Community Support

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: Questions and community support
- **Documentation**: Check existing docs first

### Maintainer Contact

- Open an issue for project-related questions
- Use GitHub Discussions for general questions
- Email for security issues only

## Recognition

Contributors are recognized in several ways:

- Listed in `CHANGELOG.md` for their contributions
- GitHub contributor statistics
- Special recognition for significant contributions

Thank you for contributing to `skv`! Your efforts help make secret management better for everyone.
