# Quick reference

* **Maintained by:**  
  [Maksym Leus](https://github.com/MaksymLeus)

* **Where to get help or file issues:**  
  [GitHub Issues](https://github.com/MaksymLeus/hostinfo/issues)

* **Supported architectures:**  
  `amd64`, `arm64`

* **Source of this description:**  
  [docs repo's / directory](https://github.com/MaksymLeus/hostinfo/tree/main/docs) ([history](https://github.com/MaksymLeus/hostinfo/commits/main/docs))

# Supported tags and respective `Dockerfile` links:
  - [`latest`, `v1.0.0`](https://github.com/MaksymLeus/hostinfo/blob/main/Dockerfile)
  
---

# What is HostInfo?

**HostInfo** is a lightweight, container-ready Go service that exposes structured system, runtime, and environment information through a Web dashboard and a versioned REST API (`/api/v1`).

It is designed for DevOps engineers, SREs, platform teams, automation systems, and homelab operators who need fast, portable host telemetry with zero external runtime dependencies.

![Dashboard](https://github.com/MaksymLeus/hostinfo/raw/main/docs/images/Dashboard.png)

### Key Features
* 🚀 Fast & minimal Go HTTP server
* 🌐 Clean Web dashboard & interactive Swagger API Docs (`/swagger/index.html`)
* ☁️ Cloud provider detection (AWS / GCP / Azure / Local)
* ☸️ Kubernetes-ready with built-in liveness & readiness probes
* 🔒 Zero external runtime dependencies

---

## How to use this image

### Run via Docker CLI
```bash
docker run -d \
  --name hostinfo \
  -p 8080:8080 \
  maximleus/hostinfo:latest
```
Access the dashboard at `http://localhost:8080` and Swagger docs at `http://localhost:8080/swagger/index.html`.

### Run via Docker Compose
```yaml
version: '3.9'

services:
  hostinfo:
    image: maximleus/hostinfo:latest
    container_name: hostinfo
    ports:
      - "8080:8080"
    environment:
      - HOSTINFO_PORT=8080
      - HOSTINFO_DEBUG=false
    restart: unless-stopped
```
Run `docker compose up -d` to start the service.

### Deploy to Kubernetes using Helm
```bash
helm upgrade --install hostinfo ./helm \
  --namespace hostinfo \
  --create-namespace
```

## Environment Variables

| Variable         | Default   | Description       |
| ---------------- | --------- | ----------------- |
| `HOSTINFO_PORT`  | `8080`    | Port to listen on inside the container |
| `HOSTINFO_HOST`  | `0.0.0.0` | Network interface bind address      |
| `HOSTINFO_DEBUG` | `false`   | Enable verbose debug logs |

For advanced configuration, check the [Configuration Guide](https://github.com/MaksymLeus/hostinfo/blob/main/docs/CONFIGURATION.md).

## License
View [license information](https://github.com/MaksymLeus/hostinfo/blob/main/LICENSE.md) for the software contained in this image.

As with all Docker images, these likely also contain other software which may be under other licenses (such as Go binaries, base OS packages, etc., along with any direct or indirect dependencies of the primary software being contained).

It is the image user's responsibility to ensure that any use of this image complies with any relevant licenses for all software contained within.