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
COPY web/ ./web/

# Build arguments for multi-platform support
ARG TARGETOS=linux
ARG TARGETARCH

# Build the binary for target platform
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-s -w -X main.Version=docker" \
    -o /app/hostinfo.go \
    ./cmd/server/hostinfo.go

# Final stage: Debian-based runtime for proper bash support
FROM debian:bookworm-slim

# Install required packages:
# - ca-certificates: for HTTPS connections
# - tzdata: for timezone support
# - bash: full bash shell for script execution
# - coreutils: standard Unix utilities
# - curl: for health checks and HTTP operations
# - openssh-client: for SSH connections to remote servers
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    bash \
    zsh \
    coreutils \
    curl \
    openssh-client \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# Create non-root user for security
# Add user to tty group for PTY access (required for interactive terminal)
RUN groupadd -g 1000 hostinfo && \
    useradd -u 1000 -g hostinfo -G tty -s /bin/bash -m hostinfo

# Create data and config directories
RUN mkdir -p /data /config && \
    chown -R hostinfo:hostinfo /data /config

# Create .ssh directory for SSH key operations
RUN mkdir -p /home/hostinfo/.ssh && \
    chown hostinfo:hostinfo /home/hostinfo/.ssh && \
    chmod 700 /home/hostinfo/.ssh

# Create tmp directory for session files (SSH configs, keys, wrappers)
# This ensures the non-root user can create temp files for terminal sessions
RUN mkdir -p /home/hostinfo/tmp && \
    chown hostinfo:hostinfo /home/hostinfo/tmp && \
    chmod 700 /home/hostinfo/tmp

WORKDIR /app

# Copy binary from builder
COPY --from=go-builder /app/hostinfo.go /app/hostinfo.go

# Set ownership
RUN chown hostinfo:hostinfo /app/hostinfo.go

# Switch to non-root user
USER hostinfo

# Expose default port
EXPOSE 8080

# Health check using curl
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -sf http://localhost:8080/healthz || exit 1

# Environment variables with defaults
# TMPDIR is set to user's home for proper temp file creation in non-root context
# hadolint ignore=DL3044
ENV TMPDIR=/home/hostinfo/tmp \
    SHELL=/bin/bash

# Run the application
ENTRYPOINT ["/app/hostinfo.go"]