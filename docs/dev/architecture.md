## Provider extensibility

- The `internal/provider` package defines a minimal `Provider` interface and a simple registry.
- Each provider implements `FetchSecret(ctx, spec)` and is registered in `cmd/skv/main.go` with a unique name and any aliases.
- A provider receives a `SecretSpec` with a `Name` and `Extras` map for provider-specific options (e.g., region, profile, project, vault_url).
- To add a new provider:
  1. Create a new package under `internal/provider/<name>/` or add a file under an existing cloud namespace.
  2. Implement `Provider` and unit tests, using small seams to mock SDK clients.
  3. Register the provider name(s) in `cmd/skv/main.go`.
  4. Document usage and extras in `docs/providers.md`.

## Multiple secret backends per cloud

- AWS: Secrets Manager (`provider: aws`) and SSM Parameter Store (`provider: aws-ssm`).
- GCP: Google Secret Manager (`provider: gcp`).
- Azure: Key Vault (`provider: azure`).
- Vault: KV v2 and generic logical reads (`provider: vault`).
- Exec: external command (`provider: exec`).

Additional providers can follow the same pattern, using `Extras` to pass provider-specific parameters like profile, region, project, namespace, etc.

# Architecture

- Commands: Cobra-based CLI (`get`, `run`, `version`, `completion`).
- Config loader: YAML parser + env interpolation + schema validation.
- Provider registry: simple map of `name -> Provider` to allow extensions.
- Providers: each implements `FetchSecret(ctx, spec)`.
- Execution: `run` builds env and execs target with `os/exec`.

Extension points:

- Add a new provider under `internal/provider/<name>` and register it in `main.go`.
- Plugin loading via Go plugins or separate binaries can be considered for advanced use cases.

Security:

- No secret persistence to disk; masking when logging.
- No shell evaluation; arguments passed directly to `exec`.
