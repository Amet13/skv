---
name: Bug report
about: Create a report to help us improve
title: "[BUG] "
labels: bug
assignees: ""
---

**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:

1. Create config file with '...'
2. Run command '....'
3. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Actual behavior**
A clear and concise description of what actually happened.

**Configuration**
Please provide your configuration file (with sensitive values redacted):

```yaml
# Your .skv.yaml configuration
secrets:
  - alias: example
    provider: aws
    name: redacted/secret/name
    # ... other config
```

**Command and output**

```bash
$ skv your-command --flags
# Paste the full command output here (redact any sensitive values)
```

**Environment (please complete the following information):**

- OS: [e.g. macOS 14.0, Ubuntu 22.04, Windows 11]
- skv version: [e.g. v1.0.0] (run `skv version`)
- Go version (if building from source): [e.g. 1.21.0]
- Provider: [e.g. AWS, GCP, Azure, Vault]

**Additional context**
Add any other context about the problem here.

**Logs**
If applicable, add logs with increased verbosity:

```bash
$ skv --log-level debug your-command
# Paste logs here (redact any sensitive values)
```

**Workaround**
If you found a workaround, please describe it here.

**Priority**

- [ ] Low - Minor inconvenience
- [ ] Medium - Affects functionality
- [ ] High - Blocks normal usage
- [ ] Critical - Security issue or data loss
