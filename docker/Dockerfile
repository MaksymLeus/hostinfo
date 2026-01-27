# Stage 1: Build the Go binary
FROM golang:1.24-bookworm AS go-builder

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    git ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*


WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod  ./
# COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
# COPY internal/ ./internal/

# Build arguments for multi-platform support
ARG TARGETOS=linux
ARG TARGETARCH

# Build the binary for target platform
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-s -w -X main.Version=docker" \
    -o /app/hostinfo \
    ./cmd/server/hostinfo.go

# Final stage: Debian-based runtime for proper bash support
FROM debian:bookworm-slim

# Create non-root user for security
# Add user to tty group for PTY access (required for interactive terminal)
RUN groupadd -g 1000 hostinfo && \
    useradd -u 1000 -g hostinfo -G tty -s /bin/bash -m hostinfo

WORKDIR /app

# Copy binary from builder
COPY --from=go-builder /app/hostinfo /app/hostinfo
COPY web ./web/

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