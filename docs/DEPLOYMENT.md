<!-- omit in toc -->
# Deployment Guide

This document covers various deployment options for the **hostinfo**.

<!-- omit in toc -->
## Table of Contents
- [Docker (Recommended)](#docker-recommended)
  - [Image Details](#image-details)
  - [Quick Start](#quick-start)
  - [Clening Up](#clening-up)
  - [Environment Variables](#environment-variables)
- [Docker Compose](#docker-compose)
  - [Standalone docker-compose.yml](#standalone-docker-composeyml)
  - [Rund](#rund)
- [Kubernetes](#kubernetes)
  - [Helm Chart](#helm-chart)
    - [Quick Start](#quick-start-1)
    - [Check Deployment Status](#check-deployment-status)
    - [Generate kubernetes configs](#generate-kubernetes-configs)
  - [Manual](#manual)
    - [Quick Start](#quick-start-2)
- [Linux Server (Binary)](#linux-server-binary)
  - [Build and Deploy](#build-and-deploy)
  - [Permissions](#permissions)
- [systemd Service](#systemd-service)
  - [Enable and Start](#enable-and-start)
- [Production Checklist](#production-checklist)
  - [Before Deployment](#before-deployment)
  - [Firewall considerations:](#firewall-considerations)
- [❗ Troubleshooting](#-troubleshooting)
  - [Common Issues](#common-issues)
  - [Kubernetes](#kubernetes-1)

## Docker (Recommended)

Hostinfo is available as a Docker image for easy deployment.

### Image Details

| Property | Value |
|----------|-------|
| Registry | Docker Hub ([`maximleus/hostinfo`](https://hub.docker.com/r/maximleus/hostinfo)) |
| Base Image | `alpine:3.20` |
| Platforms | `linux/amd64`, `linux/arm64` |
| Size | ~47 compressed |
| Tags | `latest`, version tags (e.g., `v0.2.3`) |

### Quick Start

```bash
docker run -d \
  --name hostinfo \
  -p 8080:8080 \
  -e HOSTINFO_HOST=0.0.0.0 \
  -e HOSTINFO_PORT=8080 \
  maximleus/hostinfo:latest
```
Open your browser and visit: **http://localhost:8080**

### Clening Up
```bash
docker ps -a
docker rm <container_id>
docker rmi hostinfo:latest
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HOSTINFO_PORT` | `8080` | Port to listen on |
| `HOSTINFO_HOST` | `0.0.0.0` | Host to bind to |
| `HOSTINFO_DEBUG` | `false` | Basic Debug param |
| `HOSTINFO_FRONTEND_PATH` | `./frontend/dist` | Frontend build files |
| `CORS_ALLOWED_ORIGINS` | (localhost) | Allowed CORS origins |

## Docker Compose 

### Standalone docker-compose.yml

Create a `docker-compose.yml` file anywhere on your system or use existing in project [`docker-compose.yml`](../docker-compose.yml) :

```yaml
services:
  hostinfo:
    image: maximleus/hostinfo:latest
    container_name: hostinfo
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - HOSTINFO_HOST: "0.0.0.0"
      - HOSTINFO_PORT: 8080
```
### Rund
```bash
# Start the container
docker compose up -d

# View logs
docker compose logs -f

# Stop
docker compose down
```

## Kubernetes

### Helm Chart

<!-- #### Add Helm Repository (TODO: create helm repo and release it [prefe automaticly] )
helm repo add <repo-name> <repo-url>
helm repo update -->

#### Quick Start 

```bash
# Install + Upgrade
helm upgrade --install hostinfo ./helm \
  --namespace <namespace> --create-namespace

# Uninstall Release
helm uninstall hostinfo -n <namespace>

# To override default values:
helm upgrade --install hostinfo ./helm \
  --namespace <namespace> --create-namespace -f values-dev.yaml

# Or override inline:
helm upgrade --install hostinfo ./helm \
  --namespace <namespace> --create-namespace \
  --set image.tag=1.2.0
```

#### Check Deployment Status
```bash
# Status
helm list -n <namespace>
kubectl get pods -n <namespace>
kubectl get svc -n <namespace>

# Check rollout status:
kubectl rollout status deployment/hostinfo -n <namespace>
```

#### Generate kubernetes configs
If you whant run using `kubectl` you can generate template using helm
```bash
helm template hostinfo ./helm 
```

### Manual 
Create a `hostinfo.yaml` file anywhere on your system

Paste your full manifest into it (Service + Deployment):
```yaml
# Source: hostinfo/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: hostinfo
  labels:
    helm.sh/chart: hostinfo-1.0.0
    app.kubernetes.io/name: hostinfo
    app.kubernetes.io/instance: hostinfo
    app.kubernetes.io/version: "1.0.0"
    app.kubernetes.io/managed-by: Helm
spec:
  type: ClusterIP
  ports:
    - port: 8080
      targetPort: 8080
      protocol: TCP
      name: http
  selector:
      app.kubernetes.io/name: hostinfo
      app.kubernetes.io/instance: hostinfo
---
# Source: hostinfo/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hostinfo
  labels:
    helm.sh/chart: hostinfo-1.0.0
    app.kubernetes.io/name: hostinfo
    app.kubernetes.io/instance: hostinfo
    app.kubernetes.io/version: "1.0.0"
    app.kubernetes.io/managed-by: Helm
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: hostinfo
      app.kubernetes.io/instance: hostinfo
  template:
    metadata:
      labels:
        app.kubernetes.io/name: hostinfo
        app.kubernetes.io/instance: hostinfo
    spec:  
      containers:
        - name: hostinfo
          image: "maximleus/hostinfo:latest"
          imagePullPolicy: Always
          env:
            - name: HOSTINFO_DEBUG
              value: "false"
            - name: HOSTINFO_HOST
              value: "0.0.0.0"
            - name: HOSTINFO_PORT
              value: "8080"
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace

          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
          readinessProbe:
            httpGet:
              path: /healthz/ready # The path to the readiness endpoint
              port: 8080   # The port the application is listening on
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /healthz/live # The path to the readiness endpoint
              port: 8080   # The port the application is listening on
            initialDelaySeconds: 15
            periodSeconds: 20
          resources:
            null
```

#### Quick Start

```bash
# Create namespace 
kubectl create namespace <namespace>

# Deploy to specific namespace:
# For update  just edit hostinfo.yaml (for example change image tag) and re-apply:
kubectl apply -f hostinfo.yaml -n <namespace>

# Remove Manually Installed Resources
kubectl delete -f hostinfo.yaml  -n <namespace>

```

## Linux Server (Binary)

### Build and Deploy
```bash
# build binary locally
./build.sh os linux
# Produced binaries are typically placed in `bin/`.

# copy to VM
scp bin/hostinfo-linux-x64 ubuntu@<vm-ip>:/opt/hostinfo/hostinfo

# run on VM
ssh ubuntu@<vm-ip>
cd /opt/hostinfo
./hostinfo
```

### Permissions

```bash
# Create hostinfo user
sudo useradd -r -s /bin/false hostinfo

# Set ownership
sudo chown -R hostinfo:hostinfo /opt/hostinfo
```

## systemd Service

Create `/etc/systemd/system/hostinfo.service`:

```ini
[Unit]
Description=Hostinfo Web Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/hostinfo
ExecStart=/opt/hostinfo/hostinfo
Restart=on-failure
RestartSec=3

# Environment variables (optional)
# Environment=HOSTINFO_PORT=8080
# Environment=HOSTINFO_HOST=0.0.0.0
# Evironment=HOSTINFO_DEBUG=false

[Install]
WantedBy=multi-user.target
```

### Enable and Start

```bash
# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable hostinfo
sudo systemctl start hostinfo

# Check status
sudo systemctl status hostinfo

# View logs
sudo journalctl -u hostinfo -f
```

## Production Checklist

### Before Deployment

- [ ] Configure firewall rules
- [ ] Run behind reverse proxy with HTTPS
- [ ] Use systemd for bare-metal
- [ ] Use Docker or Kubernetes for cloud
- [ ] Prefer immutable container deployments
- [ ] Do not expose port 8080 directly to public internet without TLS

### Firewall considerations:

- AWS EC2 → Security Group: allow inbound TCP `8080`
- GCP Compute → Firewall Rule: allow TCP `8080`
- Azure NSG → allow TCP `8080`

## ❗ Troubleshooting

### Common Issues
| Issue | Cause | Resolution |
|---|---|---|
| `go: no such file or directory` | Wrong working directory | `cd hostinfo` |
| `exec format error` | Wrong GOARCH/GOOS | Rebuild with correct target |
| Docker build slow | No layer cache | Enable BuildKit |
| `permission denied` | Missing execute flag | `chmod +x hostinfo` or `chmod +x build.sh` |
| CrashLoopBackOff | Misconfigured env/secret | Check secret values |
| ImagePullBackOff | Bad image URL / auth | Check registry creds |
| Ingress 404 | Wrong host/path | Match ingress host/path |
| Cloud metadata not showing | Not running on cloud | Expected behavior, fallback to local |
| Dashboard slow | Cloud metadata timed out | Network latency; consider increasing timeout if desired |
| Metadata endpoint blocked | Firewall or VPC rules | Ensure access to IMDS (AWS/Azure) or metadata server (GCP) |
| Running inside container only | Container detected but cloud unavailable | Expected; container may not have cloud metadata access |
### Kubernetes
```bash
# Check Pod Logs
kubectl logs deployment/<name> -n <namespace>

# Check Pod Events
kubectl describe pod <pod> -n <namespace>
```
---

For more configuration options, see [CONFIGURATION.md](CONFIGURATION.md).
