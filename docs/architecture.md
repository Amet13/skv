# Architecture

- Commands: Cobra-based CLI (`get`, `run`, `version`, `completion`).
- Config loader: YAML parser + env interpolation + schema validation.
- Provider registry: simple map of `name -> Provider` to allow extensions.
- Providers: each implements `FetchSecret(ctx, spec)`.
- Execution: `run` builds env and execs target with `os/exec`.

Extension points:

- Add a new provider under `internal/provider/<name>` and register it in `main.go`.
- Consider future plugin loading via Go plugins or separate binaries.

Security:

- No secret persistence to disk; masking when logging.
- No shell evaluation; arguments passed directly to `exec`.
