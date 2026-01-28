package api

import (
	"encoding/json"
	"net/http"

	"hostinfo/internal/host"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	info := host.Collect()
	writeJSON(w, info)
}

func CloudHandler(w http.ResponseWriter, r *http.Request) {
	info := host.DetectCloud()
	writeJSON(w, info)
}

func KubernetesHandler(w http.ResponseWriter, r *http.Request) {
	info := host.DetectKubernetes()
	writeJSON(w, info)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
