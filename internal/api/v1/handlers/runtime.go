package handlers

import (
	"hostinfo/internal/custom"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func Kubernetes(c echo.Context) error {
	info := DetectKubernetes()
	return c.JSON(http.StatusOK, info)
}

func detectDocker() CloudInfo {
	// Check /.dockerenv file
	if _, err := os.Stat("/.dockerenv"); err != nil {
		return CloudInfo{} // Not in Docker
	}

	// Try to read container ID from /proc/self/cgroup
	containerID := ""
	if data, err := os.ReadFile("/proc/self/cgroup"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			parts := strings.Split(line, "/")
			last := parts[len(parts)-1]
			if len(last) >= 12 {
				containerID = last
				break
			}
		}
	}

	return CloudInfo{
		Provider: "docker",
		Instance: containerID,
	}
}

// DetectRuntime detects if we are running in Kubernetes, Docker, or bare-metal
func DetectRuntime() RuntimeInfo {
	// 1️⃣ Kubernetes detection
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" &&
		custom.FileExists("/var/run/secrets/kubernetes.io/serviceaccount/token") {
		return RuntimeInfo{
			Environment: "kubernetes",
			Details:     "Kubernetes pod detected",
		}
	}

	// 2️⃣ Docker / container detection
	if custom.FileExists("/.dockerenv") {
		containerID := detectContainerID()
		return RuntimeInfo{
			Environment: "docker",
			Details:     containerID,
		}
	}

	// 3️⃣ Bare-metal / plain Go
	return RuntimeInfo{
		Environment: "bare-metal",
		Details:     "Running as standalone Go process",
	}
}

// Helper to detect container ID from /proc/self/cgroup
func detectContainerID() string {
	data, err := os.ReadFile("/proc/self/cgroup")
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "/")
		last := parts[len(parts)-1]
		if len(last) >= 12 {
			return last
		}
	}
	return ""
}

func DetectKubernetes() KubernetesInfo {
	// Hard guarantees from Kubernetes
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		return KubernetesInfo{Enabled: false}
	}

	// ServiceAccount token is mounted in all pods (unless explicitly disabled)
	if !custom.FileExists("/var/run/secrets/kubernetes.io/serviceaccount/token") {
		return KubernetesInfo{Enabled: false}
	}

	return KubernetesInfo{
		Enabled:        true,
		PodName:        os.Getenv("POD_NAME"),
		PodNamespace:   os.Getenv("POD_NAMESPACE"),
		PodIP:          os.Getenv("POD_IP"),
		NodeName:       os.Getenv("NODE_NAME"),
		ServiceAccount: os.Getenv("SERVICE_ACCOUNT"),
		Container:      os.Getenv("CONTAINER_NAME"),
	}
}
