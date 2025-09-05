## Providers

This page documents the supported providers, required environment variables, example `.skv.yaml` config snippets, and how to reference secrets.

Supported providers and storages:

- AWS Secrets Manager (`aws`)
- AWS SSM Parameter Store (`aws-ssm`, `ssm`)
- Google Secret Manager (`gcp`)
- Azure Key Vault (`azure`)
- Azure App Configuration (`azure-appconfig`, `appconfig`)
- HashiCorp Vault KV v2 / logical (`vault`)
- Exec command (`exec`)

Planned (stubs registered, not implemented yet):

- Oracle Cloud Infrastructure (OCI)
- IBM Cloud
- Alibaba Cloud

### AWS Secrets Manager

- Auth: Default AWS credential chain (env vars, shared config, metadata/IMDS).
- Name: Secret name or ARN. Optional `version_stage` in `extras`.
- Extras (optional):
  - `region`: AWS region (e.g., `us-east-1`)
  - `version_stage`: e.g., `AWSCURRENT`
  - `profile`: AWS shared config profile name (uses `~/.aws/config` and `~/.aws/credentials`)
- Example:

```yaml
secrets:
  - alias: db_password
    provider: aws
    name: myapp/prod/db_password
    env: DB_PASSWORD
    extras:
      region: us-east-1
      version_stage: AWSCURRENT
```

### Google Secret Manager (GCP)

- Auth: Application Default Credentials. Set `GOOGLE_APPLICATION_CREDENTIALS` to a service account key if needed.
- Name: `projects/<PROJECT>/secrets/<SECRET>/versions/<VERSION>` (use `latest` for newest).
- Extras (optional):
  - `credentials_file`: path to a service account JSON key file (overrides ADC)
- Example:

```yaml
secrets:
  - alias: api_key
    provider: gcp
    name: projects/myproj/secrets/api_key/versions/latest
    env: API_KEY
```

### Azure Key Vault

### Azure App Configuration (Parameter Store)

- Auth: Default Azure credentials (managed identity, environment, or CLI login).
- Name: Key name.
- Extras (required):
  - `endpoint`: e.g., https://<store>.azconfig.io
- Extras (optional):
  - `label`: App Configuration label to select
- Example:

```yaml
secrets:
  - alias: feature_flag
    provider: azure-appconfig
    name: myapp:feature:enabled
    env: FEATURE_ENABLED
    extras:
      endpoint: https://myapp.azconfig.io
      label: prod
```

- Auth: Default Azure credentials (managed identity, environment, or CLI login).
- Name: Secret name. Optional specific `version` via `extras.version`.
- Extras (required):
  - `vault_url`: e.g., <https://myvault.vault.azure.net>
- Example:

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

- Auth: `VAULT_TOKEN` or AppRole via `extras.role_id` and `extras.secret_id`.
- Address from `VAULT_ADDR` or `extras.address`.
- Name: KV v2 path, typically `<mount>/data/<path>` (e.g., `kv/data/app/password`).
- Extras (optional):
  - `address`: Vault address, e.g., <http://127.0.0.1:8200>
  - `mount`: override KV mount (if not inferrable)
  - `key`: preferred field name inside secret data
  - `namespace`: Vault Enterprise namespace to use
- Example:

```yaml
secrets:
  - alias: service_password
    provider: vault
    name: kv/data/myapp/password
    env: SERVICE_PASSWORD
    extras:
      address: http://127.0.0.1:8200
      role_id: "{{ VAULT_ROLE_ID }}"
      secret_id: "{{ VAULT_SECRET_ID }}"
```

### Exec Provider

### AWS SSM Parameter Store

- Auth: Default AWS credential chain and profiles (`AWS_PROFILE`).
- Name: Parameter name (e.g., `/myapp/prod/db_password`).
- Extras (optional):
  - `region`: AWS region
  - `profile`: shared config profile
  - `with_decryption`: `true|false` (default true)
- Example:

```yaml
secrets:
  - alias: db_password
    provider: aws-ssm
    name: /myapp/prod/db_password
    env: DB_PASSWORD
    extras:
      region: us-east-1
      with_decryption: "true"
```

- Behavior: Executes a local command and uses stdout as the secret value.
- If `extras.cmd` is omitted, `name` is treated as the command.
- Security: Ensure the script is trusted and does not log secrets.
- Name: Command path when convenient; any additional arguments can go into `extras.args`.
- Extras (optional):
  - `args`: space-separated command arguments
  - `trim`: `true` to trim whitespace from stdout
- Example:

```yaml
secrets:
  - alias: dynamic_token
    provider: exec
    name: ./scripts/fetch_token.sh
    env: DYNAMIC_TOKEN
    extras:
      args: "--timeout 5s"
      trim: "true"
```
