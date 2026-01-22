# 07 — API

This document defines the **hostinfo HTTP API** for programmatic access.  
The API complements the web dashboard by exposing host, container, and cloud information as JSON.

---

## 1. Overview

- Base URL: `/`
- Current default: **dashboard HTML**
- Future JSON endpoints follow `/api/v1/...`
- Stateless and read-only
- No authentication required by default

> All endpoints are optional; failures degrade gracefully.  

---

## 2. API Versioning

Planned versioning:

```
/api/v1/
```

- `v1` = initial stable API
- Future versions: `v2`, `v3` for extensions or breaking changes

---

## 3. Endpoints

### 3.1 `/api/v1/info`

Returns a complete snapshot of the current host, container, and cloud metadata.

**Method:** `GET`  
**Response (JSON):**

```json
{
  "host": {
    "hostname": "my-host",
    "os": "linux",
    "arch": "amd64",
    "uptime": "5h32m",
    "goVersion": "go1.24"
  },
  "container": {
    "id": "a1b2c3d4",
    "runtime": "docker",
    "uptime": "5h30m"
  },
  "cloud": {
    "provider": "aws",
    "region": "us-east-1",
    "availabilityZone": "us-east-1a",
    "instanceType": "t3.micro",
    "instanceId": "i-1234567890abcdef"
  },
  "network": [
    {
      "name": "eth0",
      "mac": "02:42:ac:11:00:02",
      "ips": ["172.17.0.2"]
    }
  ],
  "env": {
    "PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
    "HOSTINFO_ENV": "production"
  }
}
```

---

### 3.2 `/api/v1/health`

Returns server health and uptime.

**Method:** `GET`  
**Response:**

```json
{
  "status": "ok",
  "uptime": "5h32m",
  "version": "1.0.0",
  "env": "production"
}
```

- Useful for monitoring / readiness probes
- Can be extended for liveness probes

---

### 3.3 Future Endpoints

| Endpoint | Description |
|----------|------------|
| `/api/v1/metrics` | Prometheus-compatible metrics |
| `/api/v1/cloud`   | Cloud metadata only |
| `/api/v1/container` | Container info only |
| `/api/v1/network` | Network interfaces and IPs |

---

## 4. Request/Response Conventions

- All requests: `GET`
- JSON only
- Response content-type:

```http
Content-Type: application/json; charset=utf-8
```

- Errors use standard HTTP codes:

| Code | Meaning |
|------|--------|
| 200  | Success |
| 400  | Bad request |
| 404  | Endpoint not found |
| 500  | Internal server error |

Example error:

```json
{
  "error": "cloud metadata not available"
}
```

---

## 5. Query Parameters (Future)

Potential optional query parameters:

- `?format=json` → force JSON response from dashboard
- `?fields=host,cloud` → return selected sections only
- `?timeout=500ms` → adjust metadata probe timeout

---

## 6. Security Considerations

- No authentication required by default
- Exposing environment variables may leak secrets
- Recommended to run behind reverse proxy with HTTPS and optional auth
- Ensure cloud metadata endpoints are internal-only when exposing API externally

---

## 7. Versioning & Deprecation

- All endpoints must specify API version in path (`/api/v1/...`)
- Future breaking changes increment version number (`v2`, `v3`)
- Deprecated endpoints respond with HTTP `410 Gone` and redirect message

---

## 8. Rate Limiting & Performance

- Currently **no rate limiting**
- Consider implementing per-IP throttling if exposed publicly
- Metadata probes are cached in-memory for performance
- Requests are handled asynchronously to prevent blocking

---

## 9. Example Usage

### curl Example

```bash
curl http://localhost:8080/api/v1/info
```

### Python Example

```python
import requests

resp = requests.get("http://localhost:8080/api/v1/info")
data = resp.json()
print(data["cloud"]["provider"])
```

---

## 10. Logging

- API access logs can be enabled
- Standard stdout logging
- No persistent storage

---

## License

MIT — see `LICENSE.md`.
