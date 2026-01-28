package host

type HostInfo struct {
	Hostname   string            `json:"hostname"`
	IPs        []string          `json:"ips"`
	MACs       []string          `json:"macs"`
	OS         string            `json:"os"`
	Arch       string            `json:"arch"`
	GoVersion  string            `json:"goVersion"`
	StartTime  string            `json:"startTime"`
	Now        string            `json:"now"`
	Env        map[string]string `json:"env"`
	Cloud      CloudInfo         `json:"cloud"`
	Kubernetes KubernetesInfo    `json:"kubernetes"`
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
