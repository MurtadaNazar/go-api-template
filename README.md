# Go Platform Template - TUI Scaffolder

A production-ready Go API scaffolder with interactive TUI. Similar to Laravel's installer, this tool creates complete Go projects from templates with feature selection.

## Quick Start (30 seconds)

### Install

```bash
# One-line install (recommended)
curl -fsSL https://raw.githubusercontent.com/murtadanazar/go-api-template/main/scripts/install.sh | bash

# Or with Go
go install github.com/murtadanazar/go-api-template@latest

# Then run
go-platform
```

### Local Development

```bash
# Clone repo
git clone https://github.com/murtadanazar/go-api-template.git
cd go-api-template

# Build and run
make build
./go-platform

# Or run directly
go run .
```

## What It Does

The TUI creates production-ready Go API projects with:

- ✅ **Authentication** - JWT with token rotation
- ✅ **User Management** - CRUD & RBAC
- ✅ **Database** - PostgreSQL with GORM
- ✅ **File Storage** - MinIO S3-compatible
- ✅ **API Docs** - Auto-generated Swagger
- ✅ **Docker** - Docker & Docker Compose
- ✅ **Logging** - Structured logging (Zap)
- ✅ **Project Structure** - Clean architecture

## Workflow

```
$ go run .

┌─────────────────────────────────────┐
│ Go Platform Template                │
│ Production-Ready Go API Framework    │
│                                     │
│ ▶ Create New Project                │
│   Create project from template      │
│ - Help                              │
│   View keyboard shortcuts           │
│ - Exit                              │
│   Exit the scaffolder               │
└─────────────────────────────────────┘

↑/↓ Navigate • ENTER Select • CTRL+C Quit
```

### 1. Select Features

```
Choose features to include:
▶ ✓ Authentication (JWT)
  ✓ User Management (requires Auth)
  [ ] Database
  [ ] File Storage (requires Database)
  [ ] API Docs
  [ ] Docker

Dependencies auto-managed!
```

### 2. Enter Project Details

```
Project Name: my-awesome-api
Module: github.com/myorg/my-awesome-api
```

### 3. Confirm & Create

Project created in parent directory with only selected features.

```
$ cd ../my-awesome-api
$ go run ./cmd/server
# API running at http://localhost:8080
```

## Generated Project Structure

```
my-awesome-api/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── app/
│   │   ├── routes.go        # (feature-generated)
│   │   ├── middleware.go
│   │   └── bootstrap.go
│   ├── domain/
│   │   ├── auth/            # (if selected)
│   │   ├── user/            # (if selected)
│   │   └── file/            # (if selected)
│   └── platform/
│       ├── config/
│       ├── logger/
│       ├── database/
│       └── http/
├── Makefile                 # Build/dev commands
├── Dockerfile               # (if Docker selected)
├── docker-compose.yml       # (if Docker selected)
├── go.mod
└── .env.example
```

## Development (Template)

### Setup

```bash
# Install dependencies
make install-deps

# Or manually
go mod download
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Build & Run

```bash
# Build
make build
./go-platform

# Run directly
go run .

# Build for all platforms
make build-all
```

### Quality Checks

```bash
# Run all checks (format, lint, test)
make check

# Or individually
make fmt      # Format code
make lint     # Lint code
make test     # Run tests
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| ↑ / ↓ | Navigate |
| SPACE | Toggle feature |
| ENTER | Select/Proceed |
| TAB | Switch fields |
| CTRL+C | Cancel/Exit |

## Features

### Smart Dependencies

- **User Management** requires **Authentication** (auto-enabled)
- **File Storage** requires **Database** (auto-enabled)
- Deselecting a requirement auto-disables dependents
- TUI warns about missing dependencies

### Feature Details

#### Authentication (JWT)
- RS256 token signing
- Access & refresh tokens
- Token rotation
- Secure password hashing (bcrypt)

#### User Management
- User CRUD operations
- Role-based access control
- Admin user support
- Pagination & filtering

#### Database
- PostgreSQL integration
- GORM ORM
- Connection pooling
- Migrations support

#### File Storage
- MinIO S3-compatible
- Per-user isolation
- Metadata tracking
- Secure operations

#### API Docs
- Swagger/OpenAPI 3.0
- Auto-generated from comments
- Live at `/swagger/index.html`
- Easy to document

#### Docker
- Dockerfile for API
- docker-compose.yml for services
- PostgreSQL container
- MinIO container

## Created Project Usage

```bash
cd ../my-project

# Setup
cp .env.example .env
# Edit .env with your settings

# Start development (with Docker)
make dev-d

# View API
http://localhost:8080/swagger/index.html

# Run tests
make test

# Build for production
make build
```

## Keyboard Navigation

### Main Menu
- ↑/↓ - Navigate items
- ENTER - Select
- CTRL+C - Exit

### Text Input
- Type normally
- Backspace/Delete - Remove characters
- Tab - Next field
- ENTER - Submit

### Feature Selection
- ↑/↓ - Navigate features
- SPACE - Toggle selection
- ENTER - Confirm
- CTRL+C - Cancel

## Terminal Requirements

- Width: 60+ columns
- Height: 20+ lines
- Color support (16 or 256 colors)
- UTF-8 support

Works on:
- macOS (iTerm2, Terminal.app, VS Code)
- Linux (GNOME, Konsole, Kitty, Alacritty, VS Code)
- Windows (Windows Terminal, VS Code)
- tmux, screen, WSL

## Build Output

```
Binary: ~7 MB (statically linked)
Platform: Cross-platform (Linux, macOS, Windows)
Go: 1.24+
Startup: < 100ms
```

## Architecture

```
go-platform-template/
├── main.go                  # TUI entry point
├── internal/
│   └── scaffold/            # Scaffolding logic
│       ├── model.go         # TUI state machine
│       ├── theme.go         # Terminal colors
│       └── processor.go     # Project creation
├── scaffold/                # Project templates
│   ├── base/                # Base files
│   └── features/            # Feature templates
└── docs/                    # Documentation
```

## Extending the TUI

To add new menu items or features:

1. **Edit** `internal/scaffold/model.go`
2. **Add** menu item to `NewModel()`
3. **Add** handler in `Update()`
4. **Add** view function for rendering

Each feature is self-contained in `scaffold/features/`.

## Examples

### Create Minimal Project

```
Run TUI
Select: Authentication (JWT) only
Enter: my-jwt-api, github.com/me/my-jwt-api
Result: Small, focused project
```

### Create Full-Stack Project

```
Run TUI
Select: All features
Enter: my-full-app, github.com/me/my-full-app
Result: Complete production-ready API
```

### Create API + Docs

```
Run TUI
Select: Auth, User, Database, File Storage, API Docs
Enter: my-rest-api, github.com/me/my-rest-api
Result: Fully documented REST API
```

## Installation

For detailed installation instructions, see [INSTALL.md](INSTALL.md).

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/murtadanazar/go-api-template/main/install.sh | bash
```

Supports Linux, macOS, and Windows on amd64, arm64, and armv7.

## Releases & Updates

- **GitHub Releases**: [Download binaries](https://github.com/murtadanazar/go-api-template/releases)
- **CI/CD**: Automated builds for all platforms on every release tag
- **Checksums**: SHA256 verification provided for all binaries
- **Auto-Updates**: Easy upgrade path via installer script

## Support

### Troubleshooting

**Installation issues**: See [INSTALL.md](INSTALL.md)

**TUI won't run:**
```bash
# Check Go version (if building from source)
go version  # Should be 1.24+

# Run directly
go run .
```

**Project creation fails:**
- Ensure parent directory exists and is writable
- Project name must be lowercase with hyphens
- Module format: `github.com/org/project`

### Documentation

- **[INSTALL.md](INSTALL.md)** - Installation guide
- **[README.md](README.md)** - Feature overview & usage
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Contributing guide
- **[docs/](docs/)** - Complete reference

## License

MIT - See LICENSE file

---

**Go Platform Template - TUI Scaffolder**  
**Similar to Laravel Installer**  
**Production-Ready Projects in Seconds**

```bash
go run .
```
