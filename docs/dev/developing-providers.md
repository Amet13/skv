# Developing Providers

This guide explains how to add a new secret provider to `skv`.

## Overview

A provider is a component that fetches secrets from a specific backend (cloud service, local command, etc.). All providers implement a simple interface and are registered in the main application.

## Provider Interface

### Core Interface

Every provider must implement the `Provider` interface:

```go
type Provider interface {
    FetchSecret(ctx context.Context, spec SecretSpec) (string, error)
}
```

### SecretSpec Structure

The `SecretSpec` contains all information needed to fetch a secret:

```go
type SecretSpec struct {
    Alias    string            // Human-readable alias (e.g., "db_password")
    Name     string            // Provider-specific resource identifier
    Provider string            // Provider type (e.g., "aws", "gcp")
    EnvName  string            // Environment variable name
    Extras   map[string]string // Provider-specific configuration
}
```

## Implementation Steps

### 1. Create Provider Package

Create a new package under `internal/provider/<name>/`:

```
internal/provider/
|-── mycloud/
|   |-── mycloud.go      # Main implementation
|   `-── mycloud_test.go # Unit tests
```

### 2. Implement the Provider

```go
package mycloud

import (
    "context"
    "errors"
    "fmt"

    "skv/internal/provider"
)

// Client interface for testing
type secretClient interface {
    GetSecret(ctx context.Context, name string, options map[string]string) (string, error)
}

type mycloudProvider struct {
    client secretClient
}

// New creates a new MyCloud provider
func New() provider.Provider {
    return &mycloudProvider{
        client: &realSecretClient{}, // Real implementation
    }
}

func (p *mycloudProvider) FetchSecret(ctx context.Context, spec provider.SecretSpec) (string, error) {
    value, err := p.client.GetSecret(ctx, spec.Name, spec.Extras)
    if err != nil {
        // Map provider-specific errors to standard types
        if isNotFoundError(err) {
            return "", provider.ErrNotFound
        }
        return "", fmt.Errorf("mycloud provider error: %w", err)
    }

    return value, nil
}

func isNotFoundError(err error) bool {
    // Check for provider-specific "not found" error types
    var notFoundErr *MyCloudNotFoundError
    return errors.As(err, &notFoundErr)
}
```

### 3. Create Testable Client Interface

```go
// Real client implementation
type realSecretClient struct {
    // SDK client instance
}

func (c *realSecretClient) GetSecret(ctx context.Context, name string, options map[string]string) (string, error) {
    // Actual SDK calls
    return "", nil
}

// Mock client for testing
type mockSecretClient struct {
    secrets map[string]string
    err     error
}

func (c *mockSecretClient) GetSecret(ctx context.Context, name string, options map[string]string) (string, error) {
    if c.err != nil {
        return "", c.err
    }

    value, exists := c.secrets[name]
    if !exists {
        return "", &MyCloudNotFoundError{Name: name}
    }

    return value, nil
}
```

### 4. Write Comprehensive Tests

```go
package mycloud

import (
    "context"
    "testing"

    "skv/internal/provider"
)

func TestMycloudProvider_FetchSecret(t *testing.T) {
    tests := []struct {
        name    string
        spec    provider.SecretSpec
        client  secretClient
        want    string
        wantErr error
    }{
        {
            name: "successful fetch",
            spec: provider.SecretSpec{
                Name: "test-secret",
                Extras: map[string]string{"region": "us-east-1"},
            },
            client: &mockSecretClient{
                secrets: map[string]string{"test-secret": "secret-value"},
            },
            want: "secret-value",
        },
        {
            name: "secret not found",
            spec: provider.SecretSpec{Name: "nonexistent"},
            client: &mockSecretClient{
                secrets: map[string]string{},
            },
            wantErr: provider.ErrNotFound,
        },
        {
            name: "client error",
            spec: provider.SecretSpec{Name: "test-secret"},
            client: &mockSecretClient{
                err: errors.New("client error"),
            },
            wantErr: errors.New("mycloud provider error: client error"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := &mycloudProvider{client: tt.client}

            got, err := p.FetchSecret(context.Background(), tt.spec)

            if tt.wantErr != nil {
                if err == nil || err.Error() != tt.wantErr.Error() {
                    t.Errorf("FetchSecret() error = %v, wantErr %v", err, tt.wantErr)
                }
                return
            }

            if err != nil {
                t.Errorf("FetchSecret() unexpected error = %v", err)
                return
            }

            if got != tt.want {
                t.Errorf("FetchSecret() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 5. Register Provider

Add registration in `cmd/skv/main.go`:

```go
import (
    // ... other imports
    mycloudprovider "skv/internal/provider/mycloud"
)

func init() {
    // ... existing registrations

    // MyCloud provider
    provider.Register("mycloud", mycloudprovider.New())
    provider.Register("mc", mycloudprovider.New()) // Short alias
}
```

### 6. Document the Provider

Update `docs/providers.md` with a new section:

````markdown
### MyCloud

- **Auth**: MyCloud API key via `MYCLOUD_API_KEY` environment variable
- **Name**: Secret name or full resource path
- **Extras** (optional):
  - `region`: MyCloud region (e.g., `us-east-1`)
  - `project`: Project ID
  - `version`: Secret version (defaults to `latest`)

```yaml
secrets:
  - alias: api_key
    provider: mycloud
    name: myapp/prod/api-key
    env: API_KEY
    extras:
      region: us-east-1
      project: my-project
      version: latest
```
````

````

### 7. Integration Testing

Create integration tests that can run with real credentials:

```go
func TestMycloudProvider_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    apiKey := os.Getenv("MYCLOUD_API_KEY")
    if apiKey == "" {
        t.Skip("MYCLOUD_API_KEY not set")
    }

    provider := New()
    spec := provider.SecretSpec{
        Name: "integration-test-secret",
        Extras: map[string]string{
            "region": "us-east-1",
        },
    }

    value, err := provider.FetchSecret(context.Background(), spec)
    if err != nil {
        t.Fatalf("FetchSecret() error = %v", err)
    }

    if value == "" {
        t.Error("FetchSecret() returned empty value")
    }
}
````

## Best Practices

### Error Handling

- Always map provider-specific "not found" errors to `provider.ErrNotFound`
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Don't include secret values in error messages
- Provide actionable error messages when possible

### Configuration

- Use `spec.Extras` for provider-specific options
- Support common options like `region`, `profile`, `project`
- Provide sensible defaults
- Validate required configuration

### Testing

- Mock external dependencies using interfaces
- Test both success and error paths
- Include edge cases and boundary conditions
- Use table-driven tests for multiple scenarios
- Aim for >90% test coverage for provider code

### Performance

- Reuse connections when possible
- Respect context cancellation and timeouts
- Implement proper cleanup
- Consider rate limiting for API calls

### Security

- Never log secret values
- Handle authentication securely
- Validate inputs to prevent injection attacks
- Use secure defaults

## Common Patterns

### Authentication

```go
// Environment variable authentication
func (p *provider) getCredentials() (*Credentials, error) {
    apiKey := os.Getenv("MYCLOUD_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("MYCLOUD_API_KEY environment variable required")
    }
    return &Credentials{APIKey: apiKey}, nil
}

// Profile-based authentication
func (p *provider) getCredentialsFromProfile(profile string) (*Credentials, error) {
    // Load from config file based on profile
    return loadProfile(profile)
}
```

### Region/Endpoint Configuration

```go
func (p *provider) getEndpoint(extras map[string]string) string {
    if region := extras["region"]; region != "" {
        return fmt.Sprintf("https://%s.mycloud.com", region)
    }
    return "https://api.mycloud.com" // Default
}
```

### Version Handling

```go
func (p *provider) getSecretVersion(name string, extras map[string]string) string {
    if version := extras["version"]; version != "" {
        return version
    }
    return "latest" // Default
}
```

## Validation and Testing

### Pre-submission Checklist

- [ ] Provider implements `Provider` interface correctly
- [ ] Unit tests cover success and error paths
- [ ] Integration tests work with real credentials
- [ ] Provider registered in `main.go`
- [ ] Documentation updated in `providers.md`
- [ ] All linting passes (`make lint`)
- [ ] Test coverage >90% for provider code
- [ ] Error handling follows conventions
- [ ] No secret values in logs or error messages

### Testing Commands

```bash
# Run unit tests
go test ./internal/provider/mycloud/

# Run with coverage
go test -cover ./internal/provider/mycloud/

# Run integration tests
go test -tags=integration ./internal/provider/mycloud/

# Run all tests
make test

# Check linting
make lint
```

## Examples

See existing providers for reference:

- **Simple provider**: `internal/provider/exec/` - minimal implementation
- **Cloud provider**: `internal/provider/aws/` - authentication, regions, error handling
- **Multiple backends**: `internal/provider/azure/` - Key Vault and App Configuration

Each provider demonstrates different patterns and can serve as a template for new implementations.
