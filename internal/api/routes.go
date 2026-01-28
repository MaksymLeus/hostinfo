package api

import "net/http"

func RegisterRoutes() {
	http.HandleFunc("/api/v1/info", InfoHandler)
	http.HandleFunc("/api/v1/cloud", CloudHandler)
	http.HandleFunc("/api/v1/kubernetes", KubernetesHandler)
}
