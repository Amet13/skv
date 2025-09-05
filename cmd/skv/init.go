package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var (
		provider   string
		outputFile string
		force      bool
	)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Generate a configuration template",
		Long: `Generate a configuration template for skv with examples for different providers.
This helps you get started quickly with a properly structured configuration file.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			if outputFile == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					outputFile = ".skv.yaml"
				} else {
					outputFile = filepath.Join(home, ".skv.yaml")
				}
			}

			// Check if file exists
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			template := generateTemplate(provider)

			err := os.WriteFile(outputFile, []byte(template), 0600)
			if err != nil {
				return fmt.Errorf("failed to write configuration file: %w", err)
			}

			fmt.Printf("Configuration template created: %s\n", outputFile)
			fmt.Println("📝 Edit the file to add your actual secret names and configuration.")
			fmt.Println("🔍 Use 'skv validate' to check your configuration.")
			fmt.Println("🚀 Use 'skv list' to see your configured secrets.")

			return nil
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "all", "Generate template for specific provider (all, aws, gcp, azure, vault, exec)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: ~/.skv.yaml)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing file")

	return cmd
}

func generateTemplate(provider string) string {
	switch provider {
	case "aws":
		return awsTemplate
	case "gcp":
		return gcpTemplate
	case "azure":
		return azureTemplate
	case "vault":
		return vaultTemplate
	case "exec":
		return execTemplate
	default:
		return fullTemplate
	}
}

const fullTemplate = `# skv configuration file
# See https://github.com/Amet13/skv for documentation

# Global defaults (optional)
defaults:
  region: us-east-1
  extras:
    version_stage: AWSCURRENT

secrets:
  # AWS Secrets Manager
  - alias: db_password
    provider: aws
    name: myapp/prod/db_password
    env: DB_PASSWORD
    extras:
      region: us-east-1
      version_stage: AWSCURRENT
      # profile: production  # Optional AWS profile

  # AWS SSM Parameter Store
  - alias: db_host
    provider: aws-ssm
    name: /myapp/prod/db_host
    env: DB_HOST
    extras:
      region: us-east-1
      with_decryption: "true"

  # Google Cloud Secret Manager
  - alias: api_key
    provider: gcp
    name: projects/my-project/secrets/api-key/versions/latest
    env: API_KEY
    extras:
      # credentials_file: /path/to/service-account.json  # Optional

  # Azure Key Vault
  - alias: jwt_secret
    provider: azure
    name: jwt-secret
    env: JWT_SECRET
    extras:
      vault_url: https://myvault.vault.azure.net
      # version: "specific-version"  # Optional

  # Azure App Configuration
  - alias: feature_flag
    provider: azure-appconfig
    name: myapp:feature:enabled
    env: FEATURE_ENABLED
    extras:
      endpoint: https://myconfig.azconfig.io
      label: prod

  # HashiCorp Vault (KV v2)
  - alias: service_password
    provider: vault
    name: kv/data/myapp/password
    env: SERVICE_PASSWORD
    extras:
      address: https://vault.company.com
      # Authentication via AppRole
      role_id: "{{ VAULT_ROLE_ID }}"
      secret_id: "{{ VAULT_SECRET_ID }}"
      # Or use token authentication
      # token: "{{ VAULT_TOKEN }}"
      # namespace: production  # Vault Enterprise

  # Custom script execution
  - alias: dynamic_token
    provider: exec
    name: ./scripts/get-token.sh
    env: DYNAMIC_TOKEN
    extras:
      args: "--environment prod"
      trim: "true"

# Tips:
# - Use {{ VAR }} for environment variable interpolation
# - Test with: skv validate --check-secrets
# - List secrets: skv list
# - Get single secret: skv get <alias>
# - Run with secrets: skv run --all -- your-command
# - Export to file: skv export --all --format env > .env
`

const awsTemplate = `# AWS-focused skv configuration
defaults:
  region: us-east-1
  extras:
    version_stage: AWSCURRENT

secrets:
  # AWS Secrets Manager
  - alias: db_password
    provider: aws
    name: myapp/prod/db_password
    env: DB_PASSWORD
    extras:
      region: us-east-1
      version_stage: AWSCURRENT
      profile: production  # Optional AWS profile

  # AWS SSM Parameter Store
  - alias: db_host
    provider: aws-ssm
    name: /myapp/prod/db_host
    env: DB_HOST
    extras:
      region: us-east-1
      with_decryption: "true"
      profile: production

  # AWS Secrets Manager (JSON secret with key selection)
  - alias: api_credentials
    provider: aws
    name: myapp/prod/api_credentials
    env: API_KEY
    extras:
      region: us-east-1
      # For JSON secrets, you might need to parse the JSON in your application
`

const gcpTemplate = `# Google Cloud-focused skv configuration
secrets:
  # Google Secret Manager
  - alias: db_password
    provider: gcp
    name: projects/my-project/secrets/db-password/versions/latest
    env: DB_PASSWORD
    extras:
      # credentials_file: /path/to/service-account.json  # Optional override

  - alias: api_key
    provider: gcp
    name: projects/my-project/secrets/api-key/versions/1  # Specific version
    env: API_KEY

  # Using project and version in extras instead of full name
  - alias: jwt_secret
    provider: gcp
    name: jwt-secret  # Short name when using project/version extras
    env: JWT_SECRET
    extras:
      project: my-project
      version: latest
`

const azureTemplate = `# Azure-focused skv configuration
secrets:
  # Azure Key Vault
  - alias: db_password
    provider: azure
    name: db-password
    env: DB_PASSWORD
    extras:
      vault_url: https://myvault.vault.azure.net
      # version: "specific-version-id"  # Optional

  # Azure App Configuration
  - alias: feature_flag
    provider: azure-appconfig
    name: myapp:feature:enabled
    env: FEATURE_ENABLED
    extras:
      endpoint: https://myconfig.azconfig.io
      label: prod  # Optional label

  - alias: connection_string
    provider: azure-appconfig
    name: myapp:database:connectionstring
    env: CONNECTION_STRING
    extras:
      endpoint: https://myconfig.azconfig.io
`

const vaultTemplate = `# HashiCorp Vault-focused skv configuration
defaults:
  extras:
    address: https://vault.company.com

secrets:
  # Vault KV v2 with AppRole authentication
  - alias: db_password
    provider: vault
    name: kv/data/myapp/database
    env: DB_PASSWORD
    extras:
      address: https://vault.company.com
      role_id: "{{ VAULT_ROLE_ID }}"
      secret_id: "{{ VAULT_SECRET_ID }}"
      key: password  # Extract specific key from JSON response

  # Vault KV v2 with token authentication
  - alias: api_key
    provider: vault
    name: kv/data/myapp/api
    env: API_KEY
    extras:
      address: https://vault.company.com
      token: "{{ VAULT_TOKEN }}"
      namespace: production  # Vault Enterprise namespace

  # Vault logical backend (non-KV)
  - alias: dynamic_secret
    provider: vault
    name: database/creds/myapp-role
    env: DYNAMIC_SECRET
    extras:
      address: https://vault.company.com
      token: "{{ VAULT_TOKEN }}"
`

const execTemplate = `# Exec provider-focused skv configuration
secrets:
  # Simple script execution
  - alias: current_token
    provider: exec
    name: ./scripts/get-current-token.sh
    env: CURRENT_TOKEN
    extras:
      trim: "true"

  # Script with arguments
  - alias: db_password
    provider: exec
    name: ./scripts/get-secret.sh
    env: DB_PASSWORD
    extras:
      args: "--secret-name db_password --environment prod"
      trim: "true"

  # Command with full path
  - alias: api_key
    provider: exec
    name: /usr/local/bin/get-api-key
    env: API_KEY
    extras:
      args: "--format json"
      trim: "true"

  # AWS CLI example (when skv AWS provider can't be used)
  - alias: custom_secret
    provider: exec
    name: aws
    env: CUSTOM_SECRET
    extras:
      args: "secretsmanager get-secret-value --secret-id my-secret --query SecretString --output text"
      trim: "true"
`

