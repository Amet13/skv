## Contributing to skv

Thanks for your interest in contributing! Please follow these guidelines:

- Use Go 1.25.x or newer.
- Run `./build.sh host` before submitting PRs.
- Add tests for new features and bug fixes.
- Keep code readable and well-structured; avoid unnecessary complexity.

### Development

1. Fork and clone the repo.
1. Create a feature branch.
1. Build, lint, and run tests:

- `./build.sh host`
- `./scripts/lint.sh` (auto-installs tools if missing)

1. Open a PR with a clear description.

### Code Style

- Follow standard Go formatting and idioms.
- Avoid logging secret values; always mask or omit them.
- Keep provider errors specific; map not-found to exit code 4.
