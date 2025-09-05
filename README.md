<div align="center">
  <img src="images/logo.png" alt="skv Logo" width="200" style="border-radius: 20px;">
  <h1>skv</h1>
  <p><strong>CLI tool that fetches secrets from different providers</strong></p>
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

Fetch secrets from multiple providers and storages (AWS Secrets Manager, AWS SSM Parameter Store, Google Secret Manager, Azure Key Vault, HashiCorp Vault, and custom Exec) and inject them into a process environment, export them, or print to stdout.

### Why skv?

- Unify access to secrets across clouds and backends with a single CLI and config.
- Inject secrets into any process safely, with masking and dry-run.
- Concurrency, retries, timeouts, and per-provider options (region, profile, project, namespace, etc.).

### Supported providers and storages

- AWS: Secrets Manager (`provider: aws`), SSM Parameter Store (`provider: aws-ssm`)
- GCP: Secret Manager (`provider: gcp`)
- Azure: Key Vault (`provider: azure`)
- HashiCorp Vault: KV v2 and logical (`provider: vault`)
- Exec: run trusted local command to get dynamic secrets (`provider: exec`)

## 🚀 **Quick Start**

1. **Install**: Download from [releases](https://github.com/Amet13/skv/releases/latest) or see [installation guide](docs/installation.md)
2. **Configure**: `skv init` to generate config template
3. **Use**: `skv run --all -- your-command`

```bash
# Install (example for Linux)
curl -sL "https://github.com/Amet13/skv/releases/latest/download/skv_linux_amd64" -o skv
chmod +x skv && sudo mv skv /usr/local/bin/

# Setup
skv init                    # Generate ~/.skv.yaml with examples
skv validate                # Verify configuration
skv list                    # Show configured secrets

# Basic usage
skv get db-password         # Fetch single secret
skv run --all -- env        # Inject all secrets into process
skv export --all > .env     # Export to .env file
```

## 📚 **Documentation**

For comprehensive documentation, see [`docs/index.md`](docs/index.md) with links to:

- Installation guide, CLI reference, configuration details
- Provider-specific guides, examples, and troubleshooting
- Security checklist and development guidelines

## 🛠️ **Development**

```bash
make build     # Build for current platform
make lint      # Run all linters and tests
make release   # Create release artifacts
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed development setup.

## 🔒 **Security**

- ✅ **No disk writes** - Secrets are never written to disk by this tool
- 🧠 **Memory only** - Values are only present in memory and child process environment during `run`
- 🎭 **Secret masking** - Secret values are masked in dry-run output and logs by default
- 🔐 **Secure by design** - Built with security-first principles

## 📄 **License**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
