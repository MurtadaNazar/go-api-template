# Quick Start Guide

Get the Go Platform Template running in under 5 minutes.

## 30-Second Setup

```bash
# Clone and setup
git clone <repository>
cd go-platform-template
cp .env.example .env

# Start everything with Docker
make dev-d

# Run the server
make run
```

## Access Points

| Service | URL | Notes |
|---------|-----|-------|
| **Swagger UI** | http://localhost:8080/swagger/index.html | Auto-generated API docs |
| **Health Check** | http://localhost:8080/health | Server status |
| **API Base** | http://localhost:8080/api/v1 | REST endpoints |

## First API Call

```bash
# Login
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email_or_username": "admin",
    "password": "admin"
  }'

# Response includes access_token - copy it

# Get current user
curl -X GET http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer <your_access_token>"
```

## Development Workflow

### 1. Making API Changes

```go
// Edit: internal/domain/{domain}/api/handler.go

// @Summary Your endpoint description
// @Tags Category
// @Router /path [get]
func (h *Handler) YourFunction(c *gin.Context) {
    // Implementation
}
```

### 2. Auto-regenerate Swagger

Start server in debug mode - docs auto-update within 1 second:

```bash
GIN_MODE=debug go run ./cmd/server/main.go
```

### 3. View Documentation

Open http://localhost:8080/swagger/index.html to see changes

## Essential Make Commands

```bash
# Development
make run            # Run server locally
make dev-d          # Start Docker environment
make dev-down       # Stop Docker environment

# Code Quality
make fmt            # Format code
make lint           # Check for issues
make vet            # Run go vet

# Building
make build          # Build binary
make clean          # Clean artifacts

# Get Help
make help           # All available commands
```

## Common Tasks

### Start Fresh Development Session

```bash
# Terminal 1: Start services
make dev-d

# Terminal 2: Run application
GIN_MODE=debug make run

# Terminal 3: View logs (optional)
make dev-logs
```

### Add New API Endpoint

1. Create handler in `internal/domain/{domain}/api/handler.go`
2. Add Swagger annotations (`// @`)
3. Save file
4. Docs auto-regenerate
5. Test in Swagger UI

### Check Code Quality

```bash
make fmt    # Format
make vet    # Check errors
make lint   # Full linter
```

### Build for Production

```bash
make build
./go-platform-template
```

## Default Credentials

```
Username: admin
Password: admin
```

⚠️ **Change these in production!**

## Stop Development Environment

```bash
make dev-down
```

## Environment Variables

Copy `.env.example` to `.env` and modify:

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
JWT_SECRET_KEY=change-this-in-production
JWT_EXPIRE_HOURS=24
```

## Useful Links

- **Full Guide**: See `README.md`
- **Architecture**: See `ARCHITECTURE.md`
- **Contributing**: See `CONTRIBUTING.md`
- **Status**: See `PROJECT_COMPLETION.md`

## Troubleshooting

### Port Already in Use

```bash
lsof -ti:8080 | xargs kill -9
```

### Database Connection Failed

- Ensure PostgreSQL is running
- Check `DATABASE_URL` in `.env`
- Try: `psql $DATABASE_URL`

### Swagger Not Generating

- Ensure `GIN_MODE=debug`
- Check `swag` is installed: `go install github.com/swaggo/swag/cmd/swag@latest`

### Docker Issues

```bash
# Stop everything
make dev-down

# Clean up
docker system prune

# Restart
make dev-d
```

## API Examples

### Register User

```bash
curl -X POST http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "securepassword"
  }'
```

### Upload File

```bash
curl -X POST http://localhost:8080/api/v1/files/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@path/to/file.pdf"
```

### List Users

```bash
curl -X GET "http://localhost:8080/api/v1/users/?offset=0&limit=10" \
  -H "Authorization: Bearer <token>"
```

## Next Steps

1. Read `README.md` for complete setup
2. Check `ARCHITECTURE.md` to understand design
3. Review `CONTRIBUTING.md` for development guidelines
4. Explore API docs at http://localhost:8080/swagger/index.html

## Need Help?

- Check the troubleshooting section in `README.md`
- Review error logs: `make dev-logs`
- Check your `.env` configuration
- Run `make help` to see all options

---

**Status:** ✅ Production Ready  
**Documentation:** Complete  
**Build Status:** Passing  
**Linting:** Zero Issues
