# Providers

## AWS Secrets Manager

- Fields: `name` (secret ARN or name), `region`
- Not found maps to exit code `4`.
- Returns `SecretString` when available; otherwise `SecretBinary` converted to string.

## HashiCorp Vault

- Fields: `name` (path), `address`, `token`
- KV v2: `name` may be `<mount>/data/<path>` or set `metadata.mount`
- Extract a field with `metadata.key: <name>`
- For `token`, you can interpolate with `{{ VAULT_TOKEN }}` in config.
- Not found maps to exit code `4`.

## GCP Secret Manager

- Fields: `name` (resource or secret name), set `metadata.project` and optional `metadata.version` (default `latest`)
- Example:

```yaml
- alias: gcp-secret
  provider: gcp-secret-manager
  name: my-secret
  metadata:
      project: my-project
      version: latest
```

- Or full resource: `projects/123456/secrets/my-secret/versions/5`
- Not found and permission denied both map to exit code `4`.

## Azure Key Vault

- Fields: `name` (secret name), `metadata.vault_url` (e.g. <https://myvault.vault.azure.net/>), optional `metadata.version`
- Uses DefaultAzureCredential chain
- Not found and permission denied both map to exit code `4`.

## Exec Plugin Provider

- Execute an external command to fetch a secret value.
- Set `provider: exec` and specify extras:

- `metadata.cmd`: command to run
- `metadata.args`: optional space-separated arguments
- `metadata.trim`: "true" to trim whitespace from stdout
- The `name` field is appended as the last argument to the command.

Example:

```yaml
- alias: otp
  provider: exec
  name: account@example.com
  metadata:
      cmd: /usr/local/bin/secret-fetch
      args: "--issuer example"
      trim: "true"
```

Notes:

- Providers should avoid writing secret values to disk.
- Provider-specific metadata goes under `metadata:` and is passed as `extras` to the provider.
