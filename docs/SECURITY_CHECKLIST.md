## Security Checklist

- Use least-privileged IAM roles for each provider.
- Prefer short-lived credentials where possible.
- Avoid printing secrets; use `--dry-run` to preview masked values.
- Keep scripts used by the `exec` provider audited and minimal.
- Restrict network access to secret backends as needed.
- Rotate secrets regularly and enforce versioning (e.g., AWS version stages).
- Log at `info` or lower; use `debug` only in trusted environments.
- Avoid writing secrets to disk; use `run` to inject into child process envs.
