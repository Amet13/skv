# skv Documentation

## Overview

`skv` is a cross-cloud CLI to fetch secrets from multiple providers and storages, and inject them into process environments, export them, or print to stdout.

- Unified config across clouds
- Safe injection with masking and dry-run
- Concurrency, retries, timeouts
- Extensible provider model

## Quickstart

1. **Install**: see [installation.md](installation.md)
2. **Generate config template**: `skv init`
3. **Edit** `~/.skv.yaml`:

```yaml
secrets:
  - alias: db_password
    provider: aws
    name: myapp/prod/db_password
    env: DB_PASSWORD
    extras:
      region: us-east-1
```

4. **Validate and use**:

```bash
skv validate                           # Check configuration
skv list                              # Show configured secrets
skv get db_password                   # Get single secret
skv export --all --format env > .env # Export to file
skv run --all -- env | grep DB_PASSWORD  # Run with secrets
```

## Reference

- **CLI reference**: [`cli.md`](cli.md)
- **Configuration**: [`configuration.md`](configuration.md)
- **Providers**: [`providers.md`](providers.md)
- **Developer guide**: [`dev/developing-providers.md`](dev/developing-providers.md)
- **Troubleshooting**: [`troubleshooting.md`](troubleshooting.md)
- **Security**: [`security-checklist.md`](security-checklist.md)
- **Examples**: [`examples.md`](examples.md)
- **Migration**: [`migration.md`](migration.md)
