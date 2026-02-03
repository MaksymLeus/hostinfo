package host

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

var startTime = time.Now()

func Collect() HostInfo {
	now := time.Now()

	return HostInfo{
		// Identity
		Hostname: getHostname(),

		// Network
		IPs:  getIPs(),
		MACs: getMACs(),

		// System
		OS:        runtime.GOOS,
		Distro:    getDistro(),
		Arch:      runtime.GOARCH,
		GoVersion: runtime.Version(),

		// Time
		StartTime: startTime.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
		Uptime:    getUptime(startTime),

		// Runtime
		Env: getEnv(),

		// Hardware
		CPU:    getCPU(),
		Memory: getMemory(),
		Load:   getLoad(),

		// Platform
		Cloud:      DetectCloud(),
		Kubernetes: DetectKubernetes(),
	}
}

func getDistro() string {
	if runtime.GOOS != "linux" {
		return ""
	}

	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(line[13:], `"`)
		}
	}
	return ""
}

func getUptime(start time.Time) string {
	d := time.Since(start)

	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	return fmt.Sprintf("%dh %dm %ds", h, m, s)
}

func getHostname() string {
	h, _ := os.Hostname()
	return h
}

func getIPs() []string {
	var ips []string
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}

func getMACs() []string {
	var macs []string
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		if i.HardwareAddr != nil {
			macs = append(macs, i.HardwareAddr.String())
		}
	}
	return macs
}

func getEnv() map[string]string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		kv := []rune(e)
		for i, c := range kv {
			if c == '=' {
				env[string(kv[:i])] = string(kv[i+1:])
				break
			}
		}
	}
	return env
}

func DetectCloud() CloudInfo {
	if aws := detectAWS(); aws.Provider != "" {
		return aws
	}
	if gcp := detectGCP(); gcp.Provider != "" {
		return gcp
	}
	if azure := detectAzure(); azure.Provider != "" {
		return azure
	}
	if dockerInfo := detectDocker(); dockerInfo.Provider != "" {
		return dockerInfo
	}
	return CloudInfo{Provider: "local"}
}

func detectAWS() CloudInfo {
	client := http.Client{Timeout: 500 * time.Millisecond}

	req, _ := http.NewRequest("GET",
		"http://169.254.169.254/latest/meta-data/instance-id", nil)

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return CloudInfo{}
	}
	defer resp.Body.Close()

	id, _ := io.ReadAll(resp.Body)

	region := awsMeta("placement/region")
	zone := awsMeta("placement/availability-zone")

	return CloudInfo{
		Provider: "aws",
		Region:   region,
		Zone:     zone,
		Instance: string(id),
		Extra: map[string]string{
			"AMI":  awsMeta("ami-id"),
			"Type": awsMeta("instance-type"),
		},
	}
}

func awsMeta(path string) string {
	client := http.Client{Timeout: 300 * time.Millisecond}
	resp, err := client.Get("http://169.254.169.254/latest/meta-data/" + path)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b)
}

func detectGCP() CloudInfo {
	client := http.Client{Timeout: 500 * time.Millisecond}

	req, _ := http.NewRequest("GET",
		"http://metadata.google.internal/computeMetadata/v1/instance/id", nil)
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return CloudInfo{}
	}
	defer resp.Body.Close()

	id, _ := io.ReadAll(resp.Body)

	return CloudInfo{
		Provider: "gcp",
		Region:   gcpMeta("instance/region"),
		Zone:     gcpMeta("instance/zone"),
		Instance: string(id),
		Extra: map[string]string{
			"Machine": gcpMeta("instance/machine-type"),
			"Project": gcpMeta("project/project-id"),
		},
	}
}

func gcpMeta(path string) string {
	client := http.Client{Timeout: 300 * time.Millisecond}
	req, _ := http.NewRequest("GET",
		"http://metadata.google.internal/computeMetadata/v1/"+path, nil)
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b)
}

func detectAzure() CloudInfo {
	client := http.Client{Timeout: 500 * time.Millisecond}

	req, _ := http.NewRequest("GET",
		"http://169.254.169.254/metadata/instance?api-version=2021-02-01", nil)
	req.Header.Set("Metadata", "true")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return CloudInfo{}
	}
	defer resp.Body.Close()

	return CloudInfo{
		Provider: "azure",
		Extra: map[string]string{
			"VM": "Azure VM detected",
		},
	}
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

func DetectKubernetes() KubernetesInfo {
	// Always present in k8s
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		return KubernetesInfo{Enabled: false}
	}

	return KubernetesInfo{
		Enabled:        false, //true, set to false to avoid showing k8s info in the UI until all is verified
		PodName:        "POD_NAME",
		PodNamespace:   "POD_NAMESPACE",
		PodIP:          "POD_IP",
		NodeName:       "NODE_NAME",
		ServiceAccount: "SERVICE_ACCOUNT",
		Container:      "CONTAINER_NAME",
	}
}

// -------- Hardware -------- //
func getCPU() CPUInfo {
	cores, _ := cpu.Counts(true)
	info, _ := cpu.Info()
	usage, _ := cpu.Percent(0, false)

	model := ""
	if len(info) > 0 {
		model = info[0].ModelName
	}

	return CPUInfo{
		Cores:        cores,
		Model:        model,
		UsagePercent: usage[0],
	}
}

func getMemory() MemoryInfo {
	vm, _ := mem.VirtualMemory()

	return MemoryInfo{
		TotalMB: vm.Total / 1024 / 1024,
		UsedMB:  vm.Used / 1024 / 1024,
		UsedPct: vm.UsedPercent,
	}
}

func getLoad() LoadInfo {
	avg, err := load.Avg()
	if err != nil {
		return LoadInfo{}
	}

	return LoadInfo{
		Load1:  avg.Load1,
		Load5:  avg.Load5,
		Load15: avg.Load15,
	}
}
