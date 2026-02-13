# ===========================
# Stage 1: Build the Frontend
# ===========================
FROM node:24-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy frontend package files
COPY frontend/package*.json ./

# Install dependencies
RUN npm ci --legacy-peer-deps

# Copy frontend source
COPY frontend/ ./

# Build frontend
RUN npm run build

# ===========================
# Stage 2: Build the Go binary
# ===========================
FROM golang:1.24-alpine AS go-builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    curl

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY assets/ ./assets/


# Copy built frontend from previous stage
COPY --from=frontend-builder /app/frontend/dist ./assets/frontend/

# Build arguments for multi-platform support
ARG TARGETOS=linux
ARG TARGETARCH

# Build static binary for Alpine (musl)
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w -X main.Version=docker" \
    -o /app/hostinfo \
    ./cmd/server


# ================================
# Stage 3: Final Alpine runtime
# ================================
FROM alpine:3.20

# Install minimal runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl

# Create non-root user
RUN addgroup -g 1000 -S hostinfo && \
    adduser -S -u 1000 -G hostinfo -s /bin/sh hostinfo

WORKDIR /app

# Copy binary from builder
COPY --from=go-builder /app/hostinfo /app/hostinfo

# Set ownership
RUN chown -R hostinfo:hostinfo /app

# Switch to non-root user
USER hostinfo

# Expose default port
EXPOSE 8080

# Health check using curl
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -sf http://localhost:8080/healthz || exit 1

# Run the application
ENTRYPOINT ["/app/hostinfo"]
