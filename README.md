<div align="center">
  <img src="images/logo.png" alt="skv Logo" width="200" style="border-radius: 20px;">
  <h1>skv</h1>
  <p><strong>Small CLI tool that fetches secrets from provider</strong></p>
  <p>
    <a href="https://github.com/Amet13/skv/actions/workflows/ci.yml">
      <img src="https://github.com/Amet13/skv/actions/workflows/ci.yml/badge.svg" alt="CI Status">
    </a>
  </p>
</div>

skv is a small Go CLI that fetches secrets from providers like AWS Secrets Manager and HashiCorp Vault, and injects them into a process environment or prints them to stdout.

## Install

Requires Go 1.25.1 or newer. Recommended: latest Go 1.25.x.

Build locally:

```bash
bash ./build.sh host
```

Artifacts are written to `dist/`.

Optional optimizations (Go 1.25+):

- Profile-guided optimization (PGO): set `PGO_PROFILE=path/to/profile.pprof` before running `build.sh`.
- Experimental GC: you can set `GOEXPERIMENT=greenteagc` to evaluate impact (available in recent Go versions).

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

## Providers

### AWS Secrets Manager

- Uses the default AWS credential chain; override region via `region` in the secret entry.
- Returns `SecretString` if present, otherwise the binary payload converted to string.

### HashiCorp Vault

### GCP Secret Manager

- Provide `metadata.project` and optional `metadata.version` if `name` is not a full resource path.

### Azure Key Vault

- Provide `metadata.vault_url` (e.g., <https://myvault.vault.azure.net/>) and optional `metadata.version`.

## Documentation

- Overview: `docs/overview.md`
- CLI: `docs/cli.md`
- Configuration: `docs/configuration.md`
- Providers: `docs/providers.md`
- Roadmap: `ROADMAP.md`
- Changelog: `CHANGELOG.md`

- Configure `address` and `token` (token can be interpolated from env).
- Attempts KV v2 read; falls back to logical read. If the data is a map with a single string field, that value is used; otherwise the JSON-serialized map is returned.
- You can specify `extras.key` by adding `metadata: { key: someField }` in the config to extract a specific field from the secret.

## Exit Codes

- 0 success
- 2 configuration error
- 3 provider error (network, auth, etc.)
- 4 secret not found (or permission denied where indistinguishable)
- 5 command execution failure

## Security

- Secrets are never written to disk by this tool.
- Values are only present in memory and in the environment of the child process during `run`.
- Secret values are masked in dry-run output and logs by default.

## Contributing

See `.github/CONTRIBUTING.md`, `.github/CODE_OF_CONDUCT.md`, and `.github/SECURITY.md`.

### Linting locally

Install pre-commit and enable hooks:

```bash
pip install pre-commit
pre-commit install
```

Or run the lint script (auto-installs tools if missing):

```bash
bash scripts/lint.sh
```
