package handlers

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

func Info(c echo.Context) error {
	startTime := time.Now()

	return c.JSON(http.StatusOK, HostInfo{
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
		UpdatedAt: startTime.Format(time.RFC3339),
		Uptime:    getUptime(startTime),

		// Environment
		Env: getEnv(),

		// Hardware
		CPU:    getCPU(),
		Memory: getMemory(),
		Load:   getLoad(),

		// Platform
		Cloud:      DetectCloud(),
		Runtime:    DetectRuntime(),
		Kubernetes: DetectKubernetes(),
	})
}

// -------- System -------- //
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

func getEnv() map[string]string {
	// Added FF_ENABLE_ENVIERMENT_VARIABLE feature flag for security reasons
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
	if os.Getenv("FF_ENVIRONMENT_VARIABLES") != "true" {
		return make(map[string]string)
	}
	return env
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
