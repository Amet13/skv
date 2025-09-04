# CLI Reference

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

## skv list

List configured aliases. Use `-v/--verbose` to include provider and env name.

## skv export

Export selected secrets as shell `export VAR="value"` lines or `.env` style with `--env-file`.

## skv version

Print version info.

## skv completion [bash|zsh|fish|powershell]

Generate completion for shells.
