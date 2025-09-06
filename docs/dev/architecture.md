# Architecture

This document describes the overall architecture and design principles of `skv`.

## Overview

`skv` is designed as a simple, secure, and extensible CLI tool for cross-cloud secret management. The architecture follows these principles:

- **Security first**: Secrets are never written to disk, only kept in memory
- **Provider extensibility**: Easy to add new secret backends
- **Minimal interface**: Simple `Provider` interface with single method
- **Configuration-driven**: Behavior controlled via YAML configuration
- **Cross-platform**: Works on Linux, macOS, and Windows

## Project Structure

```
skv/
|-── cmd/skv/                    # CLI entry point and commands
|   |-── main.go                # Application bootstrap and provider registration
|   |-── get.go, run.go, etc.   # Individual command implementations
|   `-── *_test.go              # Command-level tests
|-── internal/
|   |-── config/                # Configuration loading and validation
|   |-── provider/              # Provider interface and implementations
|   |   |-── aws/               # AWS Secrets Manager & SSM
|   |   |-── gcp/               # Google Secret Manager
|   |   |-── azure/             # Azure Key Vault & App Config
|   |   |-── vault/             # HashiCorp Vault
|   |   |-── exec/              # External command execution
|   |   `-── mock/              # Testing mock provider
|   `-── version/               # Build-time version information
|-── docs/                      # User and developer documentation
`-── scripts/                   # Build and development scripts
```

## Core Components

### CLI Commands

Built using [Cobra](https://github.com/spf13/cobra), each command is implemented in a separate file:

- **`get`**: Fetch a single secret and print to stdout
- **`run`**: Inject secrets into environment and execute command
- **`export`**: Export secrets as environment variables or .env file
- **`list`**: List configured secret aliases
- **`validate`**: Validate configuration syntax and connectivity
- **`health`**: Check provider health and connectivity
- **`init`**: Generate configuration template
- **`version`**: Show version information
- **`completion`**: Generate shell completion scripts

### Configuration System

The configuration system (`internal/config/`) handles:

- **YAML parsing**: Load and validate configuration files
- **Environment interpolation**: Support `{{ ENV_VAR }}` syntax
- **Schema validation**: Ensure required fields are present
- **Default merging**: Apply global defaults to individual secrets
- **File discovery**: Find config in standard locations

Configuration structure:

```yaml
defaults: # Global defaults merged into each secret
  region: us-east-1
  extras:
    version_stage: AWSCURRENT

secrets: # List of secrets to manage
  - alias: db_password
    provider: aws
    name: myapp/prod/db_password
    env: DB_PASSWORD
    extras:
      region: us-east-1
```

### Provider System

The provider system (`internal/provider/`) is the core extensibility mechanism:

#### Provider Interface

```go
type Provider interface {
    FetchSecret(ctx context.Context, spec SecretSpec) (string, error)
}

type SecretSpec struct {
    Alias    string            // Human-readable alias
    Name     string            // Provider-specific secret name/path
    Provider string            // Provider type (aws, gcp, etc.)
    EnvName  string            // Environment variable name
    Extras   map[string]string // Provider-specific options
}
```

#### Provider Registry

Simple map-based registry allows multiple aliases per provider:

```go
// Register providers with multiple aliases
provider.Register("aws", awsprovider.New())
provider.Register("aws-secrets-manager", awsprovider.New())
provider.Register("aws-ssm", awsprovider.NewSSM())
provider.Register("ssm", awsprovider.NewSSM())
```

#### Supported Providers

- **AWS**: Secrets Manager (`aws`) and SSM Parameter Store (`aws-ssm`)
- **GCP**: Google Secret Manager (`gcp`)
- **Azure**: Key Vault (`azure`) and App Configuration (`azure-appconfig`)
- **HashiCorp Vault**: KV v2 and logical reads (`vault`)
- **Exec**: External command execution (`exec`)

### Security Model

Security is built into every layer:

#### Memory-Only Secrets

- Secrets are never written to disk
- Values only exist in memory and child process environments
- No temporary files or caching

#### Secret Masking

- Secret values are masked in logs and dry-run output by default
- Configurable masking behavior
- Debug mode available for troubleshooting

#### Secure Execution

- No shell evaluation - arguments passed directly to `exec`
- Context-based timeouts and cancellation
- Proper error handling without secret leakage

## Data Flow

1. **Configuration Loading**

   - Discover config file location
   - Parse YAML and validate schema
   - Interpolate environment variables
   - Merge defaults with individual secrets

2. **Provider Resolution**

   - Look up provider by name in registry
   - Create SecretSpec for each secret
   - Apply provider-specific options from extras

3. **Secret Fetching**

   - Concurrent fetching with configurable limits
   - Context-based timeouts and cancellation
   - Retry logic for transient failures
   - Error mapping to consistent types

4. **Command Execution**
   - Build environment with fetched secrets
   - Execute target command with `os/exec`
   - Clean up secrets from memory
   - Return exit code from child process

## Extension Points

### Adding New Providers

1. **Create provider package** under `internal/provider/<name>/`
2. **Implement Provider interface** with `FetchSecret` method
3. **Add unit tests** with mocked SDK clients
4. **Register in main.go** with appropriate aliases
5. **Document in providers.md** with usage examples

### Adding New Commands

1. **Create command file** in `cmd/skv/`
2. **Implement Cobra command** with proper flags and validation
3. **Add unit and integration tests**
4. **Update CLI documentation**

### Configuration Extensions

The configuration system supports:

- New provider-specific extras
- Additional global defaults
- Custom validation rules
- Environment variable interpolation

## Testing Strategy

### Unit Tests

- Each provider has comprehensive unit tests
- SDK clients are mocked using interfaces
- Error conditions and edge cases covered
- Fast execution with no external dependencies

### Integration Tests

- End-to-end command testing
- Mock provider for consistent behavior
- Configuration validation testing
- Error handling verification

### E2E Tests

- Real provider integration (when credentials available)
- Full command workflows
- Cross-platform compatibility
- Performance and concurrency testing

## Build and Release

### Build System

- **Makefile**: Single source of build commands
- **Go modules**: Dependency management
- **Cross-compilation**: Multiple OS/architecture support
- **Version injection**: Build-time version information

### Release Process

- **GoReleaser**: Automated binary building and release
- **GitHub Actions**: CI/CD pipeline
- **Binary-only releases**: No package manager dependencies
- **Checksums**: Integrity verification

## Performance Considerations

### Concurrency

- Configurable concurrent secret fetching
- Context-based cancellation
- Proper goroutine cleanup
- Rate limiting for provider APIs

### Memory Management

- Minimal memory footprint
- No secret caching or persistence
- Efficient string handling
- Garbage collection friendly

### Network Optimization

- Connection reuse where possible
- Configurable timeouts and retries
- Proper error handling for network issues
- Support for proxy configurations

## Future Considerations

### Potential Extensions

- Plugin system for external providers
- Secret caching with TTL (optional)
- Audit logging capabilities
- Integration with secret rotation systems
- Support for secret templates/transformations

### Scalability

- Batch secret operations
- Streaming for large numbers of secrets
- Distributed secret fetching
- Integration with service mesh systems
