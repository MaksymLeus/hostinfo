# 06 — Cloud Detection

hostinfo automatically detects cloud environments (AWS, GCP, Azure) and gathers metadata **without requiring credentials**.  
This document explains the detection logic, endpoints, timeouts, and fallback behavior.

---

## 1. Overview

Cloud detection is:

- **Safe**: read-only, short timeouts
- **Non-blocking**: asynchronous or concurrent calls
- **Fallback-aware**: defaults to local host if no cloud detected
- **Environment-agnostic**: works in containers, bare-metal, and VMs

Detection order:

```
1. Container detection
2. Cloud detection
3. Local fallback
```

---

## 2. Container Detection (Precedes Cloud)

hostinfo first determines if it is running inside a container:

- **Docker / Podman / LXC** detection:
  - Check `/proc/self/cgroup` for container patterns
  - Check `/proc/1/cgroup` for docker identifiers
  - Inspect hostname for container hashes

If container detected:

- Adds container uptime
- Marks environment as `container`
- Continues to cloud detection if available

---

## 3. AWS Detection

AWS metadata is fetched from **Instance Metadata Service (IMDS)**:

- Endpoint: `http://169.254.169.254/latest/meta-data/`
- Information retrieved:
  - Instance ID
  - Region & Availability Zone
  - Instance type
  - Private IP
- Requests are **short timeout** (~500ms)
- Failures do **not crash** the application

Example metadata keys:

```
instance-id
instance-type
placement/availability-zone
local-ipv4
```

---

## 4. Google Cloud Detection

GCP metadata service:

- Endpoint: `http://metadata.google.internal/computeMetadata/v1/`
- Requires header: `Metadata-Flavor: Google`
- Retrieves:
  - Project ID
  - Zone
  - Machine type
- Timeout is short (~500ms)
- Failure → fallback

Example:

```bash
curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/id
```

---

## 5. Azure Detection

Azure Instance Metadata Service (IMDS):

- Endpoint: `http://169.254.169.254/metadata/instance?api-version=2021-02-01`
- Header: `Metadata: true`
- Retrieves:
  - VM ID
  - Location
  - VM size/type
- Timeout: short (~500ms)
- Failures do **not crash** app

Example:

```bash
curl -H "Metadata:true" "http://169.254.169.254/metadata/instance?api-version=2021-02-01"
```

---

## 6. Local / Fallback

If none of the cloud endpoints respond:

- Environment marked as `local`
- Only host and container info is shown
- Ensures dashboard always displays minimal info

---

## 7. Metadata Caching

- Metadata is cached in memory per process
- Avoids repeated network calls on each HTTP request
- Cache invalidated on restart
- No persistence to disk

---

## 8. Parallel Probing Logic

Detection sequence (simplified):

```
start -> container check
      -> launch cloud probes concurrently (AWS, GCP, Azure)
      -> wait max 500ms per probe
      -> first successful response sets cloud type
      -> others cancel
      -> fallback to local if all fail
```

This prevents blocking the HTTP server.

---

## 9. Security Considerations

- No credentials required
- Only public metadata endpoints are accessed
- Cloud detection is **read-only**
- Sensitive values (API keys, secrets) are **never accessed**
- Failures are silent and logged optionally

---

## 10. Deployment Notes

- Works inside containers → detects container + cloud host if present
- Works in ephemeral CI/CD runners → falls back gracefully
- Works in multi-cloud hybrid setups → first responding cloud is chosen
- Safe for production → no network egress besides metadata endpoints

---

## 11. Troubleshooting

| Issue | Cause | Resolution |
|---|---|---|
| Cloud metadata not showing | Not running on cloud | Expected behavior, fallback to local |
| Dashboard slow | Cloud metadata timed out | Network latency; consider increasing timeout if desired |
| Metadata endpoint blocked | Firewall or VPC rules | Ensure access to IMDS (AWS/Azure) or metadata server (GCP) |
| Running inside container only | Container detected but cloud unavailable | Expected; container may not have cloud metadata access |

---

## License

MIT — see `LICENSE.md`.
