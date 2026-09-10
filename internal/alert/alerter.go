package alert

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type AlertPayload struct {
	Message   string    `json:"message"`
	Metric    string    `json:"metric"`
	Value     float64   `json:"value"`
	Threshold float64   `json:"threshold"`
	Timestamp time.Time `json:"timestamp"`
}

// StartAlertWatcher runs a background loop to check resource limits
func StartAlertWatcher() {
	// Check if alerts are explicitly enabled or if webhook URL is set
	enabledEnv := strings.ToLower(os.Getenv("HOSTINFO_ENABLE_ALERTS"))
	webhookURL := os.Getenv("ALERT_WEBHOOK_URL")

	if enabledEnv != "true" && webhookURL == "" {
		log.Println("Webhook alerting disabled (set HOSTINFO_ENABLE_ALERTS=true and ALERT_WEBHOOK_URL to enable)")
		return
	}

	if webhookURL == "" {
		log.Println("Webhook alerting disabled (ALERT_WEBHOOK_URL not set)")
		return
	}

	// Parse thresholds from env or use defaults (e.g. 90% CPU/Memory)
	cpuThreshold := getEnvFloat("ALERT_CPU_THRESHOLD", 90.0)
	memThreshold := getEnvFloat("ALERT_MEM_THRESHOLD", 90.0)

	intervalSec := getEnvInt("ALERT_CHECK_INTERVAL_SEC", 60)
	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)

	log.Printf("Starting alert watcher (CPU threshold: %.1f%%, Memory threshold: %.1f%%)", cpuThreshold, memThreshold)

	go func() {
		for range ticker.C {
			checkAndAlert(webhookURL, cpuThreshold, memThreshold)
		}
	}()
}

func checkAndAlert(webhookURL string, cpuThreshold, memThreshold float64) {
	// Check CPU
	if usage, err := cpu.Percent(0, false); err == nil && len(usage) > 0 {
		if usage[0] >= cpuThreshold {
			sendAlert(webhookURL, AlertPayload{
				Message:   "High CPU usage detected!",
				Metric:    "cpu",
				Value:     usage[0],
				Threshold: cpuThreshold,
				Timestamp: time.Now(),
			})
		}
	}

	// Check Memory
	if vm, err := mem.VirtualMemory(); err == nil && vm != nil {
		if vm.UsedPercent >= memThreshold {
			sendAlert(webhookURL, AlertPayload{
				Message:   "High Memory usage detected!",
				Metric:    "memory",
				Value:     vm.UsedPercent,
				Threshold: memThreshold,
				Timestamp: time.Now(),
			})
		}
	}
}

func sendAlert(url string, payload AlertPayload) {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal alert payload: %v", err)
		return
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to send webhook alert: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("Webhook alert returned non-success status: %d", resp.StatusCode)
	} else {
		log.Printf("Alert sent successfully for metric: %s (Value: %.2f)", payload.Metric, payload.Value)
	}
}

func getEnvFloat(key string, fallback float64) float64 {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.ParseFloat(val, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}
