# 05 — Architecture

This document outlines the architecture, execution model, core components, and cloud/container detection logic of the **hostinfo** application.

---

## 1. High-Level Overview

hostinfo is a **single-binary Go application** that:

- collects host/container/cloud metadata
- aggregates system + network information
- renders a web dashboard using templates
- performs safe non-blocking cloud metadata probing

No external services or databases are required.

---

## 2. Process Model

```
┌─────────────────────────┐
│  hostinfo (Go binary)   │
│                         │
│   ┌───────────────────┐ │
│   │ HTTP Server       │<── incoming web requests
│   └───────────────────┘ │
│   ┌───────────────────┐ │
│   │ Metadata Collect. │── system, runtime, env
│   └───────────────────┘ │
│   ┌───────────────────┐ │
│   │ Cloud Detection   │── AWS/GCP/Azure/local
│   └───────────────────┘ │
│   ┌───────────────────┐ │
│   │ Template Engine   │── render HTML dashboard
│   └───────────────────┘ │
└─────────────────────────┘
```

All components run in-process, no goroutine explosions, no IPC.

---

## 3. Source Code Layout

Relevant structure (simplified):

```
cmd/server/       → main entrypoint + HTTP bootstrap
internal/         → metadata collectors & helpers
web/templates/    → HTML layout for dashboard
web/image.png     → UI static asset
docker/           → Docker + compose definitions
.github/          → CI/CD workflows
```

Key directories are kept minimal to preserve clarity.

---

## 4. Component Breakdown

### 4.1 HTTP Server

Responsibilities:

- listen on `:8080` (default)
- serve dashboard at `/`
- render templates with collected data

### 4.2 Metadata Collector

Aggregates:

- hostname
- OS / arch / runtime
- environment variables
- uptime + process info
- network interfaces
- container markers
- cloud metadata (if available)

All calls are **non-blocking** with short timeouts.

### 4.3 Cloud Detection Layer

Checks metadata endpoints safely:

| Provider | Method |
|---|---|
| AWS | IMDSv1 (169.254.169.254) |
| GCP | metadata.google.internal |
| Azure | 169.254.169.254/metadata |
| Docker | cgroup/hostname heuristics |
| Local | fallback if no cloud metadata |

Failures never crash — they degrade to `local`.

### 4.4 Template Rendering

Backend populates a structured template context:

```go
type DashboardData struct {
    Host       HostInfo
    Container  ContainerInfo
    Cloud      CloudInfo
    Network    []NetInterface
    Runtime    GoRuntime
    Env        map[string]string
}
```

Rendered via `index.html` under `web/templates`.

---

## 5. Runtime Behavior

Startup sequence:

1. initialize collectors
2. detect execution environment
3. probe cloud metadata asynchronously
4. serve HTTP dashboard
5. update metrics on each request

### Execution Contexts

| Context | Behavior |
|---|---|
| Local bare-metal | fallback mode |
| Docker container | container info enabled |
| AWS EC2 | IMDS metadata |
| GCP | metadata server |
| Azure | IMDS metadata |
| Unknown | fallback |

Timeout defaults are small to avoid UI blocking.

---

## 6. Cloud & Container Detection Logic

Detection pipeline (simplified):

```
Container?  → Check cgroup / hostname patterns
     ↓ yes
   Docker

Cloud?      → Probe IMDS endpoints
     ↓ yes
   AWS / GCP / Azure

Else → Local fallback
```

Priority order:

```
Container → Cloud → Local
```

Cloud checks run in parallel to reduce latency.

---

## 7. Dependencies

### Language & Runtime

- **Go >= 1.21**

### External Services

None required.

### External Libraries

Minimal standard library usage preferred to keep binary small and portable.

---

## 8. Deployment Characteristics

- **stateless** (no DB, no local storage)
- **horizontal scale-friendly** (idempotent)
- **portable** across bare-metal, VMs, containers
- **read-only** by design (never mutates system)

Ideal for debugging ephemeral infra, CI workers, cloud VMs, or container clusters.

---

## 9. Security Considerations

hostinfo:

- does **not** require cloud credentials
- does **not** modify system state
- does **not** open outbound ports besides metadata probes
- does **not** expose OS secrets beyond env vars

Exposure risks:

- environment variables may contain secrets in certain deployments
- dashboard should not be public-facing without TLS/proxy

Recommendation:

Deploy behind Nginx / Traefik / Caddy for TLS + authentication.

---

## 10. Extensibility

Possible future modules:

- JSON API (`/api/info`)
- Prometheus `/metrics`
- CLI flags (`--port`, `--json`)
- WASM UI bundle
- Plugin architecture

Backend design keeps boundaries flexible:

```
Collector → Adapter → Renderer
```

No global state except basic runtime caches.

---

## License

MIT — see `LICENSE.md`.
