package metrics

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type DataPoint struct {
	Timestamp  string  `json:"timestamp"`
	CPUPercent float64 `json:"cpuPercent"`
	MemPercent float64 `json:"memPercent"`
}

type MetricsHistory struct {
	mu     sync.Mutex
	points []DataPoint
	maxLen int
}

var History = &MetricsHistory{
	points: make([]DataPoint, 0, 60),
	maxLen: 60, // Keep last 60 data points (e.g., 60 minutes if sampled minutely)
}

// StartHistoryCollector runs a background ticker to record telemetry points
func StartHistoryCollector() {
	// Take initial point immediately
	recordPoint()

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			recordPoint()
		}
	}()
}

func recordPoint() {
	cpuUsage := 0.0
	if usage, err := cpu.Percent(0, false); err == nil && len(usage) > 0 {
		cpuUsage = usage[0]
	}

	memUsage := 0.0
	if vm, err := mem.VirtualMemory(); err == nil && vm != nil {
		memUsage = vm.UsedPercent
	}

	History.mu.Lock()
	defer History.mu.Unlock()

	point := DataPoint{
		Timestamp:  time.Now().Format(time.RFC3339),
		CPUPercent: cpuUsage,
		MemPercent: memUsage,
	}

	if len(History.points) >= History.maxLen {
		// Shift slice left (drop oldest)
		copy(History.points, History.points[1:])
		History.points[History.maxLen-1] = point
	} else {
		History.points = append(History.points, point)
	}
}

func (h *MetricsHistory) GetHistory() []DataPoint {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Return a copy to avoid race conditions
	result := make([]DataPoint, len(h.points))
	copy(result, h.points)
	return result
}
