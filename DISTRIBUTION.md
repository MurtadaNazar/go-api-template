# Distribution & Release Strategy

Professional distribution system for go-platform with multiple installation methods.

## Installation Methods

### 1. One-Line Install (Recommended for End Users)

```bash
curl -fsSL https://raw.githubusercontent.com/murtadanazar/go-api-template/main/scripts/install.sh | bash
```

**Features:**
- Auto-detects OS and architecture
- Downloads appropriate binary
- Installs to `~/.local/bin`
- Verifies SHA256 checksums
- Provides setup instructions

**Supported Platforms:**
- Linux (amd64, arm64, armv7)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

### 2. GitHub Releases

Direct download from [GitHub Releases page](https://github.com/murtadanazar/go-api-template/releases)

- Pre-built binaries for all platforms
- SHA256 checksums for verification
- Release notes with changelog
- Automated builds via CI/CD

### 3. Go Install

```bash
go install github.com/murtadanazar/go-api-template@latest
```

For Go developers. Requires Go 1.24+.

## Release Process

### 1. Trigger Release

```bash
# Create patch release (v1.0.0 → v1.0.1)
make release-patch

# Or minor (v1.0.0 → v1.1.0)
make release-minor

# Or major (v1.0.0 → v2.0.0)
make release-major
```

Or manually:
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### 2. GitHub Actions Workflows

#### `release.yml` - Custom Build Pipeline
Triggered on: `push` → `tags` matching `v*`

**What it does:**
- Builds binaries for all platforms (7 architectures)
- Generates SHA256 checksums
- Creates CHECKSUMS file
- Generates release notes from git log
- Creates GitHub Release with all artifacts

**Matrix:**
- Linux: amd64, arm64, armv7
- macOS: amd64, arm64
- Windows: amd64, arm64

#### `goreleaser.yml` - GoReleaser Pipeline
Triggered on: `push` → `tags` matching `v*`

**What it does:**
- Build with GoReleaser (professional release tool)
- Create archives (optional)
- Publish to multiple release platforms

**Benefits:**
- Standardized naming
- Better changelog formatting
- Supports multiple release platforms

#### `lint-test.yml` - Quality Gates
Triggered on: `push` to main/develop, `pull_request`

**What it does:**
- Run golangci-lint
- Execute all tests with coverage
- Build check on all platforms
- Upload coverage to Codecov

### 3. Verify Release

```bash
# Check GitHub Release created
curl -s https://api.github.com/repos/murtadanazar/go-api-template/releases/latest | jq .

# Verify checksums
sha256sum -c CHECKSUMS

# Test installation
curl -fsSL https://raw.githubusercontent.com/murtadanazar/go-api-template/main/scripts/install.sh | bash
```

## Release Checklist

- [ ] Ensure main/master is clean
- [ ] Update version numbers (if not automated)
- [ ] Review and merge all PRs
- [ ] Run `make check` (lint + test)
- [ ] Create git tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
- [ ] Push tag: `git push origin v1.0.0`
- [ ] GitHub Actions builds and releases
- [ ] Verify binaries on GitHub Releases page
- [ ] Test one-line installer on multiple platforms
- [ ] Announce release (if desired)

## Automation Files

### In `scripts/` Directory

#### `scripts/install.sh`
- Auto-detects platform
- Downloads correct binary
- Verifies checksums
- Sets up PATH
- Provides clear instructions

#### `scripts/release.sh`
Semantic versioning helper:
```bash
./scripts/release.sh [major|minor|patch]
```
Automatically:
- Parses current version from git tags
- Calculates new version
- Creates annotated git tag
- Pushes to origin



### In Root Directory

### Makefile
```bash
make build              # Build local binary
make build-all         # Build for all platforms
make release-patch     # Create patch release
make release-minor     # Create minor release
make release-major     # Create major release
```

Automatically:
- Parses current version from git tags
- Calculates new version
- Creates annotated git tag
- Pushes to origin
- GitHub Actions handles the rest

## Version Management

### Semantic Versioning

Format: `v{MAJOR}.{MINOR}.{PATCH}`

- **MAJOR** (v1.0.0): Breaking changes
- **MINOR** (v1.1.0): New features, backward compatible
- **PATCH** (v1.0.1): Bug fixes

Examples:
- `v0.1.0` - Initial release
- `v1.0.0` - First stable release
- `v1.1.0` - New features
- `v1.1.1` - Bug fix
- `v2.0.0` - Major breaking changes

### Tag Naming

- Use `v` prefix: `v1.0.0` (not `1.0.0`)
- Annotated tags: `git tag -a v1.0.0 -m "Release v1.0.0"`
- Never force-push release tags

## Artifact Structure

Each release includes:

```
go-platform-v1.0.0-linux-amd64          # 7 MB binary
go-platform-v1.0.0-linux-amd64.sha256   # 65 byte checksum
go-platform-v1.0.0-linux-arm64          # ...
go-platform-v1.0.0-linux-armv7          # ...
go-platform-v1.0.0-darwin-amd64         # ...
go-platform-v1.0.0-darwin-arm64         # ...
go-platform-v1.0.0-windows-amd64.exe    # ...
go-platform-v1.0.0-windows-arm64.exe    # ...
CHECKSUMS                               # All checksums
```

Total: 8 binaries + checksums

## Continuous Deployment

### What Happens After Tag Push

1. **GitHub detects tag** `v*` pattern
2. **CI/CD workflows trigger**:
   - `release.yml` - Builds all binaries
   - `goreleaser.yml` - Creates formatted release
3. **Artifacts uploaded** to GitHub Release
4. **Checksums verified** in CHECKSUMS file
5. **Release published** to public

No manual steps needed after tag push.

## End-User Experience

### Installation (30 seconds)

```bash
curl -fsSL https://raw.githubusercontent.com/murtadanazar/go-api-template/main/scripts/install.sh | bash
```

Outputs:
```
════════════════════════════════════════
  Go Platform Template - Installer
════════════════════════════════════════

→ Fetching latest release...
→ Platform: linux/amd64
→ Version: v1.0.0
→ Download URL: https://...
→ Downloading...
✓ Installation successful

════════════════════════════════════════
  Next Steps
════════════════════════════════════════

✓ ~/.local/bin is in your PATH

Start using:
  go-platform
```

### Verification

```bash
go-platform --version
# Output: go-platform version 1.0.0 (built 2026-02-20_12:30:00)
```

## Troubleshooting

### Release didn't trigger

- Check tag format: must be `v*` (e.g., `v1.0.0`)
- Verify GitHub Actions enabled in repo settings
- Check workflow permissions (must have write access)

### Installer fails

- Check network connectivity
- Verify GitHub release was created
- Ensure binaries uploaded successfully
- Check binary permissions

### Wrong platform downloaded

Run installer with debug:
```bash
bash -x install.sh
# Shows each step with variable values
```

## Security

### Checksum Verification

All binaries signed with SHA256:

```bash
sha256sum -c go-platform-v1.0.0-linux-amd64.sha256
```

### Future: Code Signing

For enhanced security, we could add:
- GPG signatures on releases
- Notarization (for macOS)
- Code signing (for Windows)

## Maintenance

### Keep Release Tools Updated

```bash
# Update GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Update GitHub Actions
# Pinned in workflows: v4, v5, etc.
```

### Monitor Release Health

- GitHub Actions tab - Watch build status
- Releases page - Verify artifacts uploaded
- Test installer monthly on multiple platforms

## References

- [GitHub Releases](https://github.com/murtadanazar/go-api-template/releases)
- [install.sh](install.sh)
- [Makefile](Makefile)
- [.goreleaser.yaml](.goreleaser.yaml)
- [.github/workflows/release.yml](.github/workflows/release.yml)
