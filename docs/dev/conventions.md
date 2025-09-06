# Project Conventions

This document outlines coding standards, naming conventions, and best practices for the `skv` project.

## Naming Conventions

### Provider Aliases

Provider names used in configuration (`provider:` field):

- **AWS Secrets Manager**: `aws`, `aws-secrets-manager`
- **AWS SSM Parameter Store**: `aws-ssm`, `ssm`, `aws-parameter-store`
- **Google Secret Manager**: `gcp`, `gcp-secret-manager`
- **Azure Key Vault**: `azure`, `azure-key-vault`
- **Azure App Configuration**: `azure-appconfig`, `appconfig`
- **HashiCorp Vault**: `vault`
- **Exec**: `exec`

### Code Naming

- **Go packages**: lowercase, single word when possible (`provider`, `config`)
- **Go types**: PascalCase (`SecretSpec`, `Provider`)
- **Go functions**: PascalCase for exported, camelCase for private
- **Go variables**: camelCase (`secretValue`, `providerName`)
- **Go constants**: PascalCase or UPPER_SNAKE_CASE for exported constants

### Configuration

- **Config keys**: snake_case (`version_stage`, `vault_url`)
- **Extras keys**: snake_case or lowercase hyphenated (`region`, `credentials-file`)
- **Environment variables**: UPPER_SNAKE_CASE (`DB_PASSWORD`, `API_KEY`)
- **CLI flags**: kebab-case (`--dry-run`, `--log-level`)

### Files and Directories

- **Go files**: snake_case (`secret_spec.go`, `provider_test.go`)
- **Documentation**: lowercase with hyphens (`installation.md`, `security-checklist.md`)
- **Directories**: lowercase (`internal/provider/aws/`)

## Code Standards

### Go Code Style

Follow standard Go conventions with these specific guidelines:

#### Error Handling

```go
// Good: Wrap errors with context
_, err := client.FetchSecret(ctx, spec)
if err != nil {
    return fmt.Errorf("failed to fetch secret %q from %s: %w", spec.Alias, spec.Provider, err)
}

// Bad: Lose error context
_, err := client.FetchSecret(ctx, spec)
if err != nil {
    return err
}
```

#### Documentation

```go
// Good: Clear, professional documentation
// FetchSecret retrieves a secret value from the provider.
// It returns the secret value or an error if the secret cannot be retrieved.
// The context can be used for cancellation and timeout control.
func (p *provider) FetchSecret(ctx context.Context, spec SecretSpec) (string, error) {
    // Implementation
}

// Bad: Bad: Missing or unclear documentation
// Gets secret
func (p *provider) FetchSecret(ctx context.Context, spec SecretSpec) (string, error) {
    // Implementation
}
```

#### Testing

```go
// Good: Good: Table-driven tests with clear structure
func TestProviderFetchSecret(t *testing.T) {
    tests := []struct {
        name    string
        spec    provider.SecretSpec
        want    string
        wantErr bool
    }{
        {
            name: "successful secret fetch",
            spec: provider.SecretSpec{
                Name: "test-secret",
                Extras: map[string]string{"region": "us-east-1"},
            },
            want:    "secret-value",
            wantErr: false,
        },
        {
            name: "secret not found",
            spec: provider.SecretSpec{Name: "nonexistent"},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Provider Implementation

#### Interface Compliance

```go
// Good: Good: Proper interface implementation
type awsProvider struct {
    client secretsManagerClient // Interface for testing
}

func (p *awsProvider) FetchSecret(ctx context.Context, spec provider.SecretSpec) (string, error) {
    // Implementation with proper error handling
}

// Good: Good: Testable client interface
type secretsManagerClient interface {
    GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}
```

#### Error Mapping

```go
// Good: Good: Map provider errors to standard types
func (p *awsProvider) FetchSecret(ctx context.Context, spec provider.SecretSpec) (string, error) {
    output, err := p.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
        SecretId: aws.String(spec.Name),
    })
    if err != nil {
        var notFound *types.ResourceNotFoundException
        if errors.As(err, &notFound) {
            return "", provider.ErrNotFound
        }
        return "", fmt.Errorf("aws secrets manager error: %w", err)
    }

    return aws.ToString(output.SecretString), nil
}
```

## Documentation Standards

### Markdown Formatting

- **Headers**: Use `#` for top-level, `##` for sections, `###` for subsections
- **Code**: Use backticks for inline code, triple backticks for blocks
- **Links**: Use angle brackets for URLs: `<https://example.com>`
- **Emphasis**: Use `**bold**` for important terms, `*italic*` for emphasis
- **Lists**: Use `-` for unordered lists, numbers for ordered lists

### Code Examples

```yaml
# Good: Good: Complete, working examples
secrets:
  - alias: db_password
    provider: aws
    name: myapp/prod/db_password
    env: DB_PASSWORD
    extras:
      region: us-east-1
      version_stage: AWSCURRENT
```

```bash
# Good: Good: Practical, copy-pasteable commands
skv run --all --timeout 30s -- ./my-app
```

### Documentation Structure

1. **Purpose**: What the document covers
2. **Prerequisites**: What users need to know/have
3. **Examples**: Practical, working examples
4. **Reference**: Detailed parameter/option documentation
5. **Troubleshooting**: Common issues and solutions

## Testing Standards

### Unit Tests

- **Coverage**: Aim for >70% overall, 100% for critical paths
- **Isolation**: Tests must not depend on external services
- **Mocking**: Use interfaces to mock SDK clients
- **Table-driven**: Use table-driven tests for multiple scenarios

### Integration Tests

- **Mock providers**: Use mock provider for consistent behavior
- **End-to-end**: Test complete command workflows
- **Error scenarios**: Test error handling and edge cases

### Test Organization

```go
// Good: Good: Clear test organization
func TestAWSProvider(t *testing.T) {
    t.Run("FetchSecret", func(t *testing.T) {
        // Test FetchSecret method
    })

    t.Run("Authentication", func(t *testing.T) {
        // Test authentication scenarios
    })
}
```

## Build and Release Standards

### Makefile Targets

Standard targets that should be available:

- `make build`: Build for current platform
- `make build-all`: Build for all platforms
- `make test`: Run all tests
- `make lint`: Run all linters
- `make clean`: Clean build artifacts

### Version Management

- Use semantic versioning (MAJOR.MINOR.PATCH)
- Version injected at build time via ldflags
- Git tags for releases

### Dependencies

- Minimize external dependencies
- Use Go standard library when possible
- Pin dependency versions in go.mod
- Regular dependency updates and security scanning

## Security Standards

### Secret Handling

- Never log secret values
- Use masking in debug output
- Memory-only storage (no disk writes)
- Proper cleanup of sensitive data

### Error Messages

- Don't include secret values in error messages
- Provide enough context for debugging
- Don't expose internal system details

### Input Validation

- Validate all configuration inputs
- Sanitize provider-specific options
- Handle malformed data gracefully

## Performance Guidelines

### Concurrency

- Use context for cancellation and timeouts
- Implement proper goroutine cleanup
- Avoid goroutine leaks
- Use sync primitives correctly

### Memory Management

- Minimize allocations in hot paths
- Use string builders for concatenation
- Avoid unnecessary copying of large data
- Profile memory usage for optimization

## Code Review Checklist

### Before Submitting

- [ ] All tests pass
- [ ] Linting passes
- [ ] Documentation updated
- [ ] Security implications considered
- [ ] Performance impact assessed
- [ ] Error handling implemented
- [ ] Logging appropriate for level

### Review Focus Areas

- [ ] Code follows established patterns
- [ ] Error handling is consistent
- [ ] Tests cover new functionality
- [ ] Documentation is clear and accurate
- [ ] Security best practices followed
- [ ] Performance considerations addressed
