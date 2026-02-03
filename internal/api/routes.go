package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/info", corsFunc(InfoHandler))
	mux.HandleFunc("/cloud", CloudHandler)
	mux.HandleFunc("/kubernetes", KubernetesHandler)

	mux.HandleFunc("/ping", PingHandler)
	mux.HandleFunc("/curl", CurlHandler)
	mux.HandleFunc("/dns", DigHandler)
	mux.HandleFunc("/tcp", TCPHandler)

}

func corsFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}
