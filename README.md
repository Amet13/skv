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

## Installation

- Download a [latest release](https://github.com/Amet13/skv/releases/latest/).
- Rename it to `skv` and make it executable:

```bash
chmod +x ./skv_*
sudo mv ./skv_* /usr/local/bin/skv
# For macOS, add to quarantine
sudo xattr -rd com.apple.quarantine /usr/local/bin/skv

skv version
```

## Configuration

Default config discovery:

1. `--config` flag
2. `SKV_CONFIG` env var
3. `$HOME/.skv.yaml` or `$HOME/.skv.yml`

```yaml
defaults:
  region: us-east-1

secrets:
  # AWS Secrets Manager
  - alias: db_password
    provider: aws
    name: myapp/prod/db_password
    env: DB_PASSWORD
    extras:
      region: us-east-1
      version_stage: AWSCURRENT

  # AWS SSM Parameter Store
  - alias: db_host
    provider: aws-ssm
    name: /myapp/prod/db_host
    env: DB_HOST
    extras:
      region: us-east-1
      with_decryption: "true"

  # Vault KV v2
  - alias: api_token
    provider: vault
    name: kv/data/tokens/myapp
    env: MYAPP_API_TOKEN
    extras:
      address: https://vault.mycompany.com
      token: "{{ VAULT_TOKEN }}"
```

Notes:

- `{{ VAR_NAME }}` placeholders are interpolated from the environment during load; missing values cause an error.
- If `env` is omitted, the environment variable name is derived from the alias in UPPER_SNAKE_CASE.

## Usage

### Quick Start

Generate a configuration template:

```bash
skv init                    # Create ~/.skv.yaml with examples
skv init --provider aws     # AWS-focused template
skv init --output .skv.yaml # Custom output location
```

Validate your configuration:

```bash
skv validate                      # Check configuration syntax
skv validate --check-secrets      # Test connectivity to providers
skv validate --verbose            # Show detailed validation results
```

### Basic Commands

Show version and help:

```bash
skv version
skv --help
skv <command> --help
```

List configured secrets:

```bash
skv list                    # Show all configured secrets
skv list --format json     # JSON output for scripting
```

Fetch a single secret:

```bash
skv get db-password                           # Get secret value
skv get db-password --timeout 5s --retries 2 # With timeout and retries
```

### Running Applications with Secrets

Inject secrets into a process:

```bash
skv run --all -- env | grep -E "DB_PASSWORD|API_TOKEN"
skv run -s db-password -s api-token -- ./your-app --flag
skv run --secrets db-password,api-token --dry-run -- echo "hello"
```

### Export and Health Checks

Export secrets to files:

```bash
skv export --all --format env > .env
skv export --secrets a,b --format yaml --output secrets.yaml
```

Health checks for monitoring:

```bash
skv health                    # Check all secrets
skv health --secret db-pass   # Check specific secret
skv health --timeout 30s      # Custom timeout
```

### Shell Completion

Generate completion scripts:

```bash
# Bash
skv completion bash > /usr/local/etc/bash_completion.d/skv
# Zsh
mkdir -p ~/.zfunc/ && skv completion zsh > ~/.zfunc/_skv && echo 'fpath+=(~/.zfunc)' >> ~/.zshrc
# Fish
skv completion fish > ~/.config/fish/completions/skv.fish
```

### Common Flags

- `--config`: path to config file
- `--log-level`: error|warn|info|debug
- `--timeout`: e.g., 5s, 30s (per-command fetch)
- `--retries` and `--retry-delay`: simple retry policy for fetches
- `--dry-run`: show what would happen without executing

## 📚 **Documentation**

- 🏠 **Docs index**: [`docs/index.md`](docs/index.md)
- 💻 **CLI reference**: [`docs/cli.md`](docs/cli.md)
- ⚙️ **Configuration**: [`docs/configuration.md`](docs/configuration.md)
- 🔌 **Providers**: [`docs/providers.md`](docs/providers.md)
- 🛠️ **Developer guide**: [`docs/dev/developing-providers.md`](docs/dev/developing-providers.md)
- 🚨 **Troubleshooting**: [`docs/troubleshooting.md`](docs/troubleshooting.md)
- 🔒 **Security checklist**: [`docs/security-checklist.md`](docs/security-checklist.md)
- 📋 **Examples**: [`docs/examples.md`](docs/examples.md)
- 📝 **Conventions**: [`docs/dev/conventions.md`](docs/dev/conventions.md)
- 🚀 **Installation**: [`docs/installation.md`](docs/installation.md)
- 🔄 **Migration guide**: [`docs/migration.md`](docs/migration.md)

## 🛠️ **Development**

**Requirements:**

- 🐹 Go 1.25.x

```bash
make build     # 🏗️ Build for host platform
make build-all # 🌍 Cross-build for multiple platforms
```

Built artifacts are written to `dist/`.

## 🎯 **Common Tasks (Make)**

```bash
make fmt       # 🎨 Format Go, shell, markdown, YAML
make lint      # 🔍 Run all linters
make test      # 🧪 Run tests
make cover     # 📊 Run tests with coverage gate
make clean     # 🧹 Remove dist/
make release   # 📦 Local snapshot via GoReleaser (no publish)
```

## 🔒 **Security**

- ✅ **No disk writes** - Secrets are never written to disk by this tool
- 🧠 **Memory only** - Values are only present in memory and child process environment during `run`
- 🎭 **Secret masking** - Secret values are masked in dry-run output and logs by default
- 🔐 **Secure by design** - Built with security-first principles

## 📄 **License**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
