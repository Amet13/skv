## Configuration

Default discovery order:

1. `--config` flag
2. `SKV_CONFIG` env var
3. `$HOME/.skv.yaml` or `$HOME/.skv.yml`

Indentation is 2 spaces.

### Minimal example

```yaml
secrets:
  - alias: db_password
    provider: aws
    name: myapp/prod/db-password
    env: DB_PASSWORD
```

### Defaults

You can set global defaults merged into each secret (per-secret values override):

```yaml
defaults:
  region: us-east-1
  extras:
    version_stage: AWSCURRENT
```

### Schema

```yaml
secrets:
  - alias: string # local name, used by -s/--secret and logs
    provider: string # aws | aws-ssm | gcp | azure | azure-appconfig | vault | exec
    name: string # provider-specific path/name
    env: string # environment variable name to export
    extras: # optional provider-specific parameters
      key: value
```

### Skeletons

```yaml
# AWS
secrets:
  - alias: db_pass
    provider: aws
    name: myapp/prod/db_password
    env: DB_PASSWORD
    extras:
      region: us-east-1
      version_stage: AWSCURRENT

  # GCP
  - alias: api_key
    provider: gcp
    name: projects/<PROJECT>/secrets/<SECRET>/versions/latest
    env: API_KEY

  # Azure Key Vault
  - alias: jwt_secret
    provider: azure
    name: jwt-secret
    env: JWT_SECRET
    extras:
      vault_url: <https://VAULT.vault.azure.net>

  # Azure App Configuration
  - alias: feature_flag
    provider: azure-appconfig
    name: myapp:feature:enabled
    env: FEATURE_ENABLED
    extras:
      endpoint: <https://CONFIG.azconfig.io>
      label: prod

  # Vault (KV v2)
  - alias: service_password
    provider: vault
    name: kv/data/myapp/password
    env: SERVICE_PASSWORD
    extras:
      address: http://127.0.0.1:8200
      # AppRole auth (optional)
      role_id: "{{ VAULT_ROLE_ID }}"
      secret_id: "{{ VAULT_SECRET_ID }}"

  # Exec
  - alias: dynamic_token
    provider: exec
    name: ./scripts/fetch_token.sh
    env: DYNAMIC_TOKEN
    extras:
      args: "--flag1 --flag2"
      trim: "true"
```

Notes:

- `{{ VAR }}` is interpolated from the environment; missing variables cause a load error.
- If `env` is omitted, the name is derived from alias in UPPER_SNAKE_CASE.
