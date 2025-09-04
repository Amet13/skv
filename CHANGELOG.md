# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

- TBD

## [0.1.0] - 2025-09-04

- Initial MVP: CLI (get/run/version/completion/list/export)
- Providers: AWS Secrets Manager, HashiCorp Vault, GCP Secret Manager, Azure Key Vault, Exec provider
- Config: YAML with env interpolation and validation
- `run`: env injection with masking, strict/non-strict, timeout, concurrency
- CI: lint/tests/coverage upload; releases via GoReleaser
