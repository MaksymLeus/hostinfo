// internal/health/health.go
package health

import (
	"hostinfo/internal/custom"
	"net/http"
)

// /healthz/live
func LiveHandler(w http.ResponseWriter, r *http.Request) {
	custom.WriteJSON(w, http.StatusOK, CheckResult{
		Status: "ok",
	})
}

// /healthz/ready
func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	checks := BasicChecks()

	custom.WriteJSON(w, http.StatusOK, CheckResult{
		Status: "ok",
		Checks: checks,
	})
}

// /healthz
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	custom.WriteJSON(w, http.StatusOK, CheckResult{
		Status: "ok",
	})
}
