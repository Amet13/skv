# skv Overview

skv is a fast, single-binary CLI that fetches secrets from multiple providers and injects them into process environments or prints them safely.

- Simple config in `$HOME/.skv.yaml`
- Providers: AWS Secrets Manager, HashiCorp Vault, GCP Secret Manager, Azure Key Vault
- Secure by default: no disk writes, masking, no shell expansion

See `docs/configuration.md`, `docs/cli.md`, and `docs/providers.md`.
