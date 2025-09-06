# Installation Guide

This guide covers how to install `skv` on different platforms.

## Quick Install (Recommended)

### Download Binary

Download the latest release for your platform:

```bash
# Linux x64
curl -sL "https://github.com/Amet13/skv/releases/latest/download/skv_linux_amd64" -o skv
chmod +x skv
sudo mv skv /usr/local/bin/

# Linux ARM64
curl -sL "https://github.com/Amet13/skv/releases/latest/download/skv_linux_arm64" -o skv
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
```

```powershell
# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/Amet13/skv/releases/latest/download/skv_windows_amd64.exe" -OutFile "skv.exe"
# Move to a directory in your PATH
```

### Verify Installation

```bash
skv version
```

## Build from Source

### Prerequisites

- Go 1.25+ installed
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

If you get a security warning on macOS:

```bash
sudo xattr -rd com.apple.quarantine /usr/local/bin/skv
```

### Permission Denied

If you get permission errors:

```bash
# Make binary executable
chmod +x skv

# Or install to user directory
mkdir -p ~/bin
mv skv ~/bin/
export PATH="$HOME/bin:$PATH"
```

### Command Not Found

If `skv` command is not found:

```bash
# Check if binary is in PATH
which skv

# Add to PATH if needed
export PATH="/usr/local/bin:$PATH"

# Or add to shell profile
echo 'export PATH="/usr/local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Version Mismatch

If you see an old version:

```bash
# Remove old binary
sudo rm /usr/local/bin/skv

# Download and install latest
curl -sL "https://github.com/Amet13/skv/releases/latest/download/skv_linux_amd64" -o skv
chmod +x skv
sudo mv skv /usr/local/bin/
```

## Next Steps

After installation:

1. **Generate config**: `skv init`
2. **Configure providers**: Edit `~/.skv.yaml`
3. **Validate setup**: `skv validate`
4. **Test secrets**: `skv get <alias>`

See [configuration guide](configuration.md) for detailed setup instructions.
