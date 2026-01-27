<h1 align="center">hostinfo</h1>
<p align="center">Minimal Go-based system & environment info API</p>

<p align="center">
  <a href="https://github.com/MaksymLeus/hostinfo/releases"><img src="https://img.shields.io/github/v/release/MaksymLeus/hostinfo?style=for-the-badge"></a>
  <a href="https://github.com/MaksymLeus/hostinfo/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/MaksymLeus/hostinfo/ci.yml?style=for-the-badge&label=CI"></a>
  <a href="https://hub.docker.com/r/maximleus/hostinfo"><img src="https://img.shields.io/docker/pulls/maximleus/hostinfo?style=for-the-badge"></a>
  <a href="LICENSE"><img src="https://img.shields.io/github/license/MaksymLeus/hostinfo?style=for-the-badge"></a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.20+-00ADD8?style=for-the-badge&logo=go">
  <img src="https://goreportcard.com/badge/github.com/MaksymLeus/hostinfo?style=for-the-badge">
  <img src="https://img.shields.io/badge/Arch-amd64%20%7C%20arm64-orange?style=for-the-badge">
</p>

<p align="center">
  <b>HostInfo</b> is a lightweight Go service that exposes server/system information through both an HTML interface and a RESTful JSON API.  
</p>

It’s designed for DevOps/debugging scenarios, observability dashboards, and automation.

![Dashboard](docs/images/Dashboard.png)

## ✨ Features

- 🚀 Fast & lightweight Go server
- 🌐 Simple web UI
- 📡 JSON API support
- 🐳 Docker & Compose ready
- 🔒 Zero external dependencies
- 📦 CI/CD & Semantic Release compatible
- 📁 Clean repo & docs structure

See [`docs/00-overview.md`](docs/00-overview.md) for a complete overview of the project.

## Quick Start
Hostinfo is available as a Docker image for easy deployment.

**Image Details:**
- **Registry:** Docker Hub ([`maximleus/hostinfo`](https://hub.docker.com/r/maximleus/hostinfo))
- **Base Image:** `debian:bookworm-slim` 
- **Platforms:** `linux/amd64`, `linux/arm64`
- **Size:** ~47MB compressed

# Clone 
```sh
git clone https://github.com/MaksymLeus/hostinfo.git
cd  hostinfo
```
### Docker Compose Run (Recommended)
```bash
docker compose -f docker/docker-compose.yml up -d
```
### Docker Run
```bash
docker pull maximleus/hostinfo:latest
docker run -p 8080:8080 hostinfo:latest
```

### From Source
```bash
# Build
./tools/build.sh
# Run
./bin/hostinfo
```
Access at: **http://localhost:8080**


### Using Helm
```bash
helm upgrade --install hostinfo ./helm
```

For more detailed instructions on getting started, see [`docs/01-getting-started.md`](docs/01-getting-started)

Additional installation documentation is available in [`docs/02-installation.md`](docs/02-installation.md)


## 📚 Documentation

Full project documentation is available in the [`/docs`](./docs/) directory:

| Document | Description |
|---|---|
| [`00-overview.md`](./docs/00-overview.md) | Project overview, goals & scope |
| [`01-getting-started.md`](./docs/01-getting-started.md) | Minimum requirements & quick usage |
| [`02-installation.md`](./docs/02-installation.md) | Installation paths (binary, source, Docker) |
| [`03-usage.md`](./docs/04-usage.md) | Web UI, CLI, API & environment vars |
| [`04-api.md`](./docs/07-api.md) | REST API endpoints, schemas & examples |
| [`05-architecture.md`](./docs/05-architecture.md) | Internal packages, components & data flow |
| [`06-deployment.md`](./docs/03-deployment.md) | Docker, compose, Linux, systemd & cloud deployment |
| [`07-cloud-detection.md`](./docs/06-cloud-detection.md) | AWS/GCP/Azure metadata detection logic |
| [`08-git-hooks.md`](./docs/08-git-hooks.md) | Local commit enforcement & checks |
| [`09-releasing.md`](./docs/09-releasing.md) | Semantic-release, versioning & CI/CD pipelines |
| [`10-development.md`](./docs/10-development.md) | Contains guidelines for setting up and maintaining a local development environment, tooling requirements, recommended workflows, and environment configuration. |

✔ Supports developers, DevOps engineers & cloud operators  
✔ Clean separation between **usage**, **deployment** & **internals**  
✔ Ideal for open-source onboarding and maintainability


## API (in-progress)
```bash
# Health check using curl
curl http://localhost:8080/healthz
```
See [`docs/04-api.md`](docs/04-api.md) for complete documentation.

## Configuration

### Essential Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HOSTINFO_PORT` | `8080` | Listen port |
| `HOSTINFO_ADDR` | `0.0.0.0` | Bind address |
| `HOSTINFO_DEBUG` | false | Basic Debug param |

<!-- See [docs/CONFIGURATION.md](docs/CONFIGURATION.md) for all options. -->

## 🤝 Contributing
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a [Pull requests](https://github.com/MaksymLeus/hostinfo/pulls)
   
See [`docs/07-development.md`](docs/07-development.md) for development setup.


## 📄 License
MIT License — see [`LICENSE.md`](LICENSE.md) for details.