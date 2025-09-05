## Description

Brief description of what this PR does.

Fixes #(issue)

## Type of Change

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Refactoring (no functional changes)
- [ ] Performance improvement
- [ ] Test coverage improvement

## Provider Changes

If this PR affects a specific provider, please specify:

- [ ] AWS (Secrets Manager / SSM)
- [ ] GCP (Secret Manager)
- [ ] Azure (Key Vault / App Configuration)
- [ ] HashiCorp Vault
- [ ] Exec provider
- [ ] New provider: \***\*\_\_\_\*\***

## Testing

- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] I have tested this change manually with real providers (if applicable)

### Test Coverage

- [ ] Coverage remains above 70%
- [ ] All new code is covered by tests

## Configuration Changes

If this PR changes configuration format or adds new options:

- [ ] Updated documentation in `docs/configuration.md`
- [ ] Updated example in `skv-example.yml`
- [ ] Backward compatibility maintained (or breaking change noted above)

## Documentation

- [ ] Updated relevant documentation
- [ ] Updated provider documentation in `docs/providers.md` (if applicable)
- [ ] Updated CLI documentation in `docs/cli.md` (if applicable)
- [ ] Added examples to `docs/examples.md` (if applicable)

## Security

- [ ] No secrets or credentials are logged or exposed
- [ ] Changes follow the security principles in `docs/security-checklist.md`
- [ ] No new security vulnerabilities introduced

## Checklist

- [ ] My code follows the Go style guidelines
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] My changes generate no new warnings
- [ ] I have run `make lint` and `make test` locally
- [ ] Any dependent changes have been merged and published

## Breaking Changes

If this is a breaking change, describe what breaks and how users should migrate:

```
# Before
skv old-command --old-flag

# After
skv new-command --new-flag
```

## Additional Notes

Any additional information, deployment notes, etc.
