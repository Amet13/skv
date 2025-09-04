## Security Policy

- Please report vulnerabilities privately by opening a security advisory or emailing the maintainers.
- Do not open public issues for sensitive reports.
- We aim to respond within 7 days.

### Handling Secrets

- Secrets are never written to disk by this tool.
- Values are only present in memory and in the environment of the child process during `run`.
- Secret values are masked in dry-run output and logs by default.
- Config files may contain environment interpolation tokens like `{{ VAR }}`; missing values will cause load to fail for safety.
- When fetching from providers, "not found" or permission-denied are surfaced as exit code `4` to allow clear CI handling.

### Supported Versions

- Main branch and the latest released version.
