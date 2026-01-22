# Deployment

This document describes available deployment strategies for the **hostinfo** web application, including bare-metal, Docker, Docker Compose, systemd, and cloud environments.

---

## 1. Deployment Models

hostinfo supports:

- **Bare-metal binary** (Linux/Windows/macOS)
- **Docker container**
- **Docker Compose**
- **systemd** service
- **Cloud VMs** (AWS EC2, GCP Compute, Azure VM)
- (Optional) **Kubernetes** via container image

hostinfo is stateless:
- No external dependencies
- No persistent storage required
- No configuration files required by default

---

## 2. Deployment Preparation

### Clone and build:

```bash
git clone https://github.com/MaksymLeus/hostinfo.git
cd hostinfo
./build.sh all   # or: go build -o hostinfo ./cmd/server
```

Produced binaries are typically placed in `bin/`.

---

## 3. Bare-Metal Deployment (Linux)

### 3.1 Upload binary to server

```bash
scp bin/hostinfo-linux-x64 user@server:/opt/hostinfo/hostinfo
```

### 3.2 Start manually

```bash
ssh user@server
cd /opt/hostinfo
chmod +x hostinfo
./hostinfo
```

Service runs on:

```
http://<server-ip>:8080
```

---

## 4. systemd Service (Linux)

### 4.1 Create unit file

`/etc/systemd/system/hostinfo.service`:

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

[Install]
WantedBy=multi-user.target
```

### 4.2 Enable & start

```bash
sudo systemctl enable hostinfo
sudo systemctl start hostinfo
```

### 4.3 Check status and logs

```bash
sudo systemctl status hostinfo
sudo journalctl -u hostinfo -f
```

---

## 5. Docker Deployment

### 5.1 Build image locally

```bash
docker build -t hostinfo .
```

### 5.2 Run container

```bash
docker run -p 8080:8080 hostinfo
```

Access:

```
http://localhost:8080
```

---

## 6. Docker Compose Deployment

`docker-compose.yml` (existing in project):

```bash
docker compose build
docker compose up -d
```

Manage lifecycle:

```bash
docker compose logs -f
docker compose down
```

---

## 7. Cloud VM Deployment

### Supported environments (no extra config):

- AWS EC2
- GCP Compute Engine
- Azure VM

hostinfo auto-detects cloud metadata **without** credentials.

#### Example Ubuntu VM deployment:

```bash
# build binary locally
./build.sh all

# copy to VM
scp bin/hostinfo-linux-x64 ubuntu@<vm-ip>:/opt/hostinfo/hostinfo

# run on VM
ssh ubuntu@<vm-ip>
cd /opt/hostinfo
./hostinfo
```

### Firewall considerations:

- AWS EC2 → Security Group: allow inbound TCP 8080
- GCP Compute → Firewall Rule: allow TCP 8080
- Azure NSG → allow TCP 8080

---

## 8. Reverse Proxy (Optional)

For production HTTPS, recommend:

- **Nginx**
- **Caddy**
- **Traefik**

Example nginx site config:

```
location / {
    proxy_pass http://localhost:8080;
}
```

Enable TLS via Let's Encrypt or Traefik.

---

## 9. Kubernetes (Optional)

If deployed to Kubernetes, use the existing Docker image:

Example minimal manifest:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hostinfo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hostinfo
  template:
    metadata:
      labels:
        app: hostinfo
    spec:
      containers:
      - name: hostinfo
        image: maximleus/hostinfo:latest
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: hostinfo
spec:
  type: ClusterIP
  selector:
    app: hostinfo
  ports:
  - port: 80
    targetPort: 8080
```

Ingress optional depending on cluster.

---

## 10. Health & Monitoring

hostinfo exposes a UI at:

```
/       (dashboard)
```

If you later add JSON endpoints (e.g., `/api/info`), update this section.

---

## 11. Updates & Rollbacks

### Bare-Metal:

```bash
systemctl stop hostinfo
# replace binary
systemctl start hostinfo
```

### Docker:

```bash
docker pull maximleus/hostinfo:latest
docker compose up -d
```

### Kubernetes:

```bash
kubectl rollout restart deployment/hostinfo
```

Rollbacks via ReplicaSet history.

---

## 12. Deployment Recommendations

For production:

- Run behind reverse proxy with HTTPS
- Use systemd for bare-metal
- Use Docker or Kubernetes for cloud
- Prefer immutable container deployments
- Do not expose port 8080 directly to public internet without TLS

---

## License

MIT — see `LICENSE.md`.
