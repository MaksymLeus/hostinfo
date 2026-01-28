# 🖧 HostInfo API Documentation
The HostInfo API provides endpoints for retrieving system information and performing network diagnostics such as ping, DNS lookup, HTTP requests, and TCP port checks. All endpoints are served under /api/v1/.

All responses are in JSON format. Errors automatically return an error field.

### Base URL
```bash
http://<host>:<port>/api/v1/
```
- Default host: `0.0.0.0`
- Default port: `8080`
- Health check: `/healthz` (outside API)

Request Flow Diagram (High-Level):
```sh
┌──────────┐
│  Client  │
│ (curl)   │
└────┬─────┘
     │ HTTP Request
     ▼
┌──────────────┐
│  HTTP Router │
│ (chi / mux)  │
└────┬─────────┘
     │
     ▼
┌────────────────────┐
│ Rate Limit Middleware │
│  - Per IP limiter     │
│  - Token bucket       │
│  - Burst support      │
└────┬─────────────────┘
     │
     ├─❌ Limit exceeded
     │    └─► 429 JSON + X-RateLimit headers
     │
     ▼
┌──────────────┐
│ API Handler  │
│ (ping/dns/…) │
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ JSON Response│
└──────────────┘
```

## Rate Limiting


All API endpoints are protected by rate limiting per IP to prevent spam:

- Default: **15 requests per minute per IP**, burst 5 requests.

- Exceeding the limit returns **HTTP 429**:
  ```json
  {
    "error": "rate limit exceeded"
  }
  ```
- Standard headers returned for all requests:
  | Header                  | Description                                          |
  | ----------------------- | ---------------------------------------------------- |
  | `X-RateLimit-Limit`     | Maximum allowed requests per minute (burst)          |
  | `X-RateLimit-Remaining` | Number of requests left in the current window        |
  | `X-RateLimit-Reset`     | Unix timestamp when the next token will be available |
  
  Example curl to see headers:
  ```bash
  curl -s -D - "http://localhost:8080/api/v1/ping?host=google.com" -o /dev/null
  ```
  Output:
  ```http
  HTTP/1.1 200 OK
  X-RateLimit-Limit: 5
  X-RateLimit-Remaining: 4
  X-RateLimit-Reset: 1674928370
  Content-Type: application/json
  ```
  

Example Response When Rate Limited
```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1674928370

{"error":"rate limit exceeded"}
```


Configurable via: `internal/api/middleware/rate_limiter.go`
```go
var (
    rateLimit = rate.Every(time.Minute / 15) // 15 req/min
    burst     = 5                            // allow bursts
)
```
Rate Limiter Internals (Detailed):
```bash
              ┌─────────────────────┐
              │ Incoming Request IP │
              └─────────┬───────────┘
                        ▼
           ┌──────────────────────────┐
           │ clients[ip] lookup       │
           │ map[string]*client       │
           └─────────┬────────────────┘
                     │
        ┌────────────▼────────────┐
        │ Exists?                  │
        └──────┬──────────┬───────┘
               │ yes      │ no
               │          ▼
               │   ┌────────────────┐
               │   │ Create limiter │
               │   │ rate + burst   │
               │   │ lastSeen=now   │
               │   └────────────────┘
               ▼
┌──────────────────────────┐
│ limiter.Allow()          │
│ (token bucket algorithm)│
└─────────┬────────────────┘
          │
   ┌──────▼──────┐
   │ Allowed?    │
   └───┬────┬───┘
       │yes │no
       │    ▼
       │  429 Too Many Requests
       │  + X-RateLimit headers
       ▼
┌──────────────────────────┐
│ Update headers:          │
│ - Limit                  │
│ - Remaining              │
│ - Reset                  │
└─────────┬────────────────┘
          ▼
┌──────────────────────────┐
│ Forward to API handler   │
└──────────────────────────┘
```
Cleanup Goroutine (Optional Diagram):
```bash
┌────────────────────────────┐
│ Background Cleanup Loop    │
│ (every 5 minutes)          │
└─────────┬──────────────────┘
          ▼
┌────────────────────────────┐
│ Iterate clients map        │
└─────────┬──────────────────┘
          ▼
┌────────────────────────────┐
│ lastSeen > TTL ?           │
└──────┬────────────┬────────┘
       │ yes        │ no
       ▼            ▼
 Delete limiter   Keep client
```
## Endpoints

### 1. Health Check
```bash
GET /healthz
```
Response:
```json
{
  "status": "ok"
}
```
Example:
```bash
curl -s http://localhost:8080/healthz | jq
```

### 2. Ping
```bash
GET /api/v1/ping?host=<hostname>
```

**Description**: Ping a host to check availability and latency.

Query Parameters:
| Parameter | Type   | Required | Description            |
| --------- | ------ | -------- | ---------------------- |
| host      | string | Yes      | Hostname or IP to ping |

Response:
```json
{
  "host": "google.com",
  "packets_sent": 3,
  "packets_recv": 3,
  "loss_percent": 0,
  "min_rtt": "12.345ms",
  "max_rtt": "14.678ms",
  "avg_rtt": "13.456ms"
}
```
Error Response Example:
```json
{"error":"host parameter required"}
```
Curl Example:
```sh
curl -s "http://localhost:8080/api/v1/ping?host=google.com" | jq
```
### 3. DNS Lookup
```bash
GET /api/v1/dns?host=<hostname>
```
**Description**: Resolve a hostname to its IP addresses and CNAME.

Query Parameters:
| Parameter | Type   | Required | Description         |
| --------- | ------ | -------- | ------------------- |
| host      | string | Yes      | Hostname to resolve |

Response:
```json
{
  "host": "google.com",
  "ips": ["142.250.72.238", "142.250.72.206"],
  "cname": "google.com."
}
```
Curl Example:
```sh
curl -s "http://localhost:8080/api/v1/dns?host=google.com" | jq
```
### 4. HTTP Request / Curl
```bash
GET /api/v1/curl?url=<url>
```
**Description**: Perform an HTTP GET request and return status and body.

Query Parameters:
| Parameter | Type   | Required | Description  |
| --------- | ------ | -------- | ------------ |
| url       | string | Yes      | URL to fetch |

Response:
```json
{
  "url": "https://google.com",
  "status_code": 200,
  "body": "<!doctype html>...."
}
```
Error Response Example:
```json
{"error":"Get \"https://nonexistent.site\": dial tcp: lookup nonexistent.site: no such host"}
```
Curl Example:
```sh
curl -s "http://localhost:8080/api/v1/curl?url=https://google.com" | jq
```
### 5. TCP Port Check
```bash
GET /api/v1/tcp?host=<host>&port=<port>
```

**Description**: Check if a TCP port is open on a given host.

Query Parameters:
| Parameter | Type   | Required | Description    |
| --------- | ------ | -------- | -------------- |
| host      | string | Yes      | Hostname or IP |
| port      | string | Yes      | Port number    |

Response:
```json
{
  "host": "google.com",
  "port": "80",
  "open": true
}
```
Error Response Example:
```json
{
  "host": "google.com",
  "port": "801",
  "open": false,
  "error": "dial tcp 142.251.98.113:801: i/o timeout"
}
```
Curl Example:
```json
curl -s "http://localhost:8080/api/v1/tcp?host=google.com&port=80" | jq
```
### 6. Host Information
```bash
GET /api/v1/info
```
**Description**: Returns general system and runtime information about the host where the service is running.

Response:
```json
{
  "hostname": "host-name",
  "ips": [
    "192.168.1.1",
    "192.168.12.12",
    "192.168.13.13"
  ],
  "macs": [
    "62:55:44:22:g8:10",
    "72:35:34:04:d4:10"
  ],
  "os": "darwin",
  "arch": "amd64",
  "goVersion": "go1.24.12",
  "startTime": "2026-01-28T21:31:12+02:00",
  "now": "2026-01-28T21:36:14+02:00",
  "env": {
    "COLORTERM": "truecolor",
    "COMMAND_MODE": "unix2003"
  },
  "cloud": {
    "provider": "local",
    "region": "",
    "zone": "",
    "instance": "",
    "extra": null
  },
  "kubernetes": {
    "enabled": false,
    "podName": "",
    "podNamespace": "",
    "podIP": "",
    "nodeName": "",
    "serviceAccount": "",
    "container": ""
  }
}
```
Error Response Example:
```json
{ "error": "failed to collect host info" }
```
Curl Example:
```bash
curl -s "http://localhost:8080/api/v1/info" | jq
```
### 7. Cloud Provider Detection
```bash
GET /api/v1/cloud
```
**Description**: Detects whether the service is running in a supported cloud provider and returns cloud metadata.

Response:
```json
{
  "provider": "local",
  "region": "",
  "zone": "",
  "instance": "",
  "extra": null
}
```
Response (Not in Cloud):
```json
{ "provider": "unknown" }
```
Error Response Example:
```json
{ "error": "cloud metadata service unavailable" }
```
Curl Example:
```bash
curl -s "http://localhost:8080/api/v1/cloud" | jq
```
### 8. Kubernetes Environment Info
```bash
GET /api/v1/kubernetes
```
**Description**: Detects whether the application is running inside a Kubernetes cluster and returns pod-level metadata.

Response:
```json
{
  "enabled": false,
  "podName": "",
  "podNamespace": "",
  "podIP": "",
  "nodeName": "",
  "serviceAccount": "",
  "container": ""
}
```
Error Response Example:
```json
{ "error": "kubernetes API not accessible" }
```
Curl Example:
```bash
curl -s "http://localhost:8080/api/v1/kubernetes" | jq
```

## Notes / Best Practices

- All endpoints return JSON. If a parameter is missing or invalid, the response contains an error field.

- Use `jq` for pretty-printing JSON responses:
  ```sh
  curl -s "http://localhost:8080/api/v1/ping?host=google.com" | jq
  ```

- Endpoints are designed for network diagnostics and monitoring. Avoid running intensive requests in high frequency on external hosts.
