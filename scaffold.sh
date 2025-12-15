#!/bin/bash

# Go Platform Template Scaffolder
# Creates a new Go project from this template

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Helper functions
print_error() {
    echo -e "${RED}✗ Error: $1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}→ $1${NC}"
}

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
}

# Check if project name is provided
if [ -z "$1" ]; then
    print_header "Go Platform Template Scaffolder"
    echo ""
    echo -e "${YELLOW}Usage: $0 <project-name> [module-name]${NC}"
    echo ""
    echo "Arguments:"
    echo "  project-name    : Name of your project (lowercase, hyphens/underscores ok)"
    echo "  module-name     : Go module path (default: github.com/example/project-name)"
    echo ""
    echo "Examples:"
    echo "  $0 my-project"
    echo "  $0 my-project github.com/myorg/my-project"
    echo ""
    exit 1
fi

PROJECT_NAME="$1"
MODULE_NAME="${2:-github.com/example/$PROJECT_NAME}"

# Validate project name (lowercase, numbers, hyphens, underscores)
if [[ ! $PROJECT_NAME =~ ^[a-z0-9_-]+$ ]]; then
    print_error "Project name must contain only lowercase letters, numbers, hyphens, and underscores"
    exit 1
fi

# Validate module name
if [[ ! $MODULE_NAME =~ ^[a-z0-9./_-]+$ ]]; then
    print_error "Module name format invalid (use: domain.com/org/project)"
    exit 1
fi

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_DIR="$SCRIPT_DIR"
PROJECT_DIR="../$PROJECT_NAME"

# Check if directory already exists
if [ -d "$PROJECT_DIR" ]; then
    print_error "Directory '$PROJECT_NAME' already exists at $PROJECT_DIR"
    exit 1
fi

print_header "Go Platform Template Scaffolder"
echo ""
print_info "Project name: $PROJECT_NAME"
print_info "Module name: $MODULE_NAME"
print_info "Target directory: $PROJECT_DIR"
echo ""

# Create project directory
print_info "Creating project directory..."
mkdir -p "$PROJECT_DIR"
print_success "Project directory created"

# Copy root-level files
print_info "Copying template files..."
find "$TEMPLATE_DIR" -maxdepth 1 -type f \
    ! -name "scaffold.sh" \
    ! -name "server" \
    ! -name "go-platform-template" \
    ! -name "FINAL_AUDIT.md" \
    -exec cp {} "$PROJECT_DIR/" \;
print_success "Template files copied"

# Copy directories
print_info "Copying directories..."
for dir in cmd internal scripts docs docker-compose; do
    if [ -d "$TEMPLATE_DIR/$dir" ]; then
        cp -r "$TEMPLATE_DIR/$dir" "$PROJECT_DIR/"
    fi
done
print_success "Directories copied"

# Copy dotfiles
print_info "Copying dotfiles..."
cp "$TEMPLATE_DIR/.gitignore" "$PROJECT_DIR/" 2>/dev/null || true
cp "$TEMPLATE_DIR/.env.example" "$PROJECT_DIR/" 2>/dev/null || true
print_success "Dotfiles copied"

# Replace placeholders
print_info "Updating module names..."

TEMPLATE_MODULE="go_platform_template"
TEMPLATE_NAME="go-platform-template"
NEW_MODULE=$(echo "$MODULE_NAME" | sed 's/\//\\\//g')

# Replace in Go files
find "$PROJECT_DIR" -type f -name "*.go" | while read file; do
    sed -i "s|${TEMPLATE_MODULE}|${NEW_MODULE}|g" "$file"
    sed -i "s|${TEMPLATE_NAME}|${PROJECT_NAME}|g" "$file"
done

# Replace in Go mod file
if [ -f "$PROJECT_DIR/go.mod" ]; then
    sed -i "1s|^module.*|module $MODULE_NAME|" "$PROJECT_DIR/go.mod"
fi

# Replace in config files
for file in "$PROJECT_DIR"/Makefile "$PROJECT_DIR"/Dockerfile "$PROJECT_DIR"/docker-compose/*.yml; do
    if [ -f "$file" ]; then
        sed -i "s|${TEMPLATE_NAME}|${PROJECT_NAME}|g" "$file"
    fi
done

print_success "Module names updated"

# Initialize git
print_info "Initializing git repository..."
cd "$PROJECT_DIR"
rm -rf .git 2>/dev/null || true
git init > /dev/null 2>&1
git config user.email "dev@example.com" 2>/dev/null || true
git config user.name "Developer" 2>/dev/null || true
git add . > /dev/null 2>&1
git commit -m "Initial commit: created from go-platform-template" > /dev/null 2>&1
cd - > /dev/null
print_success "Git repository initialized"

echo ""
print_success "Project created successfully!"
echo ""

print_header "Next Steps"
echo ""
print_info "1. Navigate to project:"
echo "   cd $PROJECT_NAME"
echo ""
print_info "2. Configure environment:"
echo "   cp .env.example .env"
echo "   # Edit .env with your settings"
echo ""
print_info "3. Start development:"
echo "   make dev"
echo ""
print_info "4. Visit API documentation:"
echo "   http://localhost:8080/swagger/index.html"
echo ""

print_header "Project Structure"
echo ""
echo "cmd/server/               → Application entry point"
echo "internal/app/             → Application setup & bootstrap"
echo "internal/domain/          → Business logic (bounded contexts)"
echo "  ├── auth/               → Authentication & JWT"
echo "  ├── user/               → User management"
echo "  ├── file/               → File handling"
echo "  └── rbac/               → Access control"
echo "internal/platform/        → Infrastructure & external services"
echo "  ├── config/             → Configuration management"
echo "  ├── database/           → PostgreSQL & migrations"
echo "  ├── http/               → HTTP middleware"
echo "  ├── logger/             → Zap logger setup"
echo "  └── storage/            → MinIO integration"
echo "internal/shared/          → Cross-cutting concerns"
echo "  ├── errors/             → Error handling"
echo "  └── response/           → API response wrappers"
echo "docker-compose/           → Docker services"
echo ""

print_header "Useful Commands"
echo ""
echo "Development:"
echo "  make dev                → Start dev environment"
echo "  make dev-d              → Start in background"
echo "  make dev-down           → Stop dev environment"
echo "  make dev-logs           → View logs"
echo ""
echo "Building & Testing:"
echo "  make build              → Build binary"
echo "  make test               → Run tests"
echo "  make test-coverage      → Coverage report"
echo ""
echo "Code Quality:"
echo "  make fmt                → Format code"
echo "  make vet                → Run go vet"
echo "  make lint               → Run linter"
echo ""
echo "Dependencies:"
echo "  make deps               → Download dependencies"
echo "  make verify             → Verify dependencies"
echo "  make update-deps        → Tidy dependencies"
echo ""

print_header "Documentation"
echo ""
echo "→ API: http://localhost:8080/swagger/index.html"
echo "→ Health: http://localhost:8080/health"
echo ""

echo -e "${GREEN}Happy coding!${NC}"
