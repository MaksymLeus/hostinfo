package host

type HostInfo struct {
	// Identity
	Hostname string `json:"hostname"`

	// Network
	IPs  []string `json:"ips"`
	MACs []string `json:"macs"`

	// System
	OS        string `json:"os"`
	Distro    string `json:"distro,omitempty"`
	Arch      string `json:"arch"`
	GoVersion string `json:"goVersion"`

	// Time
	StartTime string `json:"startTime"`
	UpdatedAt string `json:"updatedAt"`
	Uptime    string `json:"uptime"`

	// Runtime
	Env map[string]string `json:"env,omitempty"`

	// Hardware
	CPU    CPUInfo    `json:"cpu"`
	Memory MemoryInfo `json:"memory"`
	Load   LoadInfo   `json:"load"`

	// Platform
	Cloud      CloudInfo      `json:"cloud,omitempty"`
	Kubernetes KubernetesInfo `json:"kubernetes,omitempty"`
}

type CloudInfo struct {
	Provider string            `json:"provider"`
	Region   string            `json:"region"`
	Zone     string            `json:"zone"`
	Instance string            `json:"instance"`
	Extra    map[string]string `json:"extra"`
}

type KubernetesInfo struct {
	Enabled        bool   `json:"enabled"`
	PodName        string `json:"podName"`
	PodNamespace   string `json:"podNamespace"`
	PodIP          string `json:"podIP"`
	NodeName       string `json:"nodeName"`
	ServiceAccount string `json:"serviceAccount"`
	Container      string `json:"container"`
}

type CPUInfo struct {
	Cores        int     `json:"cores"`
	Model        string  `json:"model"`
	UsagePercent float64 `json:"usagePercent"`
}

type MemoryInfo struct {
	TotalMB uint64  `json:"totalMB"`
	UsedMB  uint64  `json:"usedMB"`
	UsedPct float64 `json:"usedPercent"`
}

type LoadInfo struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}
