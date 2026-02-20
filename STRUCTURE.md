# Project Structure

Clean, unified, professional Go project structure.

## Root Level (Essentials Only)

```
go-platform-template/
├── 📖 Documentation
│   ├── README.md                  # Main docs & features
│   ├── INSTALL.md                 # Installation guide
│   ├── CONTRIBUTING.md            # How to contribute
│   ├── SECURITY.md                # Security policy
│   ├── CODE_OF_CONDUCT.md         # Community guidelines
│   └── LICENSE                    # MIT License
│
├── ⚙️  Configuration
│   ├── Makefile                   # Build automation (USE THIS)
│   ├── go.mod                     # Go dependencies
│   ├── go.sum                     # Dependency checksums
│   ├── .golangci.yml              # Linter rules
│   ├── .editorconfig              # Editor formatting
│   ├── .env.example               # Config template
│   └── .goreleaser.yaml           # Release tool config
│
├── 📋 Release & Distribution
│   ├── DISTRIBUTION.md            # Release strategy
│   └── RELEASE_CHECKLIST.md       # Step-by-step release
│
├── 🔧 Source Code (Core)
│   ├── main.go                    # Entry point
│   └── internal/                  # Application code
│       ├── app/                   # API & routes
│       ├── domain/                # Domain models
│       ├── platform/              # Services
│       ├── scaffold/              # Scaffolding engine
│       ├── shared/                # Utilities
│       └── testutil/              # Test helpers
│
├── 📦 Deployment
│   ├── scaffold/                  # Project templates
│   │   ├── base/                  # Base files
│   │   └── features/              # Feature modules
│   └── docs/                      # Documentation
│
└── 🚀 Scripts & Automation
    └── scripts/                   # ALL SCRIPTS HERE
        ├── install.sh             # One-line installer
        └── release.sh             # Version tagging helper
```

## What Goes Where

### Root Level (12 files)
✅ **Keep in root:**
- Makefile - Standard Go practice
- go.mod, go.sum - Standard Go
- main.go - Entry point
- Documentation files (README, INSTALL, etc.)
- Configuration files (go.sum, .env.example, etc.)
- License file

❌ **Do NOT keep in root:**
- Binary files (.exe, compiled binaries)
- Temp/cache files
- Old archived documentation
- Multiple redundant docs

### `scripts/` Directory
✅ **ALL scripts go here:**
- `install.sh` - Installation (called from README)
- `release.sh` - Version management
- Any future helper scripts

### `internal/` Directory
✅ **Core application code:**
- Business logic
- Domain models
- Platform services
- Scaffolding engine
- No user-facing code here

### `scaffold/` Directory
✅ **Project templates:**
- Base project templates
- Feature modules
- Generated project examples

### `docs/` Directory
✅ **API docs, guides (optional):**
- Swagger/OpenAPI specs
- Architecture diagrams
- Additional documentation
- Currently empty, ready for content

## Command Reference

```bash
# Development
make build              # Build local binary
make check             # Lint + test + format
make test              # Run tests
make lint              # Code linting
make fmt               # Format code

# Release
make release-patch     # v1.0.0 → v1.0.1
make release-minor     # v1.0.0 → v1.1.0
make release-major     # v1.0.0 → v2.0.0

# Or use script
./scripts/release.sh patch    # Same result

# Installation (for users)
curl -fsSL https://raw.githubusercontent.com/murtadanazar/go-api-template/main/scripts/install.sh | bash
```

## CI/CD Pipeline

```
.github/workflows/
├── release.yml        # Builds 7 platforms on tag
├── goreleaser.yml     # Professional releases
├── lint-test.yml      # Quality gates on push/PR
└── ci.yml             # Existing CI
```

Triggered by:
- `release.yml` & `goreleaser.yml` → git tag `v*`
- `lint-test.yml` → push to main/develop
- `ci.yml` → push/PR

## Why This Structure

| Component | Location | Why |
|-----------|----------|-----|
| Makefile | Root | Standard Go practice |
| Scripts | `scripts/` | Organized, findable, not cluttering root |
| Code | `internal/` | Go convention, only exported types |
| Docs | Root + `docs/` | README in root (visible on GitHub), detailed in docs/ |
| Config | Root | Standard Go project layout |

## File Counts

- **Root:** 12 files (clean)
- **scripts/:** 2 files (all helpers)
- **Source:** internal/ + main.go
- **Config:** 4 hidden files (.env.example, .golangci.yml, etc.)
- **Total:** 18 essential files, zero clutter

## What NOT To Keep

❌ Binary artifacts (go-platform, scaffold-tui)
❌ Temporary files (CLEAN_STATUS.md, PROJECT_STRUCTURE.txt)
❌ Old docs (GETTING_STARTED.md, START_HERE.md)
❌ Archived folders (docs/ARCHIVE)
❌ Redundant build scripts (build.sh)
❌ Generated files (go-platform, *.exe)

## Next Steps

1. **Commit:** `git add -A && git commit -m "chore: unified project structure"`
2. **Build:** `make build`
3. **Release:** `git tag -a v0.1.0 -m "Release v0.1.0" && git push origin v0.1.0`
4. **Done:** GitHub Actions builds for all platforms

---

**Structure complete. Ready to ship.** ✨
