<h1 align="center">HostInfo</h1>
<p align="center">Lightweight Go-based Host Telemetry Service</p>

<p align="center">
  <a href="https://github.com/MaksymLeus/hostinfo/releases">
    <img src="https://img.shields.io/github/v/release/MaksymLeus/hostinfo?style=for-the-badge">
  </a>
  <a href="https://github.com/MaksymLeus/hostinfo/actions/workflows/ci.yml">
    <img src="https://img.shields.io/github/actions/workflow/status/MaksymLeus/hostinfo/ci.yml?style=for-the-badge&label=CI">
  </a>
  <a href="https://hub.docker.com/r/maximleus/hostinfo">
    <img src="https://img.shields.io/docker/pulls/maximleus/hostinfo?style=for-the-badge">
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/github/license/MaksymLeus/hostinfo?style=for-the-badge">
  </a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go">
  <img src="https://img.shields.io/badge/Arch-amd64%20%7C%20arm64-orange?style=for-the-badge">
</p>

<p align="center">
  <b>HostInfo</b> is a lightweight, container-ready Go service that exposes structured system, runtime, and environment information through a Web dashboard and a versioned REST API.
</p>

It is designed for **DevOps engineers, SREs, platform teams, automation systems, and homelab operators** who need fast, portable host telemetry.

![Dashboard](docs/images/Dashboard.png)

## ✨ Features

- 🚀 Fast & minimal Go HTTP server
- 🌐 Clean Web dashboard
- 📡 Versioned REST API (`/api/v1`)
- 📊 Prometheus Metrics endpoint (`/metrics`)
- 🔔 Webhook / Alerting notifications for high CPU/Memory
- 📚 Interactive Swagger API Documentation (`/swagger/index.html`)
- ☁️ Cloud detection (AWS / GCP / Azure / Local)
- 🧱 Clean layered architecture
- 📚 Structured documentation

See [`docs/OVERVIEW.md`](docs/OVERVIEW.md) for full overview.

## 🚀 Quick Start

### Docker Compose (Recommended)
```bash
docker compose up -d
```

### Docker
```bash
docker run -p 8080:8080 maximleus/hostinfo:latest
```

### From Source
```bash
# Build
./tools/build.sh

# Run
./bin/hostinfo
```

Access at: **http://localhost:8080**

Swagger UI: **http://localhost:8080/swagger/index.html**

### ☸️ Deploy with Helm
```bash
helm upgrade --install hostinfo ./helm \
  --namespace hostinfo \
  --create-namespace
```
See [`docs/DEPLOYMENT.md`](docs/DEPLOYMENT.md) for full deployment guide.


## Screenshots

| Network | Raw Json |
|:---:|:---:|
| ![db_env](docs/images/db_network.png) | ![db_raw_json](docs/images/db_raw_json.png) |

| Command - Ping | Command List |
|:---:|:---:|
| ![nu_ping](docs/images/nu_ping.png) | ![nu_list](docs/images/nu_list_cmd.png) |

## 🔌 API Examples  & Swagger UI
HostInfo exposes a REST API for system, cloud, and network diagnostics. 
You can interactively test all endpoints using the built-in **Swagger UI** at:
📌 **`http://localhost:8080/swagger/index.html`**

## ⚙️ Configuration

| Variable         | Default   | Description       |
| ---------------- | --------- | ----------------- |
| `HOSTINFO_PORT`  | `8080`    | Port to listen on |
| `HOSTINFO_HOST`  | `0.0.0.0` | Bind address      |
| `HOSTINFO_DEBUG` | `false`   | Enable debug mode |

See [`docs/CONFIGURATION.md`](docs/CONFIGURATION.md) for full configuration reference.

## 📚 Documentation

Project documentation lives inside the `/docs` directory:

| Document           | Purpose                                 |
| ------------------ | --------------------------------------- |
| `OVERVIEW.md`      | Architecture, goals, internal structure |
| `CONFIGURATION.md` | Environment variables & configuration   |
| `DEPLOYMENT.md`    | Docker, Compose, Helm, Kubernetes       |
| `DEVELOPMENT.md`   | Local development workflow              |
| `TODO.md`          | Roadmap & planned improvements          |

## Tech Stack

**Backend:** Go 1.24+, Echo, pro-bing, gopsutil

**Frontend:** Node 23.13+, NPM 11.6+, Vite

## Contributing
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: Add amazing feature'`)
4. If you modify API endpoints, regenerate Swagger docs:
   ```bash
   swag init -g cmd/server/hostinfo.go --parseInternal --parseDepth 2
   ```
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a [Pull requests](https://github.com/MaksymLeus/hostinfo/pulls)
   
See [`docs/DEVELOPMENT.md`](docs/DEVELOPMENT.md) for development setup.

## 📄 License
MIT License — see [`LICENSE.md`](LICENSE.md) for details.

## Support

- [GitHub Issues](https://github.com/MaksymLeus/hostinfo/issues)