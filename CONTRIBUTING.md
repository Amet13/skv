## Contributing to skv

Thanks for your interest in contributing!

### Development setup

- Go 1.25+
- Make

Commands:

```bash
make fmt       # format Go, shell, markdown, YAML
make lint      # run linters and tests
make test      # run tests
make cover     # coverage with gate
make build     # build host binary
make build-all # build small cross matrix
```

### Conventional commits

Use conventional commits for clear history:
`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`.

### Pull requests

- Open a PR with a clear description
- Ensure CI is green
- Update docs when behavior changes

### Code of Conduct

This project follows the Code of Conduct in `CODE_OF_CONDUCT.md`.
