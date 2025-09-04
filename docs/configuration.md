## Configuration

Default path: `$HOME/.skv.yaml`. Override with `SKV_CONFIG` or `--config`.

Indentation is 2 spaces.

### Minimal example

```yaml
secrets:
  - alias: db_password
    provider: aws
    name: myapp/prod/db-password
    env: DB_PASSWORD
```

### Schema

```yaml
secrets:
  - alias: string # local name, used by -s/--secret and logs
    provider: string # aws | gcp | azure | vault | exec
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

  # GCP
  - alias: api_key
    provider: gcp
    name: projects/<PROJECT>/secrets/<SECRET>/versions/latest
    env: API_KEY

  # Azure
  - alias: jwt_secret
    provider: azure
    name: jwt-secret
    env: JWT_SECRET
    extras:
      vault_url: https://<VAULT>.vault.azure.net

  # Vault (KV v2)
  - alias: service_password
    provider: vault
    name: kv/data/myapp/password
    env: SERVICE_PASSWORD
    extras:
      address: http://127.0.0.1:8200

  # Exec
  - alias: dynamic_token
    provider: exec
    name: ./scripts/fetch_token.sh
    env: DYNAMIC_TOKEN
```

Notes:

- `{{ VAR }}` is interpolated from the environment; missing variables cause a load error.
- If `env` is omitted, the name is derived from alias in UPPER_SNAKE_CASE.
