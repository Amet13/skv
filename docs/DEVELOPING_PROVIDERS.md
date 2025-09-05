## Developing a new provider

This guide explains how to add a new secret provider (e.g., another cloud KMS/secret store).

### 1) Understand the interface

- Implement `provider.Provider` with a single method:

```go
type Provider interface {
    FetchSecret(ctx context.Context, spec provider.SecretSpec) (string, error)
}
```

- You receive `SecretSpec{Name, Extras, Alias, EnvName, Provider}`. Use `spec.Name` as the resource identifier and `spec.Extras` for provider-specific options (e.g., region, profile, project, namespace).

### 2) Create the implementation

- Put your code under `internal/provider/<yourprovider>/` or inside an existing cloud namespace if it shares SDK dependencies.
- Add small seams to make SDK clients replaceable in tests (function variables or small interfaces).
- Map SDK-specific "not found" errors to `provider.ErrNotFound` to enable consistent CLI behavior.

### 3) Register in the CLI

- Register the provider and any friendly aliases in `cmd/skv/main.go`:

```go
provider.Register("your-provider", yourpkg.New())
provider.Register("alias", yourpkg.New())
```

### 4) Document

- Update `docs/providers.md` with:
  - Auth expectations
  - Name format (resource path)
  - Supported `extras`
  - A minimal YAML example

### 5) Tests and coverage

- Add unit tests exercising success and not-found paths.
- Use seams to inject fake clients; avoid network calls.
- Run `make test` and ensure overall coverage > 70% with `make cover`.

### 6) Lint and build

- `make fmt && make lint && make build`

### 7) Example configuration

- Update `example-skv.yml` if appropriate, to showcase usage.
