# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

- TBD

## [0.2.1] - 2025-09-04

- Docs: expanded providers/configuration docs with detailed examples
- Examples: added `examples/` directory with per-provider `.skv.yaml` samples
- CI: removed Codecov upload; retained coverage check locally in workflow
- Release: fixed GoReleaser deprecation by switching `archives.builds` → `archives.ids`
- Tests: added Vault KVv2 empty payload test; CLI strict/non-strict partial scenarios
- Housekeeping: ignore `coverage.out` and `bin/` in `.gitignore`; removed committed tool binary

## [0.1.0] - 2025-09-04

- Initial MVP: CLI (get/run/version/completion/list/export)
- Providers: AWS Secrets Manager, HashiCorp Vault, GCP Secret Manager, Azure Key Vault, Exec provider
- Config: YAML with env interpolation and validation
- `run`: env injection with masking, strict/non-strict, timeout, concurrency
- CI: lint/tests/coverage upload; releases via GoReleaser
