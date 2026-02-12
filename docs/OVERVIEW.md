<!-- omit in toc -->
# Project Overview

HostInfo is a lightweight Go-based service that exposes system, runtime, and environment information through both an Web dashboard and  REST API.

It is designed for DevOps engineers, SREs, platform teams, and automation systems that require fast, portable, and container-ready host telemetry.

<!-- omit in toc -->
## Table of Contents
- [Project Goals](#project-goals)
- [What Information Does It Expose?](#what-information-does-it-expose)
- [🧱 Architecture Overview](#-architecture-overview)
- [🗂 Repository Structure](#-repository-structure)
- [Cloud Detection Layer](#cloud-detection-layer)
- [📦 Deployment Models](#-deployment-models)
- [🛠 Developer Experience](#-developer-experience)
- [🧑‍💻 Target Users](#-target-users)
- [🪄 Use Cases](#-use-cases)
- [⭐ Summary](#-summary)

## Project Goals

- Provide a self-hosted host inspection tool
  
- Offer both **human-readable web UI** and **machine-readable API**

- Be lightweight, container-native, Kubernetes-ready, and easy to deploy

- Follow clean Go project structure (/cmd, /internal)
  
- Be easily extendable without architectural rewrites
  
- Be safe by default (health checks, middleware, rate limiting)

## What Information Does It Expose?

HostInfo provides structured host-level telemetry including:

- Hostname
  
- OS & kernel details
  
- CPU model & core count
  
- Memory information
  
- Environment variables (optional)
  
- Go runtime metadata
  
- Network interface information
  
- Cloud provider detection (AWS, GCP, Azure, local)

- Container runtime detection (kubernetes, docker, bare-metal)

## 🧱 Architecture Overview

HostInfo follows a clean layered architecture:

| Layer             | Description                          |
| ----------------- | ------------------------------------ |
| Entry Point       | `cmd/server`                         |
| Server Layer      | HTTP server setup, graceful shutdown |
| Middleware        | Logging, CORS, rate limiting         |
| API Layer         | Versioned REST endpoints (`/api/v1`) |
| Internal Services | System & runtime inspection          |
| Health System     | Liveness & readiness checks          |
| Configuration     | Environment-driven config            |

High-level flow:

```
Client
   ↓
HTTP Server
   ↓
Middleware Stack
   ↓
API Router (/api/v1)
   ↓
Internal Modules (runtime, system, cloud, network)
```

## 🗂 Repository Structure

```bash
hostinfo/
├── cmd/
│   └── server/              # Application entrypoint
├── internal/                # Core application logic (private)
│   ├── api/                 # API layer (v1 + middleware)
│   ├── health/              # Health check system
│   ├── server/              # HTTP server logic
│   └── custom/              # Configuration handling
├── frontend/                # Static assets (dashboard)
├── helm/                    # Kubernetes Helm chart
├── tools/                   # Dev automation scripts
│   └── scripts/git_hooks/   # Git hook system
├── docs/                    # Project documentation
├── Dockerfile
├── docker-compose.yml
├── go.mod / go.sum
└── CHANGELOG.md
```
This structure follows Go best practices:

- `/cmd` for entrypoints
  
- `/internal` for encapsulated logic
  
- Modular API versioning
  
- Middleware Layer - CORS handling, Request logging, Rate limiting

## Cloud Detection Layer

HostInfo safely checks metadata endpoints and environment characteristics:

| Provider | Detection Method                                   |
| -------- | ----------------------------------------           |
| AWS      | IMDS (169.254.169.254)                             |
| GCP      | metadata.google.internal                           |
| Azure    | Azure metadata endpoint (169.254.169.254/metadata) |
| Docker   | cgroup / container detection                       |
| Local    | Fallback when no cloud metadata detected           |

Failures degrade gracefully to `local` — the service never crashes due to metadata detection.

## 📦 Deployment Models

HostInfo supports multiple deployment targets:

-  **Local Binary**
   -  For Linux/macOS/Windows workstations.

-  **Docker Container**
   -  For servers, homelabs, CI automation. 

- **Docker Compose**
  - Part of larger observability stacks.

- **Systemd Service (Optional)**
  - For persistent Linux deployments.

- **Kubernetes**
  - Ready-to-use Helm chart
  - Ready-to-use kubectl

## 🛠 Developer Experience

The tools/ directory provides:

- `prepare.sh` — environment bootstrap

- `build.sh` — build automation

- Git hooks (pre-commit, commit-msg)

- Version validation for CLI tools

- Clean developer onboarding workflow

## 🧑‍💻 Target Users

HostInfo is designed for:

- DevOps / SRE engineers
  
- Platform / Infra  engineers

- System administrators

- Observability stack maintainers
  
- Automation systems
  
- Homelab enthusiasts
  
- Edge device operators

## 🪄 Use Cases

Common usage patterns:

- Quick inspection of remote servers
  
- Cloud environment detection
  
- Infrastructure baseline verification

- Lightweight system auditing (alternative to full monitoring stacks)
  
- Debugging container runtime environments

- Edge or IoT device introspection

## ⭐ Summary

It bridges the gap between simple system commands and Full-scale monitoring stacks, offering a clean, lightweight solution for machine telemetry.

It is:

✔ Fast
✔ Portable
✔ Single-binary deployable
✔ Container-ready
✔ Kubernetes-ready
✔ Cloud-native
✔ API-friendly
✔ DevOps-oriented

HostInfo is designed to be simple enough for quick diagnostics, yet structured enough for any environments.