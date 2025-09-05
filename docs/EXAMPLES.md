## Examples

### docker-compose

Use `skv run` to inject secrets into a service:

```yaml
services:
  app:
    image: myapp:latest
    entrypoint: ["sh", "-lc", "skv run --all -- -- /app"]
    volumes:
      - ~/.skv.yaml:/root/.skv.yaml:ro
```

### GitHub Actions

```yaml
- name: Inject secrets and run tests
  run: |
    curl -L -o /usr/local/bin/skv https://github.com/Amet13/skv/releases/download/vX.Y.Z/skv_linux_amd64
    chmod +x /usr/local/bin/skv
    echo "$SKV_CONFIG_CONTENT" > ~/.skv.yaml
    skv run --all -- -- go test ./...
```
