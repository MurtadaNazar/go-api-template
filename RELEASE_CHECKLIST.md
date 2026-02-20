# Release Checklist

Complete this before each release.

## Pre-Release (Local)

- [ ] Pull latest changes: `git pull origin main`
- [ ] All changes committed and pushed
- [ ] Run full quality check: `make check`
  - [ ] Code formatted: `make fmt`
  - [ ] Linting passes: `make lint`
  - [ ] Tests pass: `make test`
- [ ] Verify README is up-to-date
- [ ] Update CHANGELOG.md (if exists)
- [ ] Verify all files have proper permissions

## Create Release

Choose one method:

### Method 1: Using Make (Recommended)

```bash
make release-patch    # v1.0.0 → v1.0.1 (bug fixes)
make release-minor    # v1.0.0 → v1.1.0 (new features)
make release-major    # v1.0.0 → v2.0.0 (breaking changes)
```

### Method 2: Manual

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### Method 3: Using Script

```bash
./scripts/release.sh patch    # or minor/major
```

## Post-Release (GitHub Actions)

- [ ] Watch Actions tab for workflow runs
- [ ] release.yml completes successfully
- [ ] All platform binaries built (7 total)
- [ ] CHECKSUMS file created
- [ ] GitHub Release created with:
  - [ ] All 8 artifacts (7 binaries + CHECKSUMS)
  - [ ] Auto-generated release notes
  - [ ] Correct version tag

## Verification

Run on 1-2 platforms:

```bash
# Test installer
curl -fsSL https://raw.githubusercontent.com/murtadanazar/go-api-template/main/scripts/install.sh | bash

# Verify binary works
go-platform --version

# Check checksums
cd ~/.local/bin
sha256sum go-platform

# Verify against release
# Compare with value on GitHub Release page
```

## Announcement (Optional)

- [ ] Tweet/post on social media
- [ ] Update website/blog
- [ ] Notify users/subscribers

## Document Release

- [ ] Create/update CHANGELOG.md:
  ```markdown
  ## [1.0.0] - 2026-02-20
  
  ### Added
  - New features here
  
  ### Fixed
  - Bug fixes here
  
  ### Changed
  - Breaking changes (if major version)
  ```

## Troubleshooting

### Release failed, want to recreate tag?

```bash
# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push origin --delete v1.0.0

# Create again
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### Binaries not building?

- [ ] Check Actions tab for errors
- [ ] Verify Go version is correct
- [ ] Ensure all workflows have permission to write

### Installer downloading wrong binary?

- [ ] Check platform detection works: `./install.sh | head -20`
- [ ] Verify binary exists on GitHub Release
- [ ] Check filename matches URL pattern

## Commands Reference

```bash
# Build/Test
make build              # Build local binary
make build-all         # Build all platforms
make test              # Run tests
make lint              # Lint code
make fmt               # Format code
make check             # All quality checks

# Release
make release-patch     # Patch version
make release-minor     # Minor version
make release-major     # Major version

# Info
make version           # Show current version
```

## Tips

- Always run `make check` before releasing
- Tag only on main/master branch
- Use semantic versioning (MAJOR.MINOR.PATCH)
- Include descriptive release notes
- Test installer on multiple platforms before announcing

---

**Release takes 30 seconds, CI/CD handles the rest!** ⚡
