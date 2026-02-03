package api

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
