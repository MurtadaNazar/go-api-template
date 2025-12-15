# Go Platform Template

A production-ready Go REST API platform with authentication, user management, file storage, and RBAC (Role-Based Access Control). Built with modern Go best practices, featuring PostgreSQL integration, JWT authentication, MinIO file storage, and auto-generated Swagger documentation.

## Table of Contents

- [Features](#features)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Development Workflow](#development-workflow)
- [Making Changes and Auto-Regenerating Swagger](#making-changes-and-auto-regenerating-swagger)
- [Database Migrations](#database-migrations)
- [Project Architecture](#project-architecture)
- [Code Quality](#code-quality)
- [Docker Deployment](#docker-deployment)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Authentication & Authorization**
  - JWT-based token authentication
  - Refresh token mechanism with token rotation
  - Role-based access control (Admin, User)
  - Secure password hashing

- **User Management**
  - User registration and login
  - User profile management
  - Pagination and filtering support
  - User type differentiation (admin/regular user)

- **File Management**
  - File upload and download
  - MinIO S3-compatible storage
  - Per-user file isolation
  - Secure file operations

- **API Features**
  - RESTful API design
  - Automatic Swagger/OpenAPI documentation
  - Request/Response validation
  - Error handling with custom error codes
  - Request ID tracking for debugging
  - CORS support

- **Development Features**
  - Auto-regenerating Swagger docs in debug mode
  - Database seeding
  - Hot reload support (with air)
  - Structured logging with Zap
  - Docker Compose for development environment

## Technology Stack

| Component | Technology |
|-----------|-----------|
| **Runtime** | Go 1.25.1 |
| **Framework** | Gin Web Framework |
| **Database** | PostgreSQL |
| **Authentication** | JWT (RS256) |
| **File Storage** | MinIO (S3-compatible) |
| **Logging** | Zap (structured logging) |
| **API Documentation** | Swagger/OpenAPI 3.0 |
| **Containerization** | Docker & Docker Compose |

## Project Structure

```
go-platform-template/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── app/                        # Application bootstrap
│   │   ├── db.go                   # Database initialization
│   │   ├── swagger.go              # Swagger setup & file watching
│   │   ├── middleware.go           # HTTP middleware
│   │   ├── routes.go               # Route registration
│   │   ├── health.go               # Health check handlers
│   │   └── server.go               # Server lifecycle management
│   ├── domain/                     # Business logic by domain
│   │   ├── auth/
│   │   │   ├── api/                # HTTP handlers
│   │   │   ├── model/              # Data models
│   │   │   └── service/            # Business logic
│   │   ├── user/
│   │   │   ├── api/
│   │   │   ├── model/
│   │   │   └── service/
│   │   └── file/
│   │       ├── api/
│   │       ├── model/
│   │       └── service/
│   └── platform/                   # Cross-cutting concerns
│       ├── config/                 # Configuration management
│       ├── logger/                 # Logging setup
│       ├── database/               # Database connections
│       └── errors/                 # Error definitions
├── docs/                           # Generated Swagger documentation
├── docker-compose/                 # Docker Compose configuration
├── scripts/                        # Utility scripts
├── Dockerfile                      # Container image definition
├── Makefile                        # Build automation
├── go.mod & go.sum                 # Dependency management
└── .env.example                    # Environment template
```

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go** 1.25.1 or higher
- **PostgreSQL** 12+ (or use Docker Compose)
- **MinIO** (or use Docker Compose)
- **Make** (for build automation)
- **Git** (for version control)

### Optional but Recommended

- **Docker** & **Docker Compose** (for containerized development)
- **golangci-lint** (for code linting)
- **swag** (auto-installed, required for Swagger generation)

## Installation

### 1. Clone the Repository

```bash
git clone <repository-url>
cd go-platform-template
```

### 2. Install Dependencies

```bash
make deps
```

Or manually:

```bash
go mod download
go mod verify
```

### 3. Set Up Environment Variables

```bash
cp .env.example .env
```

Then edit `.env` with your configuration:

```env
# Server
SERVER_ADDR=:8080
GIN_MODE=debug

# Database (required)
DATABASE_URL=postgres://postgres:postgres@localhost:5432/go_platform_template?sslmode=disable

# MinIO (file storage)
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin

# JWT
JWT_SECRET_KEY=your-secret-key-change-in-production
JWT_EXPIRE_HOURS=24

# Logging
LOG_LEVEL=info
```

### 4. Start Development Environment

**Option A: Using Docker Compose (Recommended)**

```bash
make dev-d
```

This starts PostgreSQL, MinIO, and other services in the background.

**Option B: Manual Setup**

Ensure PostgreSQL and MinIO are running on your system with the configured URLs.

## Configuration

The application uses environment variables for configuration. See `.env.example` for all available options.

### Key Configuration Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GIN_MODE` | Gin framework mode (debug, release, test) | release |
| `SERVER_ADDR` | Server address and port | :8080 |
| `DATABASE_URL` | PostgreSQL connection string | - |
| `JWT_SECRET_KEY` | Secret key for JWT signing | - |
| `JWT_EXPIRE_HOURS` | JWT expiration time in hours | 24 |
| `LOG_LEVEL` | Log level (debug, info, warn, error) | info |
| `MINIO_ENDPOINT` | MinIO server endpoint | localhost:9000 |

## Running the Application

### Development Mode

```bash
# Start with auto-reload (requires 'air' package)
air

# Or run directly
make run

# Or with GIN_MODE=debug for auto-regenerating Swagger
GIN_MODE=debug go run ./cmd/server/main.go
```

### Production Mode

```bash
# Build binary
make build

# Run binary
./go-platform-template
```

The server will start on the configured `SERVER_ADDR` (default: `http://localhost:8080`).

### Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "database": "connected",
  "storage": "connected"
}
```

## API Documentation

### Swagger/OpenAPI Documentation

The API is fully documented with Swagger/OpenAPI 3.0. The documentation is automatically generated and updated.

**Access the interactive API docs:**

1. **Start the server in debug mode:**
   ```bash
   GIN_MODE=debug go run ./cmd/server/main.go
   ```

2. **Open your browser:**
   ```
   http://localhost:8080/swagger/index.html
   ```

### Available Endpoints

#### Authentication
- `POST /api/v1/login` - User login
- `POST /api/v1/refresh` - Refresh access token
- `POST /api/v1/logout` - User logout
- `GET /api/v1/me` - Get current user info

#### Users
- `GET /api/v1/users/` - List all users (admin only)
- `POST /api/v1/users/` - Register new user
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

#### Files
- `POST /api/v1/files/upload` - Upload file
- `GET /api/v1/files/:filename` - Download file
- `GET /api/v1/files/` - List user's files
- `DELETE /api/v1/files/:filename` - Delete file

### Authentication

All protected endpoints require the `Authorization` header with a Bearer token:

```bash
curl -H "Authorization: Bearer <access_token>" \
     http://localhost:8080/api/v1/users/
```

## Development Workflow

### Making Changes and Auto-Regenerating Swagger

When you add or modify API endpoints with Swagger annotations, the documentation is **automatically regenerated** in debug mode.

#### Example: Adding a New Endpoint

1. **Create your handler with Swagger annotations:**

```go
// @Summary Get user profile
// @Description Returns the current user's profile
// @Tags Users
// @Security BearerAuth
// @Success 200 {object} model.User
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
    // Implementation
}
```

2. **Start server in debug mode:**
```bash
GIN_MODE=debug go run ./cmd/server/main.go
```

3. **Logs show auto-regeneration:**
```
{"level":"INFO","message":"Watching directory: internal/domain/user/api"}
{"level":"INFO","message":"Swagger docs regenerated"}
```

4. **Visit Swagger UI to see updated docs:**
```
http://localhost:8080/swagger/index.html
```

### Available Make Commands

```bash
# Build & Run
make build          # Build binary
make run            # Run locally
make clean          # Clean artifacts

# Testing
make test           # Run tests
make test-coverage  # Generate coverage report

# Code Quality
make fmt            # Format code
make vet            # Run go vet
make lint           # Run linter
make security       # Security checks

# Dependencies
make deps           # Download dependencies
make verify         # Verify dependencies
make update-deps    # Update dependencies

# Development
make dev            # Start Docker environment (foreground)
make dev-d          # Start Docker environment (background)
make dev-down       # Stop Docker environment
make dev-logs       # View Docker logs

# Help
make help           # Show all commands
```

### Code Formatting & Quality

Before committing, ensure your code passes quality checks:

```bash
make fmt    # Format code
make vet    # Check for errors
make lint   # Run linter
```

## Database Migrations

The application uses GORM for database management. Migrations are automatically applied on startup from the seeder.

### Creating New Migrations

1. Define your model in the appropriate domain package
2. Add migration logic to `internal/platform/database/seeder.go`
3. Restart the application

## Project Architecture

### Layered Architecture

```
HTTP Layer (Gin Handlers)
    ↓
API Layer (Request validation, Response formatting)
    ↓
Service Layer (Business logic)
    ↓
Model Layer (Data structures)
    ↓
Database Layer (GORM)
    ↓
Database (PostgreSQL)
```

### Domain-Driven Design

The project is organized by business domains:

- **Auth Domain**: Authentication and authorization logic
- **User Domain**: User management and profiles
- **File Domain**: File storage and retrieval

Each domain contains:
- `api/` - HTTP handlers and routes
- `model/` - Data structures and requests/responses
- `service/` - Business logic

### Error Handling

Custom error handling with specific error codes:

```go
// In handlers
if err := service.SomeOperation(); err != nil {
    if appErr, ok := apperrors.IsAppError(err); ok {
        _ = c.Error(appErr)
        return
    }
}
```

## Code Quality

### Linting

```bash
make lint
```

The project uses `golangci-lint` with strict rules:
- No unchecked errors
- Code formatting compliance
- Security checks
- Performance analysis

### Testing

```bash
make test            # Run all tests
make test-coverage   # Generate coverage report (coverage.html)
```

## Docker Deployment

### Development Environment

```bash
# Start all services
make dev-d

# View logs
make dev-logs

# Stop services
make dev-down
```

Services started:
- PostgreSQL (port 5432)
- MinIO (port 9000, console 9001)
- Redis (optional)

### Production Build

```bash
# Build Docker image
docker build -t go-platform-template:latest .

# Run container
docker run -p 8080:8080 \
  -e DATABASE_URL=postgres://user:pass@host:5432/db \
  -e JWT_SECRET_KEY=your-secret \
  go-platform-template:latest
```

## Troubleshooting

### Common Issues

#### 1. "Failed to connect to database"
- Ensure PostgreSQL is running
- Check `DATABASE_URL` in `.env`
- Verify database credentials

```bash
# Test connection
psql $DATABASE_URL
```

#### 2. "Swagger docs not generating"
- Ensure `swag` is installed:
  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```
- Check server is running in debug mode: `GIN_MODE=debug`
- Verify swagger annotations are correct

#### 3. "File upload failing"
- Ensure MinIO is running
- Check `MINIO_ENDPOINT` is accessible
- Verify MinIO credentials in `.env`

#### 4. "Port 8080 already in use"
- Change port: `SERVER_ADDR=:8081 go run ./cmd/server/main.go`
- Or kill process: `lsof -ti:8080 | xargs kill -9`

#### 5. Build errors
- Update dependencies: `make update-deps`
- Clear Go cache: `go clean -cache`
- Verify Go version: `go version` (requires 1.25.1+)

### Debug Mode

Enable detailed logging:

```bash
LOG_LEVEL=debug GIN_MODE=debug go run ./cmd/server/main.go
```

### Log Locations

- Console output (stdout)
- Log files in `./logs/` directory

## Contributing

### Getting Started

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes
4. Run quality checks: `make fmt make vet make lint`
5. Commit with clear messages
6. Push and create a Pull Request

### Coding Standards

- Follow Go conventions (effective Go)
- Use clear, descriptive names
- Add comments for exported functions
- Keep functions small and focused
- Write Swagger annotations for API endpoints

### Commit Message Format

```
[FEATURE|FIX|DOCS|REFACTOR] Brief description

Detailed explanation of changes and reasoning
```

### Code Review

All contributions require review. Ensure:
- Code passes all quality checks
- Tests pass (when applicable)
- Documentation is updated
- Swagger annotations are correct

## License

[Specify your license here]

---

## Support

For issues, questions, or suggestions:
- Create an GitHub issue
- Check existing documentation
- Review API documentation at `/swagger/index.html`

---

**Last Updated:** December 2024  
**Maintained By:** [Your Team/Name]
