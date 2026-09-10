package metrics

import (
	"hostinfo/internal/api/v1/handlers"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type HostInfoCollector struct {
	uptimeMetric   *prometheus.Desc
	cpuUsageMetric *prometheus.Desc
	memoryUsage    *prometheus.Desc
	diskUsage      *prometheus.Desc
}

var startTime = time.Now()

func NewHostInfoCollector() *HostInfoCollector {
	return &HostInfoCollector{
		uptimeMetric: prometheus.NewDesc(
			"hostinfo_uptime_seconds",
			"Duration since the hostinfo service started",
			nil, nil,
		),
		cpuUsageMetric: prometheus.NewDesc(
			"hostinfo_cpu_usage_percent",
			"Current CPU usage percentage",
			nil, nil,
		),
		memoryUsage: prometheus.NewDesc(
			"hostinfo_memory_used_bytes",
			"Memory used in bytes",
			nil, nil,
		),
		diskUsage: prometheus.NewDesc(
			"hostinfo_disk_used_percent",
			"Disk usage percentage by mount path",
			[]string{"path", "fstype"}, nil,
		),
	}
}

func (c *HostInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.uptimeMetric
	ch <- c.cpuUsageMetric
	ch <- c.memoryUsage
	ch <- c.diskUsage
}

func (c *HostInfoCollector) Collect(ch chan<- prometheus.Metric) {
	// 1. Uptime
	ch <- prometheus.MustNewConstMetric(
		c.uptimeMetric,
		prometheus.GaugeValue,
		time.Since(startTime).Seconds(),
	)

	// 2. CPU Usage
	cpuInfo := handlers.GetCPUForMetrics() // We'll add this helper below, or calculate it
	ch <- prometheus.MustNewConstMetric(
		c.cpuUsageMetric,
		prometheus.GaugeValue,
		cpuInfo.UsagePercent,
	)

	// 3. Memory Usage
	memInfo := handlers.GetMemoryForMetrics()
	ch <- prometheus.MustNewConstMetric(
		c.memoryUsage,
		prometheus.GaugeValue,
		float64(memInfo.UsedMB)*1024*1024,
	)

	// 4. Disk Usage
	disks := handlers.GetDiskForMetrics()
	for _, d := range disks {
		ch <- prometheus.MustNewConstMetric(
			c.diskUsage,
			prometheus.GaugeValue,
			d.UsedPercent,
			d.Path,
			d.FSType,
		)
	}
}
