
# Getting Started

This guide will help you get HostInfo running on your machine using the quickest and simplest methods.  
Whether you prefer to build from source or use Docker, this document has you covered.

---

## ✅ Prerequisites

HostInfo supports multiple environments, but depending on your preferred setup, you may need:

### **For local builds**
- Go **1.22+**
- Git

### **For containerized deployments**
- Docker **20+**
- (Optional) Docker Compose **v2+**

### **Supported Operating Systems**
| OS       | Status      |
|----------|------------|
| Linux    | ✔ Supported |
| macOS    | ✔ Supported |
| Windows  | ✔ Supported |

No external Go dependencies are required for the default build.

---

## 📦 Step 1 — Clone Repository
```bash
git clone https://github.com/yourname/hostinfo.git
cd hostinfo
```
If you plan to contribute, fork first and clone your fork instead.

## 🧱 Step 2 — Choose an Installation Method

### HostInfo can be run in three ways:

#### Option A — 🏁 Local Go Build (Recommended for Dev)
This method compiles and runs the binary directly.

**Build**
```bash
go build -o hostinfo ./cmd/server
```
**Run**
```bash
./hostinfo
```
Open your browser and visit: `http://localhost:8080`

#### Option B — 🐳 Docker Container (No Go Needed)
**Build image**
```bash
docker build -t hostinfo:latest ./docker
```
**Run container**
```bash
docker run -it --rm -p 8080:8080 hostinfo:latest
```
Open your browser and visit: `http://localhost:8080`

#### Option C — 🐙 Docker Compose (Services Setup)

This method is useful for multi-service environments.

**Start**
```bash
docker compose up -d
```
**Stop**
```bash
docker compose down
```
Open your browser and visit: `http://localhost:8080`


## ⚙️ Configuration Basics
HostInfo can be configured via:
- Environment variables
- CLI flags

Defaults are sensible, so configuration is optional for onboarding.

**Example CLI overrides**
```bash
./hostinfo --port 9090 --debug
```

## 🧪 Verify Functionality
After running HostInfo by any method, verify that:

**✔ Web UI is up**

Visit: `http://localhost:8080`
You should see a system information dashboard.


**✔ API responds**
```bash
curl http://localhost:8080/api/v1/info | jq
```
Expected output (varies per system):
```json
{
  "hostname": "myhost",
  "os": "linux",
  "cpu": {
    "cores": 8
  },
  "memory": "16GB"
}
```
## 🧹 Cleaning Up (Optional)

**Docker container cleanup**
```bash
docker ps -a
docker rm <container_id>
docker rmi hostinfo:latest
```
**Binary cleanup**
```bash
rm hostinfo
```

## 🙋 Need More Details?
Additional documentation is available in:
- `02-installation.md`
- `03-configuration.md`
- `04-usage.md`
- `05-api.md`

## 🎉 You're Ready!

You now have HostInfo running — either as a native binary or inside a Docker container.
Next steps depend on your use case:

➡ For config options → read `03-configuration.md`
➡ For API usage → read `05-api.md`
➡ For development → read `07-development.md`
