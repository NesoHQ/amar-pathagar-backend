# Multi-stage Dockerfile for Amar Pathagar Backend

# --------------------------------------------------
# Base Stage - Common dependencies
# --------------------------------------------------
FROM golang:1.25-alpine AS base
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache \
    gcc \
    musl-dev \
    bash \
    curl \
    git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# --------------------------------------------------
# Development Stage - With hot reload
# --------------------------------------------------
FROM base AS dev

# Install air for hot reload and goose for migrations
RUN go install github.com/air-verse/air@latest && \
    go install github.com/pressly/goose/v3/cmd/goose@latest

# Expose port
EXPOSE 8080

# Run with air
CMD ["air", "-c", ".air.toml"]

# --------------------------------------------------
# Builder Stage - Build the binary
# --------------------------------------------------
FROM base AS builder

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/server \
    ./cmd

# --------------------------------------------------
# Production Stage - Minimal runtime
# --------------------------------------------------
FROM alpine:latest AS prod

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    wget

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=appuser:appuser /app/server .

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./server"]
