# 🖥️ HostInfo

<p>
  <a href="#"><img src="https://img.shields.io/github/v/release/MaksymLeus/hostinfo?style=flat-square" alt="release"></a>
  <a href="#"><img src="https://img.shields.io/github/actions/workflow/status/MaksymLeus/hostinfo/ci.yml?style=flat-square&label=CI" alt="ci status"></a>
  <a href="#"><img src="https://img.shields.io/github/license/MaksymLeus/hostinfo?style=flat-square" alt="license"></a>
  <a href="#"><img src="https://img.shields.io/codecov/c/github/MaksymLeus/hostinfo?style=flat-square" alt="coverage"></a>
  <a href="#"><img src="https://img.shields.io/docker/pulls/maximleus/hostinfo?style=flat-square" alt="docker pulls"></a>
  <a href="#"><img src="https://img.shields.io/badge/semantic-release-enabled-e10079?logo=semantic-release&style=flat-square" alt="semantic-release"></a>
</p>

**HostInfo** is a lightweight Go service that exposes server/system information through both an HTML interface and a RESTful JSON API.  
It’s designed for DevOps/debugging scenarios, observability dashboards, and automation.

---





## ✨ Features

- 🚀 Fast & lightweight Go server
- 🌐 Simple web UI
- 📡 JSON API support
- 🐳 Docker & Compose ready
- 🔒 Zero external dependencies
- 📦 CI/CD & Semantic Release compatible
- 📁 Clean repo & docs structure

---

## 📦 Installation

### Option A — Local Build (Go)

```bash
git clone https://github.com/MaksymLeus/hostinfo.git
cd hostinfo
go build ./cmd/server
./server
```

### Option B — Docker
```bash
docker build -t hostinfo:latest ./docker
docker run -p 8080:8080 hostinfo:latest
```

### Option C — Docker Compose
```bash
docker compose up -d
```

## 🏃 Getting Started
```bash
./server
```
Open in browser: `http://localhost:8080`

You’ll see system info like CPU, RAM, OS, disk, hostname, etc.

## 🧰 Usage
### Web UI

Provides a human-friendly view of system info.

### JSON API
```http
GET /api/v1/info
```
### Sample response:
```http
{
  "hostname": "mylaptop",
  "cpu": {
    "cores": 8,
    "model": "Intel(R) Core(TM) i7"
  },
  "memory": "16GB",
  "os": "Linux"
}
```
## ⚙️ Configuration

### Environment variables:

| Variable | Default | Description |
|:---|:---:|:---|
| `PORT` | `8080` | HTTP port |
| `DEBUG` | `false` | Debug logs |

### CLI flags:

```bash
./server --port 9090 --debug
```

## 🧱 Repository Structure

```bash
hostinfo
├── cmd/server
├── internal
├── web
├── docker
├── docs
├── .github/workflows
├── scripts/hooks
├── tools
└── ...
```
For full breakdown, see: [`docs/00-overview.md`](docs/00-overview.md)

## ⚙️ Development

### Prerequisites:
- Go 1.22+
- Docker (optional)
- Linux/macOS/Windows

### Run Tests
```bash
go test ./...
```
### Git Hooks (pre-commit & commit-msg)
```bash
./tools/setup-hooks.sh
```
## 🚀 CI/CD & Releases
This project supports:

✔ GitHub Actions CI

✔ Docker image builds

✔ Semantic Versioning

✔ Automated changelog generation

Semantic Release is used for tagging & changelog:

- feat: → minor version
- fix: → patch version
- BREAKING CHANGE: → major version

## 🐳 Docker Image
When published, you’ll be able to pull:
```bash
docker pull maksymleus/hostinfo:latest
```
Or run:
```bash
docker run -p 8080:8080 MaksymLeus/hostinfo:latest
```
## 📚 Documentation
See the docs/ folder for:
- [Overview](./docs/00-overview.md)
- [Getting Started](./docs/01-getting-started.md)
- Installation 
- API Reference 
- CI/CD 
- Deployment 
- Releasing

## 🤝 Contributing
Contributions are welcome!

Steps:

1. Fork the project
2. Create your feature branch
3. Commit using Conventional Commits
4. Push & open PR

Please follow: [`docs/07-development.md`](docs/07-development.md)

## 📄 License
This project is licensed under the MIT License — see [`LICENSE.md`](LICENSE.md) for details.

## ⭐ Support

If you find this project helpful:
- Star the repo
- Open issues or feature requests
- Share with others

