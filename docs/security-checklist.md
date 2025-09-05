# Security Checklist

This document provides security best practices and guidelines for using `skv` safely in production environments.

## Core Security Principles

### Secret Handling

- **Never write secrets to disk** - `skv` keeps secrets in memory only
- **Mask secrets in output** - Use `--mask` flag (enabled by default) for dry-run and logs
- **Inject into process environment** - Use `skv run` instead of exporting to shell
- **Fail fast on missing secrets** - Use `--strict` mode (default) to catch issues early

### Authentication & Authorization

- **Use least-privileged IAM roles** for each provider
- **Prefer short-lived credentials** where possible (STS tokens, service accounts)
- **Rotate credentials regularly** and enforce versioning (e.g., AWS version stages)
- **Use provider-specific authentication** (profiles, roles, service accounts)

### Network Security

- **Restrict network access** to secret backends as needed
- **Use TLS/HTTPS endpoints** for all provider communications
- **Validate TLS certificates** - avoid `--insecure` flags in production
- **Use VPC endpoints** when available (AWS PrivateLink, etc.)

### Configuration Security

- **Protect configuration files** with appropriate file permissions (600/640)
- **Use environment interpolation** `{{ VAR }}` instead of hardcoding credentials
- **Validate configuration** before deployment
- **Version control configuration** (without secrets) for auditability

### Logging & Monitoring

- **Log at `info` or lower** in production; use `debug` only in trusted environments
- **Monitor secret access patterns** for anomalies
- **Set up alerting** for failed secret retrievals
- **Audit configuration changes** and access patterns

### Exec Provider Security

- **Audit scripts thoroughly** used by the `exec` provider
- **Keep scripts minimal** and focused on secret retrieval only
- **Use absolute paths** for script execution
- **Validate script permissions** and ownership
- **Avoid shell injection** - use proper argument passing

## Provider-Specific Security

### AWS

- Use IAM roles with minimal permissions
- Enable CloudTrail logging for API calls
- Use VPC endpoints for Secrets Manager/SSM
- Enable secret rotation where possible
- Use specific version stages (not `AWSPENDING`)

```yaml
# Secure AWS configuration
secrets:
  - alias: db_password
    provider: aws
    name: myapp/prod/db_password
    extras:
      region: us-east-1
      version_stage: AWSCURRENT # Explicit version
      profile: production # Specific profile
```

### Google Cloud

- Use service accounts with minimal permissions
- Enable audit logging for Secret Manager
- Use Workload Identity in GKE
- Specify exact secret versions when possible

```yaml
# Secure GCP configuration
secrets:
  - alias: api_key
    provider: gcp
    name: projects/my-project/secrets/api-key/versions/1 # Specific version
    extras:
      credentials_file: /path/to/service-account.json # If needed
```

### Azure

- Use managed identities when possible
- Enable Key Vault logging and monitoring
- Use specific secret versions
- Configure network access restrictions

```yaml
# Secure Azure configuration
secrets:
  - alias: jwt_secret
    provider: azure
    name: jwt-secret
    extras:
      vault_url: <https://myvault.vault.azure.net>
      version: "specific-version-id" # Pin to specific version
```

### HashiCorp Vault

- Use AppRole authentication for services
- Enable audit logging
- Use short-lived tokens
- Configure proper ACL policies
- Use TLS for all communications

```yaml
# Secure Vault configuration
secrets:
  - alias: service_password
    provider: vault
    name: kv/data/myapp/password
    extras:
      address: <https://vault.company.com> # HTTPS only
      role_id: "{{ VAULT_ROLE_ID }}" # From environment
      secret_id: "{{ VAULT_SECRET_ID }}" # From environment
      namespace: production # Vault Enterprise
```

## Deployment Security

### Container Security

```dockerfile
# Secure container practices
FROM alpine:latest
RUN adduser -D -s /bin/sh skv
USER skv
COPY --chown=skv:skv skv /usr/local/bin/
COPY --chown=skv:skv --chmod=600 .skv.yaml /home/skv/
ENTRYPOINT ["skv", "run", "--all", "--"]
```

### Kubernetes Security

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: skv-service-account
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::ACCOUNT:role/skv-role
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      serviceAccountName: skv-service-account
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
        fsGroup: 1000
      containers:
        - name: app
          image: my-app
          command: ["skv", "run", "--all", "--", "./my-app"]
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities:
              drop: ["ALL"]
```

### CI/CD Security

```yaml
# GitHub Actions example
steps:
  - name: Configure AWS credentials
    uses: aws-actions/configure-aws-credentials@v4
    with:
      role-to-assume: arn:aws:iam::ACCOUNT:role/github-actions-role
      role-session-name: skv-deployment

  - name: Deploy with secrets
    run: |
      skv run --all --timeout 30s -- ./deploy.sh
    env:
      SKV_CONFIG: ./deployment/.skv.yaml
```

## Security Incident Response

### If Secrets Are Compromised

1. **Immediately rotate affected secrets** in the provider
2. **Revoke access** for compromised credentials
3. **Update configuration** to use new secret versions
4. **Audit access logs** to understand the scope
5. **Review and update** security policies

### If Configuration Is Exposed

1. **Rotate any hardcoded credentials** (if any)
2. **Review environment variables** that might be exposed
3. **Update authentication tokens** used for interpolation
4. **Audit recent access** to secret providers

## Security Audit Checklist

### Pre-Production

- [ ] Configuration files have proper permissions (600/640)
- [ ] No hardcoded secrets in configuration
- [ ] All providers use least-privilege authentication
- [ ] TLS/HTTPS enabled for all provider endpoints
- [ ] Logging configured appropriately for environment
- [ ] Scripts used by exec provider are audited and minimal
- [ ] Timeout and retry settings are reasonable
- [ ] Dry-run testing completed successfully

### Production Deployment

- [ ] Service accounts/IAM roles have minimal required permissions
- [ ] Network access is restricted (VPC endpoints, firewalls)
- [ ] Monitoring and alerting configured for failures
- [ ] Secret rotation schedules established
- [ ] Incident response procedures documented
- [ ] Regular security reviews scheduled

### Ongoing Maintenance

- [ ] Regular credential rotation
- [ ] Monitoring secret access patterns
- [ ] Updating skv to latest versions
- [ ] Reviewing and updating IAM policies
- [ ] Auditing configuration changes
- [ ] Testing disaster recovery procedures

## Security Testing

### Test Commands

```bash
# Test configuration without executing
skv run --all --dry-run -- env

# Test with timeout to prevent hanging
skv get secret-alias --timeout 10s

# Test with strict mode to catch missing secrets
skv run --all --strict -- echo "test"

# Verify masking works
skv run --all --dry-run --mask -- echo "test"
```

### Automated Security Testing

```bash
# In CI/CD pipeline
skv list --format json | jq '.[] | select(.provider == "exec") | .name' | xargs -I {} echo "Audit script: {}"
skv run --all --dry-run --timeout 5s -- echo "Security test passed"
```

Remember: Security is an ongoing process, not a one-time setup. Regularly review and update these practices as your environment and requirements change.
