# Installation

This document describes how to install and build the **hostinfo** web application across multiple environments.

---

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
Access in browser: `http://localhost:8080`

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

Access:

```
http://localhost:8080
```

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

---

## 10. ❗ Troubleshooting

| Issue | Cause | Resolution |
|---|---|---|
| `go: no such file or directory` | Wrong working directory | `cd hostinfo` |
| `exec format error` | Wrong GOARCH/GOOS | Rebuild with correct target |
| Docker build slow | No layer cache | Enable BuildKit |
| `permission denied` | Missing execute flag | `chmod +x hostinfo` or `chmod +x build.sh` |

---

## 11. 🧼 Uninstallation

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

---

## 📄 License

MIT — see `LICENSE.md` for details.
