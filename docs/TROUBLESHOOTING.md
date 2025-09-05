## Troubleshooting

### Common issues

- Unknown provider

  - Ensure the `provider` matches a registered name (`aws`, `aws-ssm`, `gcp`, `azure`, `vault`, `exec`).
  - Verify providers are registered in `cmd/skv/main.go`.

- Auth/Permissions

  - AWS: check `AWS_PROFILE`, `AWS_REGION` or `extras.region`, and IAM policies.
  - GCP: ensure ADC is configured or set `extras.credentials_file`.
  - Azure: login (CLI/env) and set `extras.vault_url`.
  - Vault: set `VAULT_ADDR`/`extras.address` and token/AppRole.

- Not found errors

  - Double-check secret/parameter names and versions.
  - For Vault KV v2, ensure path format `<mount>/data/<path>`.

- CLI errors
  - Use `--log-level debug` to get more details.
  - Add `--retries` and `--retry-delay` for transient failures.

### Getting help

- Run with `--dry-run` to preview env additions.
- Use `skv list -v` to verify aliases and env mappings.
