# Architecture Documentation

## Overview

This document describes the architectural design and patterns used in the Go Platform Template.

## Design Patterns

### 1. Layered Architecture

The application follows a clean layered architecture pattern:

```
┌─────────────────────────────────┐
│    HTTP Layer (Gin Router)      │
├─────────────────────────────────┤
│   API Layer (Handlers)          │
│   - Request validation          │
│   - Response formatting         │
│   - Error handling              │
├─────────────────────────────────┤
│   Service Layer (Business)      │
│   - Core business logic         │
│   - Data validation             │
│   - Cross-domain operations     │
├─────────────────────────────────┤
│   Model/Repository Layer        │
│   - Data access                 │
│   - Query building              │
├─────────────────────────────────┤
│    Database (PostgreSQL)        │
└─────────────────────────────────┘
```

**Benefits:**
- Clear separation of concerns
- Easy to test each layer independently
- Database changes don't affect handlers
- Business logic is reusable

### 2. Domain-Driven Design (DDD)

The codebase is organized around business domains:

```
internal/domain/
├── auth/          # Authentication domain
├── user/          # User management domain
└── file/          # File management domain
```

Each domain is self-contained with its own:
- **API Layer** (`api/`) - HTTP handlers
- **Service Layer** (`service/`) - Business logic
- **Model Layer** (`model/`) - Data structures
- **Repository** - Data access (implicit via service)

**Advantages:**
- Easy to locate related code
- Domains can be developed independently
- Clear boundaries and dependencies

### 3. Handler/Service/Model Pattern

```
Handler (HTTP)
    ↓
Service (Business Logic)
    ↓
Model (Data)
    ↓
Database
```

**Example Flow:**

```go
// Handler: Accept HTTP request
func (h *UserHandler) GetUser(c *gin.Context) {
    userID := c.Param("id")
    
    // Call service
    user, err := h.service.GetByID(c.Request.Context(), userID)
    if err != nil {
        c.Error(err)
        return
    }
    
    // Return response
    c.JSON(http.StatusOK, user)
}

// Service: Execute business logic
func (s *UserService) GetByID(ctx context.Context, id string) (*User, error) {
    // Validation
    if id == "" {
        return nil, errors.New("invalid id")
    }
    
    // Query database
    var user User
    if err := s.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
        return nil, err
    }
    
    return &user, nil
}

// Model: Data structure
type User struct {
    ID        string    `gorm:"primaryKey"`
    Email     string    `gorm:"uniqueIndex"`
    Password  string    `gorm:"-"` // Sensitive field
    CreatedAt time.Time
}
```

## Error Handling

### Error Categories

```
apperrors/
├── BadRequestError      (400) - Invalid input
├── UnauthorizedError    (401) - Authentication failed
├── ForbiddenError       (403) - Insufficient permissions
├── NotFoundError        (404) - Resource not found
├── ConflictError        (409) - Resource conflict
└── InternalError        (500) - Server error
```

### Error Flow

```
Database Error
    ↓
Service converts to AppError
    ↓
Handler receives AppError
    ↓
Middleware formats response
    ↓
Client receives structured error
```

### Example

```go
// Service detects error
if user == nil {
    return nil, apperrors.NewAppError(
        apperrors.NotFoundError, 
        "User not found",
    )
}

// Handler receives error
if appErr, ok := apperrors.IsAppError(err); ok {
    _ = c.Error(appErr)
    return
}

// Middleware handles it (via error handler)
```

## API Design

### RESTful Conventions

```
GET     /api/v1/users/          - List users
POST    /api/v1/users/          - Create user
GET     /api/v1/users/:id       - Get user
PUT     /api/v1/users/:id       - Update user
DELETE  /api/v1/users/:id       - Delete user
```

### Request/Response Format

All API responses follow a standard format:

```json
{
  "data": {...},
  "error": null,
  "request_id": "12345..."
}
```

### Status Codes

- **200** - OK (successful GET, PUT)
- **201** - Created (successful POST)
- **204** - No Content (successful DELETE)
- **400** - Bad Request (validation error)
- **401** - Unauthorized (auth failed)
- **403** - Forbidden (insufficient permissions)
- **404** - Not Found (resource missing)
- **409** - Conflict (duplicate resource)
- **500** - Internal Server Error

## Authentication & Authorization

### JWT Flow

```
Login Request
    ↓
Validate Credentials
    ↓
Generate JWT Tokens (access + refresh)
    ↓
Return Tokens to Client
    ↓
Client uses access token in Authorization header
    ↓
Middleware verifies token signature
    ↓
Extract claims (user ID, role)
    ↓
Proceed to handler
```

### Token Structure

```
Access Token:
- Type: Short-lived (configurable, default: 24h)
- Contains: user_id, role
- Used for: Authenticating requests

Refresh Token:
- Type: Long-lived
- Used for: Getting new access tokens
- Stored in: Database (for revocation)
```

### Authorization Check

```go
// Middleware extracts claims
if claims := extractClaims(token); claims.Role != "admin" {
    return apperrors.ForbiddenError
}
```

## Database Design

### ORM Usage

The project uses **GORM** for database operations:

```go
// Query
var user User
db.Where("email = ?", email).First(&user)

// Create
db.Create(&user)

// Update
db.Model(&user).Updates(updates)

// Delete
db.Delete(&user)
```

### Schema Management

Migrations are auto-applied on startup:

```go
// In seeder.go
db.AutoMigrate(&User{}, &File{})
```

### Indexes

Key fields are indexed for performance:

```go
type User struct {
    Email string `gorm:"uniqueIndex"`
    // Other fields
}
```

## File Storage

### MinIO Integration

Files are stored in MinIO (S3-compatible):

```
File Upload Flow:
    ↓
Validate file
    ↓
Generate unique filename
    ↓
Upload to MinIO
    ↓
Store metadata in database
    ↓
Return download URL

File Download Flow:
    ↓
Verify user owns file
    ↓
Get presigned URL from MinIO
    ↓
Redirect to download
```

### Isolation

Each user's files are isolated:

```
/bucket/
├── user-1/
│   ├── file1.pdf
│   └── file2.doc
├── user-2/
│   └── file3.jpg
```

## Swagger/OpenAPI Integration

### Auto-Generation

Swagger docs are generated from code annotations:

```go
// @Summary Get user
// @Description Returns user by ID
// @Tags Users
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    // Implementation
}
```

### File Watching

In debug mode, a file watcher monitors:
- `internal/domain/*/api/` - Handler files
- `internal/domain/*/model/` - Model files  
- `internal/domain/*/dto/` - DTO files

When changes are detected, `swag init` is automatically triggered.

## Middleware Stack

Request flows through middleware in order:

```
Request
    ↓
CORS Middleware
    ↓
Request ID Middleware (adds tracing)
    ↓
Logging Middleware
    ↓
Error Middleware
    ↓
Auth Middleware (if protected route)
    ↓
Handler
    ↓
Response
```

## Configuration Management

Configuration is loaded from environment variables with `godotenv`:

```go
type Config struct {
    ServerAddr  string
    DatabaseURL string
    GinMode     string
    // ... other fields
}

func LoadConfig() *Config {
    // Loads from .env or environment
}
```

## Logging Strategy

Uses **Zap** for structured logging:

```go
// Info level
logger.Info("User created successfully")

// With fields
logger.Infow("User logged in", 
    "user_id", userID,
    "ip", clientIP,
)

// Error with stack trace
logger.Errorw("Database error", 
    "error", err,
    "request_id", requestID,
)
```

## Testing Strategy

### Recommended Test Coverage

```
Models      → Data validation (unit tests)
Services    → Business logic (unit + integration tests)
Handlers    → HTTP layer (integration tests)
Middleware  → Request processing (integration tests)
```

### Testing Patterns

```go
// Table-driven tests
var tests = []struct {
    name    string
    input   interface{}
    want    interface{}
    wantErr bool
}{
    {"case 1", input1, want1, false},
    {"case 2", input2, want2, true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test logic
    })
}
```

## Security Considerations

### Password Security
- Passwords hashed with bcrypt
- Never stored in plain text
- Never logged or returned in responses

### JWT Security
- Tokens signed with secret key
- Expiration time enforced
- Refresh token rotation supported

### SQL Injection Prevention
- GORM parameterized queries
- User input always validated
- Never build SQL strings manually

### CORS
- Configured per environment
- Whitelist allowed origins
- Control HTTP methods

## Performance Optimizations

### Database
- Connection pooling (DB_MAX_OPEN_CONNS)
- Indexes on frequently queried fields
- Pagination for large result sets

### Caching
- Consider Redis for session tokens
- Cache frequently accessed data

### API Design
- Pagination (offset/limit)
- Selective field inclusion
- Compression support

## Deployment Architecture

### Development
```
Docker Compose
├── Go App
├── PostgreSQL
└── MinIO
```

### Production
```
Load Balancer
    ↓
├── Go App Instance 1
├── Go App Instance 2
└── Go App Instance N
    ↓
PostgreSQL (managed)
    ↓
MinIO / S3
    ↓
Redis (optional)
```

## Future Considerations

- [ ] Caching layer (Redis)
- [ ] Message queue (RabbitMQ/Kafka)
- [ ] GraphQL API layer
- [ ] Rate limiting
- [ ] WebSocket support
- [ ] Multi-tenancy support
- [ ] Audit logging
- [ ] API versioning strategy

---

For code examples and implementation details, refer to the specific domain packages in `internal/domain/`.
