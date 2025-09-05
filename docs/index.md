## skv Documentation

## Overview

skv is a cross-cloud CLI to fetch secrets from multiple providers and storages, and inject them into process environments, export them, or print to stdout.

- Unified config across clouds
- Safe injection with masking and dry-run
- Concurrency, retries, timeouts
- Extensible provider model

## Quickstart

1. Install: see README Installation.
2. Create config `~/.skv.yaml`:

```yaml
secrets:
  - alias: db_password
    provider: aws
    name: myapp/prod/db_password
    env: DB_PASSWORD
    extras:
      region: us-east-1
```

1. Use the CLI:

```bash
skv list
skv get db_password
skv export --all --env-file > .env
skv run --all -- -- env | grep DB_PASSWORD
```

## Reference

- CLI reference: `docs/cli.md`
- Configuration: `docs/configuration.md`
- Providers and storages: `docs/providers.md`
- Developer guide: `docs/DEVELOPING_PROVIDERS.md`
