# Contributing Guide

Thank you for your interest in contributing to the Go Platform Template! This guide will help you get started.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Coding Standards](#coding-standards)
- [Git Workflow](#git-workflow)
- [Pull Request Process](#pull-request-process)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)

## Code of Conduct

Be respectful, inclusive, and professional. We appreciate all contributions regardless of experience level.

## Getting Started

### Fork and Clone

```bash
# Fork the repository on GitHub

# Clone your fork
git clone https://github.com/YOUR_USERNAME/go-platform-template.git
cd go-platform-template

# Add upstream remote
git remote add upstream https://github.com/ORIGINAL_OWNER/go-platform-template.git
```

### Create a Feature Branch

```bash
# Update main from upstream
git fetch upstream
git checkout main
git merge upstream/main

# Create feature branch
git checkout -b feature/your-feature-name
```

## Development Setup

### Prerequisites

- Go 1.25.1+
- PostgreSQL 12+
- Docker & Docker Compose (recommended)
- Make

### Initial Setup

```bash
# Install dependencies
make deps

# Copy environment file
cp .env.example .env

# Start development environment
make dev-d

# Run the application
make run
```

## Making Changes

### Understanding the Project Structure

Familiarize yourself with the architecture:

```
internal/domain/    → Business logic organized by domain
internal/app/       → Application bootstrap
internal/platform/  → Cross-cutting concerns
cmd/server/         → Application entry point
```

See `ARCHITECTURE.md` for detailed information.

### Adding a New Feature

#### Example: Adding a New Endpoint

**1. Create the Handler**

File: `internal/domain/{domain}/api/handler.go`

```go
// @Summary Get resource
// @Description Get resource by ID
// @Tags Resources
// @Param id path string true "Resource ID"
// @Success 200 {object} model.Resource
// @Failure 404 {object} response.ErrorResponse
// @Router /resources/{id} [get]
func (h *ResourceHandler) GetResource(c *gin.Context) {
    id := c.Param("id")
    
    resource, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        _ = c.Error(err)
        return
    }
    
    c.JSON(http.StatusOK, resource)
}
```

**2. Add Service Logic**

File: `internal/domain/{domain}/service/service.go`

```go
func (s *ResourceService) GetByID(ctx context.Context, id string) (*Resource, error) {
    if id == "" {
        return nil, apperrors.NewAppError(
            apperrors.BadRequestError,
            "Resource ID is required",
        )
    }
    
    var resource Resource
    if err := s.db.WithContext(ctx).First(&resource, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, apperrors.NewAppError(
                apperrors.NotFoundError,
                "Resource not found",
            )
        }
        return nil, apperrors.NewAppError(
            apperrors.InternalError,
            "Failed to fetch resource",
        )
    }
    
    return &resource, nil
}
```

**3. Define Models**

File: `internal/domain/{domain}/model/model.go`

```go
type Resource struct {
    ID        string    `gorm:"primaryKey" json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

**4. Register Route**

File: `internal/app/routes.go`

```go
// In RegisterRoutes function
resourceHandler := handlers.NewResourceHandler(db, logger)
resourceGroup.GET("/:id", resourceHandler.GetResource)
```

**5. Add Swagger Annotations**

Swagger docs auto-regenerate when you save with:

```bash
GIN_MODE=debug go run ./cmd/server/main.go
```

## Coding Standards

### Go Conventions

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting: `make fmt`
- Use `go vet` for error checking: `make vet`

### File Naming

```
handler.go          → HTTP handlers
service.go          → Business logic
model.go            → Data structures
response.go         → Response types
request.go          → Request DTOs
repository.go       → Data access (when needed)
```

### Package Organization

```
internal/domain/auth/
├── api/
│   └── handler.go           # HTTP handlers
├── model/
│   ├── model.go             # Domain entities
│   ├── request.go           # Request DTOs
│   └── response.go          # Response types
└── service/
    └── service.go           # Business logic
```

### Naming Conventions

```go
// Interfaces (add -er suffix)
type UserService interface {}
type Validator interface {}

// Constants (CamelCase)
const DefaultPageSize = 20

// Unexported (lowercase)
func (h *Handler) handleError(err error) {}

// Exported (Capitalize)
func (h *Handler) HandleError(err error) {}
```

### Function Structure

```go
// Order: validation → logic → response
func (s *Service) DoSomething(ctx context.Context, req *Request) (*Response, error) {
    // 1. Validate
    if req.Name == "" {
        return nil, apperrors.NewAppError(...)
    }
    
    // 2. Execute logic
    result, err := s.process(ctx, req)
    if err != nil {
        return nil, err
    }
    
    // 3. Return response
    return &Response{Data: result}, nil
}
```

### Error Handling

```go
// Always check errors
if err != nil {
    // Handle or return
    return apperrors.NewAppError(...)
}

// Use type assertions for custom errors
if appErr, ok := apperrors.IsAppError(err); ok {
    // Handle application error
}

// Wrap errors with context when needed
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

### Comments

```go
// Package-level comment (required for exported packages)
// Package auth provides authentication and authorization functionality.
package auth

// Exported function must have comment (go vet requirement)
// Validate checks if the user credentials are valid.
func (s *Service) Validate(ctx context.Context, credentials *Credentials) error {
    // Implementation detail comments
    // explain WHY, not WHAT (code shows what)
}
```

### Type Definitions

```go
type User struct {
    // Public fields
    ID        string    `gorm:"primaryKey" json:"id"`
    Email     string    `gorm:"uniqueIndex" json:"email"`
    Username  string    `json:"username"`
    
    // Password excluded from JSON
    Password  string    `gorm:"column:password_hash" json:"-"`
    
    // Timestamps
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## Git Workflow

### Commit Messages

Follow conventional commits format:

```
[FEATURE|FIX|DOCS|REFACTOR|PERF|TEST|CHORE] Description

Detailed explanation (wrap at 72 chars)

Fixes #123
```

Examples:

```
[FEATURE] Add user profile endpoint

Implement GET /api/v1/users/profile endpoint to return
current user information with role details.

Fixes #45
```

```
[FIX] Handle nil pointer in user service

Add nil check before accessing user properties to prevent
panic in edge case when user is not found.

Fixes #67
```

### Commit Best Practices

- Make atomic commits (one logical change per commit)
- Commit frequently
- Write clear messages explaining the change
- Reference related issues

```bash
# Good
git commit -m "[FEATURE] Add email validation"

# Bad
git commit -m "updates"
```

### Keep Branch Updated

```bash
# Before creating PR, sync with main
git fetch upstream
git rebase upstream/main

# Force push if needed (on your branch only!)
git push origin feature/your-feature --force
```

## Pull Request Process

### Before Submitting

1. **Run quality checks:**
   ```bash
   make fmt      # Format code
   make vet      # Check errors
   make lint     # Run linter
   make test     # Run tests (if applicable)
   ```

2. **Update documentation:**
   - Swagger annotations for API changes
   - README.md if user-facing
   - ARCHITECTURE.md if design changes

3. **Test thoroughly:**
   ```bash
   # Manual testing
   GIN_MODE=debug go run ./cmd/server/main.go
   
   # Test with Swagger UI
   curl http://localhost:8080/swagger/index.html
   ```

### Creating a PR

**Title Format:**
```
[FEATURE|FIX|DOCS] Brief description
```

**Description Template:**
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Feature
- [ ] Bug fix
- [ ] Documentation
- [ ] Refactoring

## Related Issue
Fixes #123

## Testing
How to test this change:
- Steps to reproduce (for bugs)
- How to verify (for features)

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Swagger annotations added (if API change)
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Tests pass (if applicable)
```

### PR Review

- Be open to feedback
- Respond to all comments
- Make requested changes
- Push updates to same branch
- Request re-review

## Testing Guidelines

### Test Structure

```go
// Table-driven tests
func TestGetUser(t *testing.T) {
    tests := []struct {
        name    string
        userID  string
        want    *User
        wantErr bool
    }{
        {
            name:   "valid user",
            userID: "user123",
            want:   &User{ID: "user123"},
        },
        {
            name:    "empty user id",
            userID:  "",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := service.GetUser(context.Background(), tt.userID)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error: %v, wantErr: %v", err, tt.wantErr)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
make test

# Run specific test file
go test ./internal/domain/user/service -v

# Run with coverage
make test-coverage
```

## Documentation

### API Documentation (Swagger)

Always add Swagger annotations to new endpoints:

```go
// @Summary User login
// @Description Authenticates user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginReq body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
```

### Code Documentation

Document exported functions:

```go
// UserService handles user-related business operations.
type UserService struct {
    db *gorm.DB
}

// GetByID retrieves a user by their ID.
// Returns NotFoundError if user doesn't exist.
func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) {
```

### README Updates

Update README.md if your change:
- Adds configuration options
- Changes installation steps
- Adds new API endpoints
- Requires new dependencies

## Questions?

- Check existing issues and PRs
- Review ARCHITECTURE.md
- Ask in PR comments
- Create a discussion issue

---

Thank you for contributing! 🎉
