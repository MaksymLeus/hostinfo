package health

import "time"

// Lightweight checks only
func BasicChecks() map[string]string {
	return map[string]string{
		"uptime": time.Since(startTime).String(),
	}
}

var startTime = time.Now()
