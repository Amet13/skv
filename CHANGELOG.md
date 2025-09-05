# 📋 **Changelog**

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### ✨ **Added**

- 🎯 **New Commands**: `init`, `validate`, and `health` for better user experience
- 🧪 **Enhanced test coverage** across all providers (71.9% total coverage)
- 📈 **Comprehensive E2E and integration tests** with realistic scenarios
- 📦 **Enhanced GoReleaser** configuration with multi-platform support
- 🔧 **Shell completion** scripts generation for all major shells
- 📋 **SBOM generation** for supply chain security compliance
- ✍️ **Package signing** with Cosign for enhanced security
- 🍺 **Homebrew tap** support for easy macOS installation
- 🪣 **Scoop bucket** support for Windows package management
- 📦 **Linux package** distribution (deb, rpm, apk)
- 🎨 **Consistent emoji formatting** across all documentation and code
- 📁 **Standardized file naming** (lowercase with dashes for docs)

### 🔄 **Changed**

- 📚 **Updated documentation** with comprehensive examples and emoji formatting
- 🔧 **Improved error handling** and user experience across all commands
- 🚀 **Enhanced CI/CD pipeline** with security scanning and comprehensive checks
- 📝 **Restructured documentation** with consistent naming and better organization
- 🎯 **Optimized build process** with better caching and parallel execution

### 🧹 **Removed**

- 🗑️ **Unimplemented provider stubs** (OCI, IBM, Alibaba) for cleaner codebase
- 🐳 **Docker support** - removed to focus on core CLI functionality
- 📝 **Future functionality mentions** replaced with current capabilities
- 🧽 **Unused code** and outdated documentation sections

### 🐛 **Fixed**

- ✅ **Test coverage** for mock provider and version package (now 100%)
- 📖 **Documentation consistency** across all files with proper linking
- 🔗 **Broken documentation links** after file renaming
- 🔧 **Linting issues** and code quality improvements
- 🧪 **Test reliability** with better mocking and error handling

## [0.1.0] - 2024-01-01

### Added

- Initial release of skv CLI tool
- Support for multiple secret providers:
  - AWS Secrets Manager
  - AWS SSM Parameter Store
  - Google Cloud Secret Manager
  - Azure Key Vault
  - Azure App Configuration
  - HashiCorp Vault (KV v2)
  - Exec provider for custom scripts
- Configuration file support with environment variable interpolation
- Concurrent secret fetching with configurable concurrency
- Retry mechanism with configurable delays
- Dry-run mode for testing configurations
- Secret value masking in logs and output
- Shell completion support (bash, zsh, fish, powershell)
- Comprehensive CLI commands:
  - `get` - Fetch single secret
  - `run` - Inject secrets into process environment
  - `list` - List configured secrets
  - `export` - Export secrets in various formats
  - `version` - Show version information
- Cross-platform builds (Linux, macOS, Windows)
- Comprehensive test suite with good coverage
- Security-focused design with no disk writes
- Extensive documentation and examples

### Security

- Secrets never written to disk
- Values only present in memory and child process environment
- Secret masking in dry-run output and logs
- Secure credential handling for all providers
