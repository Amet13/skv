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
- Planned: Oracle (OCI), IBM Cloud, Alibaba Cloud (currently registered as stubs)

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

Show version:

```bash
skv version
```

Generate completion:

```bash
# Bash
skv completion bash > /usr/local/etc/bash_completion.d/skv
# Zsh
mkdir -p ~/.zfunc/ && skv completion zsh > ~/.zfunc/_skv && echo 'fpath+=(~/.zfunc)' >> ~/.zshrc
# Fish
skv completion fish > ~/.config/fish/completions/skv.fish
```

Fetch a single secret and print to stdout:

```bash
skv get db-password --timeout 5s --retries 2 --retry-delay 250ms
```

Inject secrets into a process:

```bash
skv run --all -- env | grep -E "DB_PASSWORD|MYAPP_API_TOKEN"
skv run -s db-password -s api-token --retries 2 -- ./your-app --flag
skv run --secrets db-password,api-token --dry-run -- echo "hello"
```

List and export:

```bash
skv list --format json
skv export --all --format env > .env
skv export --secrets a,b --format yaml --output secrets.yaml
```

Common flags:

- `--config`: path to config file
- `--log-level`: error|warn|info|debug
- `--timeout`: e.g., 5s, 30s (per-command fetch)
- `--retries` and `--retry-delay`: simple retry policy for fetches
- `--dry-run`: show what would happen without executing

## Documentation

- Docs index: `docs/index.md`
- CLI: `docs/cli.md`
- Configuration: `docs/configuration.md`
- Providers and storages: `docs/providers.md`
- Developer guide (extensibility): `docs/DEVELOPING_PROVIDERS.md`
- Troubleshooting: `docs/TROUBLESHOOTING.md`
- Security checklist: `docs/SECURITY_CHECKLIST.md`
- Examples: `docs/EXAMPLES.md`
- Conventions: `docs/CONVENTIONS.md`

## Development

Requirements:

- Go 1.25.x

```bash
make build     # host
make build-all # small cross-build matrix
```

Built artifacts are written to `dist/`.

## Common tasks (Make)

```bash
make fmt       # format Go, shell, markdown, YAML
make lint      # run all linters
make test      # run tests
make cover     # run tests with coverage gate
make clean     # remove dist/
make release   # local snapshot via GoReleaser (no publish)
```

## Security

- Secrets are never written to disk by this tool.
- Values are only present in memory and in the environment of the child process during `run`.
- Secret values are masked in dry-run output and logs by default.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
