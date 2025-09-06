# skv Documentation

**skv** (Secure Key/Value Manager) is a cross-cloud CLI that unifies secret management across multiple providers. Fetch secrets from AWS, GCP, Azure, HashiCorp Vault, or custom commands, then inject them securely into processes, export to files, or print to stdout.

## Key Features

- **Cross-cloud unified config** - Single YAML for all providers
- **Secure by design** - Memory-only secrets, never written to disk
- **Process injection** - Safely inject secrets into command environments
- **Flexible output** - Print, export to .env, or inject into processes
- **Provider extensibility** - Easy to add new secret backends
- **Production ready** - Timeouts, retries, health checks, validation

## Quickstart

1. **Install**: see [installation.md](installation.md)
2. **Generate config template**: `skv init`
3. **Edit** `~/.skv.yaml`:

```yaml
secrets:
  - alias: db_password
    provider: aws
    name: myapp/prod/db_password
    env: DB_PASSWORD
    extras:
      region: us-east-1
```

4. **Validate and use**:

```bash
skv validate                           # Check configuration
skv list                              # Show configured secrets
skv get db_password                   # Get single secret
skv export --all --format env > .env # Export to file
skv run --all -- env | grep DB_PASSWORD  # Run with secrets
```

## Documentation

### Getting Started

- **[Installation Guide](installation.md)** - Platform-specific installation instructions
- **[Configuration](configuration.md)** - YAML config reference and examples
- **[CLI Reference](cli.md)** - Complete command documentation

### Providers

- **[Providers Overview](providers.md)** - AWS, GCP, Azure, Vault, and Exec provider guides
- **[Examples](examples.md)** - Real-world usage scenarios

### Operations

- **[Troubleshooting](troubleshooting.md)** - Common issues and solutions
- **[Security Checklist](security-checklist.md)** - Security best practices

### Development

- **[Architecture](dev/architecture.md)** - Project structure and design
- **[Developing Providers](dev/developing-providers.md)** - Adding new providers
- **[Conventions](dev/conventions.md)** - Code style and project conventions
- **[Cross-Platform Development](dev/cross-platform.md)** - Development on Windows/macOS/Linux
