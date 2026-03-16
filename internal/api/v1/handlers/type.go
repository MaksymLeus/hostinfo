package handlers

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

	// Environment
	Env map[string]string `json:"env,omitempty"`

	// Hardware
	CPU    CPUInfo    `json:"cpu"`
	Memory MemoryInfo `json:"memory"`
	Load   LoadInfo   `json:"load"`

	// Platform
	Cloud      CloudInfo      `json:"cloud,omitempty"`
	Runtime    RuntimeInfo    `json:"runtime,omitempty"`
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

// RuntimeInfo represents the detected runtime environment
type RuntimeInfo struct {
	Environment string `json:"environment"` // "kubernetes", "docker", "bare-metal"
	Details     string `json:"details"`     // Extra info, like container ID
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

// -------------------- Commands -------------------- \\
type PingResult struct {
	Host        string  `json:"host"`
	PacketsSent int     `json:"packets_sent"`
	PacketsRecv int     `json:"packets_recv"`
	Loss        float64 `json:"loss_percent"`
	MinRTT      string  `json:"min_rtt"`
	MaxRTT      string  `json:"max_rtt"`
	AvgRTT      string  `json:"avg_rtt"`
}

type DNSResult struct {
	Host  string   `json:"host"`
	IPs   []string `json:"ips"`
	CNAME string   `json:"cname"`
}

type CurlResult struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}

type TCPResult struct {
	Host  string `json:"host"`
	Port  string `json:"port"`
	Open  bool   `json:"open"`
	Error string `json:"error,omitempty"`
}
