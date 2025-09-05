# Project Conventions

## Naming

- Provider aliases (config `provider:`):
  - AWS Secrets Manager: `aws`
  - AWS SSM Parameter Store: `aws-ssm` (alias `ssm`)
  - Google Secret Manager: `gcp`
  - Azure Key Vault: `azure`
  - Azure App Configuration: `azure-appconfig` (alias `appconfig`)
  - HashiCorp Vault: `vault`
  - Exec: `exec`
- Config keys and extras use `snake_case` or lowercase hyphenated keys as documented per provider.
- Environment variable names are `UPPER_SNAKE_CASE`.
- CLI flags use `kebab-case`.

## Docs

- Top-level headings start with `#` and sub-sections use `##`, `###`.
- Use backticks for file paths, commands, and inline code.
- Keep examples minimal and verified.

## Code

- Go naming follows standard conventions; exported symbols have doc comments.
- Avoid panics; return errors with context using `fmt.Errorf("... %w ...")`.
- Tests must be hermetic; use seams/fakes for SDK clients.
