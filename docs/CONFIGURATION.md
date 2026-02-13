<!-- omit in toc -->
# Configuration Guide


<!-- omit in toc -->
## Table of Contents
- [Environment Variables](#environment-variables)
  - [Server Configuration](#server-configuration)
  - [Example Usage](#example-usage)
- [CORS Configuration](#cors-configuration)
  - [Default Behavior](#default-behavior)
  - [Custom Origins](#custom-origins)

## Environment Variables

All configuration options can be set via environment variables.

### Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `HOSTINFO_PORT` | `8080` | Port to listen on |
| `HOSTINFO_HOST` | `0.0.0.0` | Host to bind to |
| `HOSTINFO_DEBUG` | `false` | Basic Debug param |
| `HOSTINFO_FRONTEND_PATH` | `./frontend/dist` | Frontend build files |
| `FF_ENVIRONMENT_VARIABLES` | `false` | return env vars if set to "true" |


### Example Usage

```bash
# Standard environment variables (recommended)
HOSTINFO_PORT=8080 \
HOSTINFO_HOST=localhost \
HOSTINFO_DEBUG=false \
./hostinfo
```

## CORS Configuration

Configure Cross-Origin Resource Sharing for frontend access.

### Default Behavior

By default, only localhost origins are allowed:
- `http://localhost:8080`
- `http://127.0.0.1:8080`

### Custom Origins

```bash
# Single origin
export CORS_ALLOWED_ORIGINS="https://hostinfo.example.com"

# Multiple origins (comma-separated)
export CORS_ALLOWED_ORIGINS="https://hostinfo.example.com,https://admin.example.com"

./hostinfo
```

---

For more information on security configuration, see [SECURITY.md](SECURITY.md).
For deployment instructions, see [DEPLOYMENT.md](DEPLOYMENT.md).