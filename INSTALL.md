# Installation Guide

Go Platform Template provides multiple installation methods for different use cases.

## Quick Install (Recommended)

### One-Line Install

```bash
curl -fsSL https://raw.githubusercontent.com/murtadanazar/go-api-template/main/scripts/install.sh | bash
```

The installer will:
- Detect your OS and architecture
- Download the latest compatible binary
- Place it in `~/.local/bin`
- Provide setup instructions for PATH

### From GitHub Releases

Download the binary for your platform from [GitHub Releases](https://github.com/murtadanazar/go-api-template/releases):

```bash
# Linux amd64
wget https://github.com/murtadanazar/go-api-template/releases/download/v1.0.0/go-platform-v1.0.0-linux-amd64
chmod +x go-platform-v1.0.0-linux-amd64
mv go-platform-v1.0.0-linux-amd64 ~/.local/bin/go-platform

# macOS arm64 (Apple Silicon)
wget https://github.com/murtadanazar/go-api-template/releases/download/v1.0.0/go-platform-v1.0.0-darwin-arm64
chmod +x go-platform-v1.0.0-darwin-arm64
mv go-platform-v1.0.0-darwin-arm64 ~/.local/bin/go-platform

# Windows (PowerShell)
# Download manually from releases page or use curl
curl -L https://github.com/murtadanazar/go-api-template/releases/download/v1.0.0/go-platform-v1.0.0-windows-amd64.exe -o go-platform.exe
```



### Using Go Install

```bash
go install github.com/murtadanazar/go-api-template@latest
```

## Verify Installation

```bash
go-platform --version
# or
go-platform version
```

## Supported Platforms

| OS | Architecture | Status |
|---|---|---|
| Linux | amd64 | ✅ Supported |
| Linux | arm64 | ✅ Supported |
| Linux | armv7 | ✅ Supported |
| macOS | amd64 | ✅ Supported |
| macOS | arm64 | ✅ Supported |
| Windows | amd64 | ✅ Supported |
| Windows | arm64 | ✅ Supported |

## Building from Source

### Prerequisites

- Go 1.24 or higher
- Git

### Clone & Build

```bash
git clone https://github.com/murtadanazar/go-api-template.git
cd go-api-template

# Build
make build

# Or with the build script
./build.sh go-platform

# Run
./go-platform
```

### Build for Specific Platform

```bash
# Linux
make build-linux-amd64
make build-linux-arm64

# macOS
make build-darwin-amd64
make build-darwin-arm64

# Windows
make build-windows-amd64

# All platforms
make build-all
```

## Verify Binary Integrity

Each release includes SHA256 checksums for verification:

```bash
# Download binary and checksum
wget https://github.com/murtadanazar/go-api-template/releases/download/v1.0.0/go-platform-v1.0.0-linux-amd64
wget https://github.com/murtadanazar/go-api-template/releases/download/v1.0.0/go-platform-v1.0.0-linux-amd64.sha256

# Verify
sha256sum -c go-platform-v1.0.0-linux-amd64.sha256
# Should output: go-platform-v1.0.0-linux-amd64: OK
```

## PATH Setup

If the installer couldn't automatically add `~/.local/bin` to your PATH, add it manually:

### Bash (~/.bashrc or ~/.bash_profile)
```bash
export PATH="$HOME/.local/bin:$PATH"
```

### Zsh (~/.zshrc)
```bash
export PATH="$HOME/.local/bin:$PATH"
```

### Fish (~/.config/fish/config.fish)
```fish
set -gx PATH $HOME/.local/bin $PATH
```

Then reload your shell:
```bash
source ~/.bashrc  # or ~/.zshrc, etc.
```

## Troubleshooting

### "command not found: go-platform"

1. Verify installation:
   ```bash
   ls -la ~/.local/bin/go-platform
   ```

2. Check PATH:
   ```bash
   echo $PATH
   ```

3. If `~/.local/bin` is not in PATH, add it as shown above

### Wrong binary downloaded

1. Check your OS and architecture:
   ```bash
   uname -s  # OS
   uname -m  # Architecture
   ```

2. Download the correct binary from [GitHub Releases](https://github.com/murtadanazar/go-api-template/releases)

### Permission denied

Make the binary executable:
```bash
chmod +x /path/to/go-platform
```

## Uninstall

```bash
rm -f ~/.local/bin/go-platform
# or
rm -f $(which go-platform)


```

## Update

To update to the latest version:

```bash
# Using installer script
curl -fsSL https://raw.githubusercontent.com/murtadanazar/go-api-template/main/install.sh | bash

# Using Go
go install github.com/murtadanazar/go-api-template@latest

# Or download from GitHub Releases
```

## Support

- **Issues**: [GitHub Issues](https://github.com/murtadanazar/go-api-template/issues)
- **Documentation**: [README.md](README.md)
- **Quick Start**: [README.md#quick-start](README.md#quick-start)
