# {{.ProjectName}}

A production-ready Go API built with the Go Platform Template scaffolder.

## Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose (optional, for containerized development)
- PostgreSQL 12+ (if not using Docker)
- MinIO (if using file storage without Docker)

### Setup

```bash
# Install dependencies
go mod download

# Copy environment variables
cp .env.example .env

# Edit .env with your settings
# vim .env

# Start development environment
make dev-d

# Check API is running
curl http://localhost:8080/health

# View API documentation
# http://localhost:8080/swagger/index.html
```

## Available Commands

```bash
# Build
make build              # Build the application

# Testing
make test               # Run all tests
make test-coverage      # Run tests with coverage

# Development
make dev                # Start dev environment (foreground)
make dev-d              # Start dev environment (background)
make dev-down           # Stop dev environment
make dev-logs           # View logs

# Code Quality
make fmt                # Format code
make lint               # Run linter
make vet                # Run go vet
make security           # Security checks

# Dependencies
make deps               # Download dependencies
make verify             # Verify dependencies
make update-deps        # Update dependencies
```

## Project Structure

```
{{.ProjectName}}/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── app/                 # Application bootstrap
│   │   ├── bootstrap.go
│   │   ├── routes.go        # Route definitions
│   │   └── middleware.go
│   ├── domain/              # Business logic
│   │   ├── auth/            # Authentication (if selected)
│   │   ├── user/            # User management (if selected)
│   │   └── file/            # File handling (if selected)
│   ├── platform/            # Infrastructure
│   │   ├── config/          # Configuration
│   │   ├── logger/          # Logging
│   │   ├── database/        # Database
│   │   └── http/            # HTTP utilities
│   └── shared/              # Cross-cutting concerns
├── Makefile                 # Build automation
├── Dockerfile               # Container image
├── docker-compose.yml       # Service orchestration
├── go.mod & go.sum          # Dependencies
└── .env.example             # Configuration template
```

## Configuration

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

### Environment Variables

- **SERVER_PORT** - API port (default: 8080)
- **SERVER_ENV** - Environment (development/production)
- **DB_HOST** - Database host
- **DB_PORT** - Database port
- **DB_USER** - Database user
- **DB_PASSWORD** - Database password
- **DB_NAME** - Database name
- **JWT_SECRET** - JWT signing secret
- **JWT_EXPIRY** - Access token expiry
- **JWT_REFRESH_EXPIRY** - Refresh token expiry
- **MINIO_ENDPOINT** - MinIO endpoint (if using file storage)
- **MINIO_ACCESS_KEY** - MinIO access key
- **MINIO_SECRET_KEY** - MinIO secret key
- **MINIO_BUCKET** - MinIO bucket name

## Development Workflow

### Start Development Environment

```bash
make dev-d
```

This starts:
- PostgreSQL database (port 5432)
- MinIO file storage (ports 9000/9001)
- Your API server (port 8080)

### Stop Development Environment

```bash
make dev-down
```

### View Logs

```bash
make dev-logs
```

### Run Tests

```bash
make test

# With coverage
make test-coverage
```

### Code Quality

```bash
make fmt                # Format code
make lint               # Run linter
make vet                # Run go vet
make security           # Security checks
```

## API Documentation

When running the development environment, visit:

```
http://localhost:8080/swagger/index.html
```

Swagger documentation is auto-generated from code comments.

## Building for Production

### Build Binary

```bash
make build
```

### Build Docker Image

```bash
docker build -t {{.ProjectName}}:latest .
```

### Run Container

```bash
docker run -p 8080:8080 {{.ProjectName}}:latest
```

## Database Migrations

Migrations are handled automatically on startup.

To manually run migrations:

```go
// In your bootstrap code
if err := db.AutoMigrate(&models.User{}, &models.File{}); err != nil {
    log.Fatal(err)
}
```

## Features

### Included
- ✅ Clean Architecture (DDD)
- ✅ Structured Logging (Zap)
- ✅ Error Handling
- ✅ Request/Response wrappers

### Optional (selected during scaffolding)
- ✅ JWT Authentication
- ✅ User Management & RBAC
- ✅ PostgreSQL Database
- ✅ MinIO File Storage
- ✅ Swagger API Docs
- ✅ Docker & Docker Compose

## Security

- JWT RS256 token signing
- Secure password hashing (bcrypt)
- CORS support
- SQL injection prevention (GORM)
- Request validation

## Troubleshooting

### Port Already in Use

Change `SERVER_PORT` in `.env`:

```env
SERVER_PORT=8081
```

### Database Connection Failed

Ensure PostgreSQL is running and credentials are correct in `.env`.

### MinIO Connection Failed

Ensure MinIO is running on the configured endpoint.

### Tests Failing

```bash
# Clean and rebuild
make clean
make test
```

## Contributing

1. Create a feature branch
2. Make your changes
3. Run tests: `make test`
4. Run linter: `make lint`
5. Submit a pull request

## License

MIT License - See LICENSE file for details.

## Support

For issues or questions:
1. Check the error message carefully
2. Review the configuration in `.env`
3. Check application logs: `make dev-logs`
4. Refer to the architecture documentation

---

Built with ❤️ using Go 1.25+
