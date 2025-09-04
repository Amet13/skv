<div align="center">
  <img src="images/logo.png" alt="skv Logo" width="200" style="border-radius: 20px;">
  <h1>skv</h1>
  <p><strong>CLI tool that fetches secrets from different providers</strong></p>
  <p>
    <a href="https://github.com/Amet13/skv/actions/workflows/ci.yml">
      <img src="https://github.com/Amet13/skv/actions/workflows/ci.yml/badge.svg" alt="CI Status">
    </a>
  </p>
</div>

Fetch secrets from providers like AWS Secrets Manager and HashiCorp Vault, and inject them into a process environment or print them to stdout.

## Installation

- Download a [latest release](https://github.com/Amet13/skv/releases/latest/).
- Rename it to `skv` and make it executable:

```bash
chmod +x ./skv_*
mv ./skv_* /usr/local/bin/skv
skv version
```

## Configuration

Default config file path: `$HOME/.skv.yaml`. Override with `SKV_CONFIG` or `--config`.

```yaml
secrets:
  - alias: "db-password"
    provider: "aws-secrets-manager"
    name: "myapp/prod/db-password"
    region: "us-east-1"

  - alias: "api-token"
    provider: "vault"
    name: "secret/data/tokens/myapp"
    address: "https://vault.mycompany.com"
    token: "{{ VAULT_TOKEN }}"
    env: "MYAPP_API_TOKEN"
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
skv completion zsh > ~/.zfunc/_skv && echo 'fpath+=(~/.zfunc)' >> ~/.zshrc
# Fish
skv completion fish > ~/.config/fish/completions/skv.fish
```

Fetch a single secret and print to stdout:

```bash
skv get db-password
```

Inject secrets into a process:

```bash
skv run --all -- env | grep -E "DB_PASSWORD|MYAPP_API_TOKEN"
skv run -s db-password -s api-token -- ./your-app --flag
skv run --secrets db-password,api-token --dry-run -- echo "hello"
```

Common flags:

- `--config`: path to config file
- `--log-level`: error|warn|info|debug
- `--timeout`: e.g., 5s, 30s
- `--dry-run`: show what would happen without executing

## Documentation

- CLI: `docs/cli.md`
- Configuration: `docs/configuration.md`
- Providers: `docs/providers.md`

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
