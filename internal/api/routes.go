package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/info", InfoHandler)
	mux.HandleFunc("/cloud", CloudHandler)
	mux.HandleFunc("/kubernetes", KubernetesHandler)

	mux.HandleFunc("/ping", PingHandler)
	mux.HandleFunc("/curl", CurlHandler)
	mux.HandleFunc("/dns", DigHandler)
	mux.HandleFunc("/tcp", TCPHandler)

}
