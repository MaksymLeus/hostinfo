# 04 — Usage

This document describes how to run, access, and interact with the **hostinfo** application in different environments.

---

## 1. Default Behavior

When started without arguments, hostinfo:

- runs an HTTP server on `:8080`
- renders the host information dashboard at `/`

Example:

```bash
./hostinfo
```

Access in browser:

```
http://localhost:8080
```

No configuration files or environment variables are required.

---

## 2. Supported UI Modes

The dashboard renders the following categories (depending on environment):

- **Host System**
- **Container Runtime**
- **Cloud Metadata**
- **Network Interfaces**
- **Environment Variables**
- **Runtime Details (Go / OS / Arch)**

If running inside:

- Docker → shows container metadata
- AWS EC2 → shows EC2 metadata
- GCP → shows project/zone/machine
- Azure → VM metadata
- Bare-metal → falls back to local host info

No cloud credentials are needed for metadata detection.

---

## 3. JSON API (If Enabled / Future)

If JSON output is enabled or supported in the future, examples may include:

```
GET /api/info
```

Response might include:

```json
{
  "hostname": "ip-10-0-1-15",
  "os": "linux",
  "arch": "amd64",
  "cloud": "aws",
  "region": "us-east-1",
  "instanceType": "t3.micro"
}
```

> If JSON output is not currently implemented, this section remains reserved for future support.

---

## 4. Environment Behavior

hostinfo auto-detects cloud/platform context:

| Platform | Detection |
|---|---|
| Docker | `/proc/self/cgroup`, container hostname |
| AWS EC2 | IMDS (169.254.169.254) |
| GCP | GCP metadata server |
| Azure | Azure Metadata Service |
| Local | no metadata endpoints |

All metadata requests are **safe**:
- timeouts are short
- failures do not crash the app

---

## 5. Port Configuration (If Provided)

If port configuration is supported later, recommended format:

```bash
./hostinfo --port 9090
```

Access:

```
http://localhost:9090
```

If not yet implemented, this section becomes future reserved.

---

## 6. Logging

hostinfo prints startup info to stdout:

Example:

```
[hostinfo] starting server on :8080
[hostinfo] environment: docker
[hostinfo] cloud: aws (us-east-1)
```

Logs can be collected via:

- systemd (`journalctl`)
- Docker logs
- Kubernetes logs

---

## 7. Production Notes

When running in production:

- run behind reverse proxy (nginx, traefik, caddy)
- restrict public access unless required
- enable TLS termination on proxy layer
- do not expose raw port `8080` to internet without TLS

---

## 8. Troubleshooting

| Issue | Cause | Resolution |
|---|---|---|
| Blank UI | Cloud metadata timeout | Wait ~1s or disable metadata calls |
| `bind: address already in use` | Port in use | Change port or stop other service |
| Docker no metadata | No cloud server | Expected inside local docker runs |
| Cloud info missing | Not running on cloud | Expected fallback |

---

## License

MIT — see `LICENSE.md`.
