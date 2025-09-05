# Migration Guide

This guide helps you migrate from other secret management tools to `skv`.

## From AWS CLI / aws-vault

### Before (AWS CLI)

```bash
# Export secrets manually
export DB_PASSWORD=$(aws secretsmanager get-secret-value --secret-id myapp/prod/db_password --query SecretString --output text)
export API_KEY=$(aws ssm get-parameter --name /myapp/prod/api_key --with-decryption --query Parameter.Value --output text)

# Run application
./my-app
```

### After (skv)

```yaml
# ~/.skv.yaml
secrets:
  - alias: db_password
    provider: aws
    name: myapp/prod/db_password
    env: DB_PASSWORD
  - alias: api_key
    provider: aws-ssm
    name: /myapp/prod/api_key
    env: API_KEY
```

```bash
# Run application with secrets
skv run --all -- ./my-app
```

## From HashiCorp Vault CLI

### Before (Vault CLI)

```bash
# Login and export secrets
vault auth -method=userpass username=myuser
export DB_PASSWORD=$(vault kv get -field=password secret/myapp/db)
export API_TOKEN=$(vault kv get -field=token secret/myapp/api)

# Run application
./my-app
```

### After (skv)

```yaml
# ~/.skv.yaml
secrets:
  - alias: db_password
    provider: vault
    name: secret/data/myapp/db
    env: DB_PASSWORD
    extras:
      address: https://vault.company.com
      key: password
  - alias: api_token
    provider: vault
    name: secret/data/myapp/api
    env: API_TOKEN
    extras:
      address: https://vault.company.com
      key: token
```

```bash
# Run application with secrets
export VAULT_TOKEN=$(vault auth -method=userpass username=myuser -format=json | jq -r .auth.client_token)
skv run --all -- ./my-app
```

## From Google gcloud / berglas

### Before (gcloud)

```bash
# Export secrets manually
export DB_PASSWORD=$(gcloud secrets versions access latest --secret=db-password)
export API_KEY=$(gcloud secrets versions access latest --secret=api-key --project=my-project)

# Run application
./my-app
```

### After (skv)

```yaml
# ~/.skv.yaml
secrets:
  - alias: db_password
    provider: gcp
    name: projects/my-project/secrets/db-password/versions/latest
    env: DB_PASSWORD
  - alias: api_key
    provider: gcp
    name: projects/my-project/secrets/api-key/versions/latest
    env: API_KEY
```

```bash
# Run application with secrets
skv run --all -- ./my-app
```

## From Azure CLI

### Before (Azure CLI)

```bash
# Export secrets manually
export DB_PASSWORD=$(az keyvault secret show --vault-name MyVault --name db-password --query value -o tsv)
export API_KEY=$(az appconfig kv show --name MyConfig --key api-key --query value -o tsv)

# Run application
./my-app
```

### After (skv)

```yaml
# ~/.skv.yaml
secrets:
  - alias: db_password
    provider: azure
    name: db-password
    env: DB_PASSWORD
    extras:
      vault_url: https://MyVault.vault.azure.net
  - alias: api_key
    provider: azure-appconfig
    name: api-key
    env: API_KEY
    extras:
      endpoint: https://MyConfig.azconfig.io
```

```bash
# Run application with secrets
skv run --all -- ./my-app
```

## From Environment Files (.env)

### Before (.env files)

```bash
# .env
DB_PASSWORD=secret123
API_KEY=key456
DEBUG=true
```

```bash
# Load and run
source .env
./my-app
```

### After (skv)

```yaml
# ~/.skv.yaml
secrets:
  - alias: db_password
    provider: aws # or your provider
    name: myapp/prod/db_password
    env: DB_PASSWORD
  - alias: api_key
    provider: aws
    name: myapp/prod/api_key
    env: API_KEY
```

```bash
# Run application with secrets (plus regular env vars)
env DEBUG=true skv run --all -- ./my-app

# Or export to .env for compatibility
skv export --all --format env > .env
source .env
./my-app
```

## From Docker Secrets

### Before (Docker Compose)

```yaml
# docker-compose.yml
version: "3.8"
services:
  app:
    image: my-app
    secrets:
      - db_password
      - api_key
    environment:
      DB_PASSWORD_FILE: /run/secrets/db_password
      API_KEY_FILE: /run/secrets/api_key

secrets:
  db_password:
    external: true
  api_key:
    external: true
```

### After (skv in container)

```yaml
# docker-compose.yml
version: "3.8"
services:
  app:
    image: my-app
    volumes:
      - ~/.skv.yaml:/root/.skv.yaml:ro
    environment:
      - AWS_REGION
      - AWS_ACCESS_KEY_ID
      - AWS_SECRET_ACCESS_KEY
    command: ["skv", "run", "--all", "--", "./my-app"]
```

## From Kubernetes Secrets

### Before (Kubernetes)

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
data:
  db-password: <base64-encoded>
  api-key: <base64-encoded>
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: app
          envFrom:
            - secretRef:
                name: app-secrets
```

### After (skv in Kubernetes)

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: skv-config
data:
  skv.yaml: |
    secrets:
      - alias: db_password
        provider: aws
        name: myapp/prod/db_password
        env: DB_PASSWORD
      - alias: api_key
        provider: aws
        name: myapp/prod/api_key
        env: API_KEY
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      serviceAccountName: skv-service-account # with AWS IAM role
      containers:
        - name: app
          image: my-app
          command: ["skv", "run", "--all", "--", "./my-app"]
          volumeMounts:
            - name: config
              mountPath: /root/.skv.yaml
              subPath: skv.yaml
      volumes:
        - name: config
          configMap:
            name: skv-config
```

## From Custom Scripts

### Before (Custom bash script)

```bash
#!/bin/bash
# fetch-secrets.sh

DB_PASSWORD=$(some-custom-command --get db-password)
API_KEY=$(another-command --fetch api-key)

export DB_PASSWORD API_KEY
exec "$@"
```

### After (skv with exec provider)

```yaml
# ~/.skv.yaml
secrets:
  - alias: db_password
    provider: exec
    name: some-custom-command
    env: DB_PASSWORD
    extras:
      args: "--get db-password"
      trim: "true"
  - alias: api_key
    provider: exec
    name: another-command
    env: API_KEY
    extras:
      args: "--fetch api-key"
      trim: "true"
```

```bash
# Replace script usage
skv run --all -- ./my-app
```

## Migration Checklist

### 1. Inventory Current Secrets

- [ ] List all secrets currently used
- [ ] Identify their sources (AWS, Vault, files, etc.)
- [ ] Document environment variable names
- [ ] Note any transformation/processing needed

### 2. Create skv Configuration

- [ ] Choose appropriate providers
- [ ] Map secret names/paths
- [ ] Configure authentication (profiles, tokens, etc.)
- [ ] Set up environment variable names
- [ ] Test configuration with `skv list`

### 3. Test Secret Retrieval

- [ ] Test each secret individually: `skv get <alias>`
- [ ] Verify values match existing sources
- [ ] Test with dry-run: `skv run --all --dry-run -- env`
- [ ] Check authentication and permissions

### 4. Update Applications

- [ ] Replace custom secret-loading code
- [ ] Update deployment scripts/containers
- [ ] Modify CI/CD pipelines
- [ ] Update documentation

### 5. Gradual Migration

- [ ] Start with development/staging environments
- [ ] Run both systems in parallel initially
- [ ] Monitor for issues or missing secrets
- [ ] Migrate production last

### 6. Cleanup

- [ ] Remove old secret management code
- [ ] Clean up unused environment files
- [ ] Update team documentation
- [ ] Archive old configurations

## Common Migration Issues

### Authentication

**Problem**: Different authentication methods between tools
**Solution**: Use provider-specific authentication in `extras`:

```yaml
secrets:
  - alias: example
    provider: aws
    extras:
      profile: production # Use specific AWS profile
      region: us-east-1
```

### Secret Name Mapping

**Problem**: Different naming conventions
**Solution**: Use aliases to maintain compatibility:

```yaml
secrets:
  - alias: old_name_format # Keep old env var name
    provider: aws
    name: new/secret/path # New provider path
    env: OLD_NAME_FORMAT # Existing env var
```

### Missing Secrets

**Problem**: Some secrets don't exist in new provider
**Solution**: Use exec provider as bridge:

```yaml
secrets:
  - alias: legacy_secret
    provider: exec
    name: ./scripts/get-legacy-secret.sh
    env: LEGACY_SECRET
```

### Batch Operations

**Problem**: Need to process many secrets at once
**Solution**: Use `--all` flag and filtering:

```bash
# Get all secrets
skv run --all -- ./my-app

# Get specific subset
skv run --secrets db_password,api_key -- ./my-app

# Exclude certain secrets
skv run --all --all-except debug_flag,dev_token -- ./my-app
```

## Best Practices

1. **Start Small**: Migrate one application at a time
2. **Test Thoroughly**: Use dry-run mode extensively
3. **Keep Backups**: Maintain old configurations until confident
4. **Document Changes**: Update team documentation immediately
5. **Monitor**: Watch for authentication or permission issues
6. **Gradual Rollout**: Use feature flags or blue-green deployments
