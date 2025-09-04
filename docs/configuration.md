# Configuration

Default path: `$HOME/.skv.yaml`. Override with `SKV_CONFIG` or `--config`.

```yaml
secrets:
    - alias: db-password
      provider: aws-secrets-manager
      name: myapp/prod/db-password
      region: us-east-1

    - alias: api-token
      provider: vault
      name: secret/data/tokens/myapp
      address: https://vault.example.com
      token: "{{ VAULT_TOKEN }}"
      env: MYAPP_API_TOKEN
```

- `{{ VAR }}` is interpolated from the environment; missing variables cause a load error.
- If `env` is omitted, the variable name is derived from the alias in UPPER_SNAKE_CASE.
