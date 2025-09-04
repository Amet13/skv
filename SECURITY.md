# Security Policy

## Reporting a Vulnerability

Please do not open public issues for security reports.

Email the maintainers or use GitHub Security Advisories to privately disclose.

We aim to respond within 72 hours.

## Notes

- Secrets are never written to disk by `skv`.
- Values are only present in memory and the child process environment when using `run`.
- Avoid sharing debug logs that might contain secret values.
