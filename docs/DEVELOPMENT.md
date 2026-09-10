<!-- omit in toc -->
# Development Guide

This document covers the development setup and workflow for **hostinfo**.

<!-- omit in toc -->
## Table of Contents
- [📦 Prerequisites](#-prerequisites)
- [Quick Setup](#quick-setup)
- [Development Workflow](#development-workflow)
  - [Recommended: Two Terminal Setup](#recommended-two-terminal-setup)
- [Project Structure](#project-structure)
- [Backend Development](#backend-development)
  - [Running](#running)
  - [Testing](#testing)
  - [Linting](#linting)
- [Frontend Development](#frontend-development)
  - [Running](#running-1)
  - [Testing](#testing-1)
  - [Linting](#linting-1)
- [Building](#building)
  - [Development Build](#development-build)
  - [Build Script Options](#build-script-options)
  - [Docker Build](#docker-build)
- [Contributing](#contributing)
  - [Getting Started](#getting-started)
  - [Code Style](#code-style)
  - [Commit messages follow](#commit-messages-follow)
  - [Pull Request Guidelines](#pull-request-guidelines)

## 📦 Prerequisites
- **Go**: Version 1.21 or higher ([Download](https://golang.org/dl/))
- **Node.js**: Version 24 or higher ([Download](https://nodejs.org/))
- **npm**: Version 11 or higher (comes with Node.js)

## Quick Setup

```bash
# Clone the repository
git clone https://github.com/MaksymLeus/hostinfo.git
cd hostinfo

# Run prepare script 
# The prepare.sh script bootstraps the development environment for the Go project.
# For more description check ./tools/TOOLS.md
./tools/prepare.sh

# Build and run
./build.sh
./bin/hostinfo
```
Access at: `http://localhost:8080`

## Development Workflow

### Recommended: Two Terminal Setup

**Terminal 1 - Backend:**

```bash
go run cmd/server/hostinfo.go
```

**Terminal 2 - Frontend (with hot reload):**

```bash
cd frontend && npm run dev
```

Access frontend dev server at: `http://localhost:5173`
(API calls are proxied to the Go backend on port `8080`)

## Project Structure
```
hostinfo/
├── CHANGELOG.md
├── Dockerfile
├── LICENSE.md
├── README.md
├── cmd
│   └── server
│       └── hostinfo.go
├── docker-compose.yml
├── docs
├── frontend
│   └── assets
├── go.mod
├── go.sum
├── helm
│   ├── Chart.yaml
│   ├── templates
│   │   ├── _helpers.tpl
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   └── values.yaml
├── internal
│   ├── api
│   │   ├── middleware
│   │   │   ├── cors.go
│   │   │   ├── logging.go
│   │   │   └── ratelimit.go
│   │   └── v1
│   │       ├── handlers
│   │       │   ├── cloud.go
│   │       │   ├── commands.go
│   │       │   ├── info.go
│   │       │   ├── network.go
│   │       │   ├── runtime.go
│   │       │   └── type.go
│   │       └── routes.go
│   ├── custom
│   │   └── config.go
│   ├── health
│   │   ├── checks.go
│   │   ├── handlers.go
│   │   └── types.go
│   └── server
│       └── server.go
└── tools
    ├── TOOLS.md
    ├── build.sh
    ├── prepare.sh
    └── scripts
        └── git_hooks
            ├── GIT_HOOKS.md
            ├── commit-msg
            └── pre-commit
```

## Backend Development

### Running

```bash
# Run with go run
go run cmd/hostinfo/main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/validation/...

# Verbose output
go test -v ./...
```

### Linting

```bash
# Format code
go fmt ./...

# Run vet
go vet ./...

# Run golangci-lint (if installed)
golangci-lint run
```

## Frontend Development

### Running

```bash
cd frontend

# Start dev server with hot reload
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Testing

```bash
cd frontend

# Run tests
npm test

# Run tests with coverage
npm test -- --coverage
```

### Linting

```bash
cd frontend

# Run ESLint
npm run lint

# Fix auto-fixable issues
npm run lint -- --fix
```

## Building

### Development Build

```bash
# Quick build for current platform
go build -o hostinfo cmd/hostinfo/main.go
```

### Build Script Options

```bash
  ./build.sh                      # Quick local build
  ./build.sh all                  # Build everything into ./bin/
  ./build.sh os linux             # Build all linux binaries
  ./build.sh arch arm64           # Build all arm64 binaries
  ./build.sh target darwin amd64  # Build macOS Intel binary
  ./build.sh checksums            # Generate SHA256 checksums for bin/ artifacts
  ./build.sh clean                # Remove build artifacts
  ./build.sh list                 # List all supported OS/ARCH combinations
```
Build for all platforms
```bash
# Build for all platforms
./build.sh all

# Outputs:
# bin/hostinfo-linux-x64
# bin/hostinfo-darwin-x64
# bin/hostinfo-darwin-arm64
```

### Docker Build

```bash
# Build Docker image
docker build -t hostinfo .

# Or using docker compose
docker compose build
```

## Contributing

### Getting Started

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Write/update tests
5. Commit: `git commit -m 'Add some amazing feature'`
6. Push: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Code Style

- **Go**: Follow standard Go formatting (`go fmt`)
- **JavaScript/React**: Follow ESLint rules
- Write tests for new features
- Update documentation as needed

### Commit messages follow 

**Conventional Commits**:
- `breaking:` → breaking change
- `feat:` → new feature
- `fix:` → bug fix
- `chore:` → maintenance
- `docs:` → documentation update
- `refactor:` → code refactor
- `test:` → tests only
- `perf:` → performance improvement

This ensures **automatic versioning** during release.

| Change | Commit | Version Bump |
|---|---|---|
| Add new cloud provider support | `feat: add oracle cloud detection` | `1.2.0 → 1.3.0` |
| Fix EC2 region extraction | `fix: correct EC2 region detection` | `1.3.0 → 1.3.1` |
| Remove field or change API | `breaking: remove instance_type field or footer` | `1.3.1 → 2.0.0` |

### Pull Request Guidelines

- Ensure all tests pass
- Update documentation if needed
- Keep PRs focused (one feature/fix per PR)
- Provide clear description of changes

---

For configuration options, see [CONFIGURATION.md](CONFIGURATION.md).

For deployment instructions, see [DEPLOYMENT.md](DEPLOYMENT.md).