<div align="center">
  <img src="images/logo.png" alt="skv Logo" width="200" style="border-radius: 20px;">
  <h1>skv</h1>
  <p><strong>Secure Key/Value Manager</strong></p>
  <p>
    <a href="https://github.com/Amet13/skv/actions/workflows/ci.yml">
      <img src="https://github.com/Amet13/skv/actions/workflows/ci.yml/badge.svg" alt="CI Status">
    </a>
    <a href="https://github.com/Amet13/skv/actions/workflows/codeql.yml">
      <img src="https://github.com/Amet13/skv/actions/workflows/codeql.yml/badge.svg" alt="CodeQL">
    </a>
    <a href="https://github.com/Amet13/skv/releases">
      <img src="https://img.shields.io/github/v/release/Amet13/skv?label=version" alt="Latest Release">
    </a>
  </p>
</div>

**skv** unifies secret management across cloud providers with a single configuration file and CLI. Fetch secrets from AWS, GCP, Azure, HashiCorp Vault, or custom commands, then inject them securely into processes, export to files, or print to stdout.

### Supported providers and storages

- **AWS**: Secrets Manager (`aws`), SSM Parameter Store (`aws-ssm`)
- **GCP**: Secret Manager (`gcp`)
- **Azure**: Key Vault (`azure`), App Configuration (`azure-appconfig`)
- **HashiCorp Vault**: KV v2 and logical (`vault`)
- **Exec**: run trusted local command to get dynamic secrets (`exec`)

## Quick Start

```bash
# Install
curl -sL "https://github.com/Amet13/skv/releases/latest/download/skv_linux_amd64" -o skv
chmod +x skv && sudo mv skv /usr/local/bin/

# Setup and use
skv init                          # Generate ~/.skv.yaml template
skv doctor                        # Run diagnostics and health checks
skv completion install            # Install shell completions
skv get db-password               # Fetch single secret
skv run --all -- env              # Inject all secrets into process
skv watch --all -- echo "changed" # Watch secrets for changes
```

See [installation guide](docs/installation.md) for other platforms and [documentation](docs/index.md) for full usage.

## Documentation

- **[Installation Guide](docs/installation.md)** - Platform-specific installation instructions
- **[Configuration](docs/configuration.md)** - YAML config reference and examples
- **[Providers](docs/providers.md)** - AWS, GCP, Azure, Vault, and Exec provider guides
- **[CLI Reference](docs/cli.md)** - Complete command documentation
- **[Examples](docs/examples.md)** - Real-world usage scenarios

Full documentation index: [`docs/index.md`](docs/index.md)

## Key Features

- **Cross-cloud unified config** - Single YAML for all providers
- **Secure by design** - Memory-only secrets, never written to disk
- **Process injection** - Safely inject secrets into command environments
- **Flexible output** - Print, export to .env, or inject into processes
- **Provider extensibility** - Easy to add new secret backends
- **Production ready** - Timeouts, retries, health checks, validation

## Development

```bash
make build && make lint
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and [docs/dev/](docs/dev/) for architecture details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
