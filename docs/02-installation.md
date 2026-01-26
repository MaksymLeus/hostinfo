# Installation

This document describes how to install and build the **hostinfo** web application across multiple environments.

## Table of Contents
- [Installation](#installation)
  - [Table of Contents](#table-of-contents)
  - [1. 📦 Prerequisites](#1--prerequisites)
  - [2. 🔽 Clone Repository](#2--clone-repository)
  - [3. Backend Dependencies](#3-backend-dependencies)
  - [4. 🛠️ Local Build Options](#4-️-local-build-options)
    - [Option A — Using build.sh (recommended)](#option-a--using-buildsh-recommended)
    - [Option B — Manual go build](#option-b--manual-go-build)
  - [5. Platform-Specific Builds](#5-platform-specific-builds)
    - [Linux (from macOS/Windows)](#linux-from-macoswindows)
  - [6. ⚙️ Install to GOPATH](#6-️-install-to-gopath)
  - [7. Deployment Targets](#7-deployment-targets)
    - [Linux Server (manual)](#linux-server-manual)
    - [systemd (optional)](#systemd-optional)
  - [8. Docker Installation](#8-docker-installation)
    - [Build locally:](#build-locally)
  - [9. Docker Compose](#9-docker-compose)
  - [10. Helm Chart — Installation](#10-helm-chart--installation)
    - [Install Release](#install-release)
    - [Modify helm chart using values](#modify-helm-chart-using-values)
    - [Verify](#verify)
    - [Generate kubernetes configs](#generate-kubernetes-configs)
  - [11. ❗ Troubleshooting](#11--troubleshooting)
  - [12. 🧼 Uninstallation](#12--uninstallation)
  - [📄 License](#-license)

## 1. 📦 Prerequisites

Required:
You need the following tools installed:

| Tool | Purpose | Check |
|---|---|---|
| **Go** `>= 1.21` | Build & run hostinfo | `go version` |
| **Git** | Clone repository | `git --version` |

Optional but recommended (for deployment / CI):

- `make` — to use Makefile automation (if present)
- `docker` — for container builds (optional)
- `docker compose`
- `scp` (for remote deployment)
- `systemd` (Linux service)

## 2. 🔽 Clone Repository

```bash
git clone https://github.com/MaksymLeus/hostinfo.git
cd hostinfo
```

## 3. Backend Dependencies

```bash
go mod download
```

## 4. 🛠️ Local Build Options

### Option A — Using build.sh (recommended)

```bash
./build.sh
```

This produces a binary based on host OS/architecture, e.g.:
```bash
./bin/hostinfo
```

### Option B — Manual go build

```bash
go build -o ./bin/hostinfo ./cmd/server
```

Run it:
```bash
./bin/hostinfo
```
Open your browser and visit: `http://localhost:8080`

## 5. Platform-Specific Builds

### Linux (from macOS/Windows)

```bash
GOOS=linux GOARCH=amd64 go build -o hostinfo-linux-x64 ./cmd/server
```
Available arch targets include:

- `amd64`
- `arm64`

## 6. ⚙️ Install to GOPATH

```bash
go install ./cmd/server
```

Binary installs into:

```
$(go env GOPATH)/bin/hostinfo
```

Add GOPATH bin to PATH (if missing):

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## 7. Deployment Targets

### Linux Server (manual)

```bash
./build.sh all
scp bin/hostinfo-linux-x64 user@server:/opt/hostinfo/hostinfo
```

Run on server:

```bash
ssh user@server
cd /opt/hostinfo
./hostinfo
```

### systemd (optional)

Create `/etc/systemd/system/hostinfo.service`:

```ini
[Unit]
Description=Web Hostinfo Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/hostinfo
ExecStart=/opt/hostinfo/hostinfo
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Enable + start:

```bash
sudo systemctl enable hostinfo
sudo systemctl start hostinfo
```

## 8. Docker Installation

Hostinfo can be containerized using Docker.

### Build locally:

```bash
docker build -t hostinfo . -f /docker/Dockerfile
```

Run:

```bash
docker run -p 8080:8080 hostinfo
```
Open your browser and visit: `http://localhost:8080`


---

## 9. Docker Compose

```bash
docker compose build -f ./docker/docker-compose.yml
docker compose up -d
docker compose logs -f
```

Stop:

```bash
docker compose down
```

## 10. Helm Chart — Installation

<!-- ### Add Helm Repository
helm repo add <repo-name> <repo-url>
helm repo update -->

### Install Release
```bash
helm upgrade --install hostinfo ./helm <repo-name>/<chart-name> \
  --namespace <namespace> --create-namespace
```
### Modify helm chart using values
```bash
# inside ./helm/values.yaml you can see an example of how to change the helm release as you needed
helm upgrade --install hostinfo ./helm -f ./helm/values.yaml 
```

### Verify
kubectl get pods -n <namespace>
kubectl get svc -n <namespace>

### Generate kubernetes configs
If you whant run using `kubectl` you can generate template using helm
```bash
helm template hostinfo ./helm 
```

## 11. ❗ Troubleshooting

| Issue | Cause | Resolution |
|---|---|---|
| `go: no such file or directory` | Wrong working directory | `cd hostinfo` |
| `exec format error` | Wrong GOARCH/GOOS | Rebuild with correct target |
| Docker build slow | No layer cache | Enable BuildKit |
| `permission denied` | Missing execute flag | `chmod +x hostinfo` or `chmod +x build.sh` |

---

## 12. 🧼 Uninstallation

Remove binary:

```bash
rm -f hostinfo
```

Remove build artifacts:

```bash
rm -rf bin/
```

Remove systemd service:

```bash
sudo systemctl disable --now hostinfo
sudo rm -f /etc/systemd/system/hostinfo.service
```

Remove helm release:
```bash
#Check what’s installed
helm list -A
#remove (uninstall) a Helm release
helm uninstall <release-name> -n <namespace>
```

## 📄 License

MIT — see [`LICENSE.md`](../LICENSE.md) for details.
