# Add Kubernetes
### ☸️ Kubernetes (NO RBAC)
Uses **Downward API only**
- Pod name
- Namespace
- Pod IP
- Node name
- Service account
- Container name

If you later deploy to K8s, probes map directly:
```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /healthz
    port: 8080
  initialDelaySeconds: 3
  periodSeconds: 5
```

# Upgrade UI
- Make new ui on js
- favicon.ico 


# CI
- Docker cache
- speed up CI

Scenario B — .releaserc stored in another repo

Not supported out of the box, because semantic-release always evaluates config against the current working directory (current repo codebase).

BUT you can make it work by:

Option 1 — Checkout config repo first

In workflow:
```yaml
- uses: actions/checkout@v4
  with:
    path: current

- uses: actions/checkout@v4
  with:
    repository: myorg/semantic-config
    path: semantic-config

- name: Use external .releaserc
  run: cp semantic-config/.releaserc current/.releaserc
```

This approach is used in mono-repos and org-wide standards.