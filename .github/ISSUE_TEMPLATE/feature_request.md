---
name: Feature request
about: Suggest an idea for this project
title: "[FEATURE] "
labels: enhancement
assignees: ""
---

**Is your feature request related to a problem? Please describe.**
A clear and concise description of what the problem is. Ex. I'm always frustrated when [...]

**Describe the solution you'd like**
A clear and concise description of what you want to happen.

**Describe alternatives you've considered**
A clear and concise description of any alternative solutions or features you've considered.

**Provider-specific requests**
If this is related to a specific provider (AWS, GCP, Azure, Vault, etc.), please specify:

- Provider name:
- Specific service/feature:
- Use case:

**Additional context**
Add any other context, screenshots, or examples about the feature request here.

**Configuration example (if applicable)**

```yaml
# Example of how you'd like to configure this feature
secrets:
  - alias: example
    provider: new-provider
    name: example-secret
    extras:
      new_option: value
```

**Priority**

- [ ] Low - Nice to have
- [ ] Medium - Would be helpful
- [ ] High - Blocking my use case
- [ ] Critical - Security or data loss concern
