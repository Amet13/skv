## CLI Reference

## skv get <alias>

Fetch a single secret and print it to stdout.

Flags:

- `--config` path to config
- `--newline` append trailing newline
- `--raw` print raw value (default true)

## skv run [flags] -- <command> [args...]

Inject selected secrets into the command's environment.

Selection:

- `--all` inject all
- `--secrets` a,b,c or `-s` repeatable

Flags:

- `--dry-run` show env additions, masked
- `--strict` fail on missing (default true)
- `--mask` mask values in logs (default true)
- `--timeout` fetch timeout
- `--concurrency` number of concurrent provider calls (default 4)
- `--retries` number of retries; `--retry-delay` between retries (e.g., 200ms)
- `--require-env` ensure specific env names are present after fetch
- `--require-alias` ensure specific aliases are selected

## skv list

List configured aliases. Use `-v/--verbose` to include provider and env name.

## skv export

Export selected secrets as shell `export VAR="value"` lines or `.env` style with `--env-file`.

## skv version

Print version info.

## skv completion [bash|zsh|fish|powershell]

Generate completion for shells.

### Examples

```bash
# Fetch one secret
skv get db_password

# Export all secrets to .env format
skv export --all --env-file > .env

# Run with secrets from config, with concurrency
skv run --concurrency 8 -s db_password -- ./bin/app

# Exclude some aliases while using --all
skv run --all --all-except db_password,api_key -- -- printenv | grep -E 'JWT_SECRET|SERVICE_PASSWORD'

# Retries and timeouts
skv get db_password --retries 2 --retry-delay 300ms --timeout 5s
```
