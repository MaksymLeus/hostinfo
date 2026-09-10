package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	gopsnet "github.com/shirou/gopsutil/v3/net"
)

type HostInfoCollector struct {
	uptimeMetric   *prometheus.Desc
	cpuUsageMetric *prometheus.Desc
	memoryUsage    *prometheus.Desc
	diskUsage      *prometheus.Desc
	netBytesSent   *prometheus.Desc
	netBytesRecv   *prometheus.Desc
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
		netBytesSent: prometheus.NewDesc(
			"hostinfo_net_bytes_sent_total",
			"Total network bytes sent by interface",
			[]string{"interface"}, nil,
		),
		netBytesRecv: prometheus.NewDesc(
			"hostinfo_net_bytes_recv_total",
			"Total network bytes received by interface",
			[]string{"interface"}, nil,
		),
	}
}

func (c *HostInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.uptimeMetric
	ch <- c.cpuUsageMetric
	ch <- c.memoryUsage
	ch <- c.diskUsage
	ch <- c.netBytesSent
	ch <- c.netBytesRecv
}

func (c *HostInfoCollector) Collect(ch chan<- prometheus.Metric) {
	// 1. Uptime
	ch <- prometheus.MustNewConstMetric(
		c.uptimeMetric,
		prometheus.GaugeValue,
		time.Since(startTime).Seconds(),
	)

	// 2. CPU Usage
	cpuUsage := 0.0
	if usage, err := cpu.Percent(0, false); err == nil && len(usage) > 0 {
		cpuUsage = usage[0]
	}
	ch <- prometheus.MustNewConstMetric(
		c.cpuUsageMetric,
		prometheus.GaugeValue,
		cpuUsage,
	)

	// 3. Memory Usage
	if vm, err := mem.VirtualMemory(); err == nil && vm != nil {
		ch <- prometheus.MustNewConstMetric(
			c.memoryUsage,
			prometheus.GaugeValue,
			float64(vm.Used),
		)
	}

	// 4. Disk Usage
	if partitions, err := disk.Partitions(false); err == nil {
		for _, p := range partitions {
			if isIgnoredFilesystem(p.Fstype) {
				continue
			}
			if usage, err := disk.Usage(p.Mountpoint); err == nil {
				ch <- prometheus.MustNewConstMetric(
					c.diskUsage,
					prometheus.GaugeValue,
					usage.UsedPercent,
					p.Mountpoint,
					p.Fstype,
				)
			}
		}
	}

	// 5. Network I/O
	if ioCounters, err := gopsnet.IOCounters(true); err == nil {
		for _, io := range ioCounters {
			ch <- prometheus.MustNewConstMetric(
				c.netBytesSent,
				prometheus.CounterValue,
				float64(io.BytesSent),
				io.Name,
			)
			ch <- prometheus.MustNewConstMetric(
				c.netBytesRecv,
				prometheus.CounterValue,
				float64(io.BytesRecv),
				io.Name,
			)
		}
	}
}

func isIgnoredFilesystem(fstype string) bool {
	ignored := map[string]bool{
		"squashfs": true,
		"tmpfs":    true,
		"devtmpfs": true,
		"overlay":  true,
		"aufs":     true,
	}
	return ignored[fstype]
}
