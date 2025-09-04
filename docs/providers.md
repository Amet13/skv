## Providers

This page documents the supported providers, required environment variables, example `.skv.yaml` config snippets, and how to reference secrets.

### AWS Secrets Manager

- **Auth**: Default AWS credential chain (env vars, shared config, metadata/IMDS).
- **Name**: Secret name or ARN. Optional `version_stage` in `extras`.
- **Extras (optional)**:
  - `region`: AWS region (e.g., `us-east-1`)
  - `version_stage`: e.g., `AWSCURRENT`
- **Example**:

```yaml
secrets:
  - alias: db_password
    provider: aws
    name: myapp/prod/db_password
    env: DB_PASSWORD
    extras:
      region: us-east-1
```

### Google Secret Manager (GCP)

- **Auth**: Application Default Credentials. Set `GOOGLE_APPLICATION_CREDENTIALS` to a service account key if needed.
- **Name**: `projects/<PROJECT>/secrets/<SECRET>/versions/<VERSION>` (use `latest` for newest).
- **Example**:

```yaml
secrets:
  - alias: api_key
    provider: gcp
    name: projects/myproj/secrets/api_key/versions/latest
    env: API_KEY
```

### Azure Key Vault

- **Auth**: Default Azure credentials (managed identity, environment, or CLI login).
- **Name**: Secret name. Optional specific `version` via `extras.version`.
- **Extras (required)**:
  - `vault_url`: e.g., <https://myvault.vault.azure.net>
- **Example**:

```yaml
secrets:
  - alias: jwt_secret
    provider: azure
    name: jwt-secret
    env: JWT_SECRET
    extras:
      vault_url: https://myvault.vault.azure.net
```

### HashiCorp Vault (KV v2)

- **Auth**: `VAULT_TOKEN` (or token in `extras.token`); address from `VAULT_ADDR` or `extras.address`.
- **Name**: KV v2 path, typically `<mount>/data/<path>` (e.g., `kv/data/app/password`).
- **Extras (optional)**:
  - `address`: Vault address, e.g., <http://127.0.0.1:8200>
- **Example**:

```yaml
secrets:
  - alias: service_password
    provider: vault
    name: kv/data/myapp/password
    env: SERVICE_PASSWORD
    extras:
      address: http://127.0.0.1:8200
```

### Exec Provider

- **Behavior**: Executes a local command and uses stdout as the secret value.
- **Security**: Ensure the script is trusted and does not log secrets.
- **Name**: Command path when convenient; any additional arguments can go into `extras.args`.
- **Extras (optional)**:
  - `timeout`: duration like `5s`
- **Example**:

```yaml
secrets:
  - alias: dynamic_token
    provider: exec
    name: ./scripts/fetch_token.sh
    env: DYNAMIC_TOKEN
    extras:
      timeout: 5s
```
