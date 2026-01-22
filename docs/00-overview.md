
# HostInfo — Project Overview

HostInfo is a lightweight Go-based service that exposes system information through an HTML dashboard and a JSON API.  
It is designed for DevOps engineers, SREs, system administrators, and automation pipelines that need quick and portable access to machine telemetry.

## 🎯 Project Goals

- Provide a self-hosted tool to inspect system information
- Offer both **human-readable web UI** and **machine-readable API**
- Be lightweight, dependency-free, and easy to deploy
- Support Docker, Docker Compose, and standard Linux service management
- Integrate cleanly into CI/CD and automation environments

## 🧩 Key Features

- 🌐 **Web Dashboard** — view host info in a browser
- 📡 **REST API** — extract data programmatically
- 🐳 **Container-Ready** — minimal Docker image support
- 📦 **Binary or Docker Deployment**
- ⚙️ **Configurable** via env or CLI flags
- 🔐 **Zero External Dependencies**
- 📂 **Clean Codebase** with docs + CI

## 🖥 What Information Does It Show?

HostInfo exposes hardware and OS metrics such as:

- Hostname
- OS & kernel details
- CPU model & core count
- Memory information
- Envierment veriabels
- Disk space (planned)
- Network details (planned)

These can be improved or extended over time (CPU usage %, disk IO, network stats, etc).

## 🧱 Architecture Overview

HostInfo is structured as a simple web server with the following logical layers:

| Layers | Description |
|:---|:---|
| Web UI  | ➜ HTML Templates |
| HTTP API | ➜ JSON Responses |
| System Information | ➜ OS / CPU / Memory |
| Runtime | ➜ Go 1.22+ |


## 🗂 Repository Structure (High-Level)
```bash
hostinfo/
├── cmd/server # Application entrypoint
├── internal/ # Core internal logic
├── web/ # HTML templates, static assets
├── docker/ # Docker + Compose files
├── docs/ # Documentation
├── .github/workflows/ # CI/CD pipelines
├── scripts/hooks/ # Git hooks (pre-commit, commit-msg)
└── tools/ # Helper scripts (bootstrap, dev)
```

More details in: `07-development.md`

## 🛠 Tech Stack

| Category    | Choice        |
|-------------|---------------|
| Language    | Go (1.22+)    |
| Runtime     | Standard Lib  |
| UI          | HTML Templates|
| Packaging   | Docker & Go   |
| CI/CD       | GitHub Actions|
| Release     | semantic-release|

No external dependencies are required for core features.

## 📦 Deployment Models

### HostInfo supports multiple deployment targets:

-  **Local Binary**
   -  For Linux/macOS/Windows workstations.

-  **Docker Container**
   -  For servers, homelabs, CI automation. 

- **Docker Compose**
  - Part of larger observability stacks.

- **Systemd Service (Optional)**
  - For persistent Linux deployments.

Kubernetes deployment manifests may be added later.

## 🧑‍💻 Target Users

HostInfo is intended for:

- DevOps / SREs
- Platform / Infra engineers
- System administrators
- Automation pipelines
- Observability stack maintainers
- Homelab enthusiasts

## 🪄 Use Cases

Common usage patterns include:

- Checking remote machine details via web browser
- Collecting telemetry in CI/CD jobs
- Integrating system info into dashboards
- Baseline validation for provisioning
- Self-hosted server inventory in homelabs
- Lightweight monitoring for edge devices

## 🏁 Project Status

**Current Stage:** Early Development  
Core features are functional, and additional modules (disk, network, metrics, authentication) are planned.

Upcoming enhancements include:

- Disk usage collection
- Network stats
- Auth (optional Basic/OAuth)
- Metrics (Prometheus endpoint)
- Swagger/OpenAPI API reference
- Improved UI front-end

---

## 📜 License

This project is licensed under the **MIT License**.  
See `LICENSE.md` for full details.

---

## 🤝 Contributing

Contributions are welcome!  
Please see:

- `docs/07-development.md`
- `docs/09-releasing.md` (if semantic-release involved)

---

## ⭐ Summary

HostInfo aims to streamline the collection of host-level data in a way that is:

✔ Fast  
✔ Portable  
✔ Single-binary  
✔ Self-hosted  
✔ API-friendly  
✔ DevOps-ready

It bridges the gap between “simple Linux commands” and “full monitoring stacks” by providing a clean, lightweight, and flexible utility for machine introspection.




