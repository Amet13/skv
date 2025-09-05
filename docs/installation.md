# Installation Guide

This guide covers various ways to install `skv` on different platforms.

## Quick Install (Recommended)

### Download Binary

Download the latest release for your platform:

```bash
# Linux x64
curl -sL "https://github.com/Amet13/skv/releases/latest/download/skv_linux_amd64" -o skv
chmod +x skv
sudo mv skv /usr/local/bin/

# macOS (Intel)
curl -sL "https://github.com/Amet13/skv/releases/latest/download/skv_darwin_amd64" -o skv
chmod +x skv
sudo mv skv /usr/local/bin/
sudo xattr -rd com.apple.quarantine /usr/local/bin/skv

# macOS (Apple Silicon)
curl -sL "https://github.com/Amet13/skv/releases/latest/download/skv_darwin_arm64" -o skv
chmod +x skv
sudo mv skv /usr/local/bin/
sudo xattr -rd com.apple.quarantine /usr/local/bin/skv

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/Amet13/skv/releases/latest/download/skv_windows_amd64.exe" -OutFile "skv.exe"
```

### Verify Installation

```bash
skv version
```

## Package Managers

### Homebrew (macOS/Linux)

```bash
# Add tap
brew tap amet13/tap

# Install
brew install skv

# Update
brew upgrade skv
```

### Scoop (Windows)

```bash
# Add bucket
scoop bucket add amet13 https://github.com/Amet13/scoop-bucket

# Install
scoop install skv

# Update
scoop update skv
```

### Linux Package Managers

#### Debian/Ubuntu (APT)

```bash
# Download and install .deb package
curl -sLO "https://github.com/Amet13/skv/releases/latest/download/skv_linux_amd64.deb"
sudo dpkg -i skv_linux_amd64.deb

# Or install dependencies if needed
sudo apt-get install -f
```

#### RHEL/CentOS/Fedora (YUM/DNF)

```bash
# Download and install .rpm package
curl -sLO "https://github.com/Amet13/skv/releases/latest/download/skv_linux_amd64.rpm"
sudo rpm -i skv_linux_amd64.rpm

# Or with yum/dnf
sudo yum install ./skv_linux_amd64.rpm
sudo dnf install ./skv_linux_amd64.rpm
```

#### Alpine Linux (APK)

```bash
# Download and install .apk package
curl -sLO "https://github.com/Amet13/skv/releases/latest/download/skv_linux_amd64.apk"
sudo apk add --allow-untrusted ./skv_linux_amd64.apk
```

## Build from Source

### Prerequisites

- Go 1.21+ installed
- Git

### Steps

```bash
# Clone repository
git clone https://github.com/Amet13/skv.git
cd skv

# Build for your platform
make build

# Or build for all supported platforms
make build-all

# Install to /usr/local/bin
sudo cp dist/skv_$(go env GOOS)_$(go env GOARCH) /usr/local/bin/skv
```

### Development Build

```bash
# Install directly with go
go install github.com/Amet13/skv/cmd/skv@latest

# Rename if needed
mv $GOPATH/bin/skv $GOPATH/bin/skv
```

## Container Images

### Docker

```bash
# Run with Docker (coming soon)
docker run --rm -v ~/.skv.yaml:/root/.skv.yaml ghcr.io/amet13/skv:latest version
```

### Kubernetes

```yaml
# Example Kubernetes Job (coming soon)
apiVersion: batch/v1
kind: Job
metadata:
  name: skv-secrets
spec:
  template:
    spec:
      containers:
        - name: skv
          image: ghcr.io/amet13/skv:latest
          command: ["skv", "run", "--all", "--", "my-app"]
          volumeMounts:
            - name: config
              mountPath: /root/.skv.yaml
              subPath: skv.yaml
      volumes:
        - name: config
          configMap:
            name: skv-config
      restartPolicy: Never
```

## Shell Completions

After installation, enable shell completions:

### Bash

```bash
# Install completion
skv completion bash | sudo tee /etc/bash_completion.d/skv > /dev/null

# Or for user only
mkdir -p ~/.local/share/bash-completion/completions
skv completion bash > ~/.local/share/bash-completion/completions/skv
```

### Zsh

```bash
# Install completion
mkdir -p ~/.zfunc
skv completion zsh > ~/.zfunc/_skv

# Add to .zshrc if not already present
echo 'fpath+=(~/.zfunc)' >> ~/.zshrc
echo 'autoload -Uz compinit && compinit' >> ~/.zshrc
```

### Fish

```bash
# Install completion
skv completion fish > ~/.config/fish/completions/skv.fish
```

### PowerShell

```powershell
# Add to PowerShell profile
skv completion powershell | Out-String | Invoke-Expression

# Or save to profile
skv completion powershell >> $PROFILE
```

## Verification

After installation, verify everything works:

```bash
# Check version
skv version

# Check help
skv --help

# Test with a simple config
echo 'secrets:
  - alias: test
    provider: exec
    name: echo "hello world"
    env: TEST_SECRET' > test.yaml

skv --config test.yaml get test
# Should output: hello world

# Clean up
rm test.yaml
```

## Troubleshooting

### macOS Quarantine

If you get a "cannot be opened because the developer cannot be verified" error:

```bash
sudo xattr -rd com.apple.quarantine /usr/local/bin/skv
```

### Permission Denied

If you get permission denied errors:

```bash
chmod +x skv
```

### Command Not Found

Ensure the binary is in your PATH:

```bash
echo $PATH
which skv

# Add to PATH if needed (add to ~/.bashrc or ~/.zshrc)
export PATH="/usr/local/bin:$PATH"
```

### Version Mismatch

If you have multiple installations:

```bash
# Find all installations
which -a skv

# Remove old versions
sudo rm /old/path/to/skv
```

## Next Steps

- [Configuration Guide](configuration.md)
- [Provider Setup](providers.md)
- [CLI Reference](cli.md)
- [Examples](EXAMPLES.md)
