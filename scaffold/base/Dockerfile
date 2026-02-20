# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.Version=docker -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" \
    -o app cmd/server/main.go

# Final stage
FROM alpine:3.18

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/app .
COPY --from=builder /build/.env.example .

EXPOSE 8080

CMD ["./app"]
