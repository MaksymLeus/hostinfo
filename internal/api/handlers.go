package api

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"hostinfo/internal/host"

	probing "github.com/prometheus-community/pro-bing"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	info := host.Collect()
	writeJSON(w, http.StatusOK, info)
}

// -------------------- Ping --------------------

func PingHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		writeJSON(w, http.StatusBadRequest, "`host` query parameter required")
		return
	}

	pinger, err := probing.NewPinger(host)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	pinger.Count = 3
	pinger.Timeout = time.Second * 5
	pinger.SetPrivileged(false)

	if err := pinger.Run(); err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	stats := pinger.Statistics()
	json.NewEncoder(w).Encode(PingResult{
		Host:        host,
		PacketsSent: stats.PacketsSent,
		PacketsRecv: stats.PacketsRecv,
		Loss:        stats.PacketLoss,
		MinRTT:      stats.MinRtt.String(),
		MaxRTT:      stats.MaxRtt.String(),
		AvgRTT:      stats.AvgRtt.String(),
	})
}

// -------------------- Curl / HTTP GET --------------------

func CurlHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		writeJSON(w, http.StatusBadRequest, "`url` query parameter required")
		return
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	json.NewEncoder(w).Encode(CurlResult{
		URL:        url,
		StatusCode: resp.StatusCode,
		Body:       string(body),
	})
}

// -------------------- Dig / DNS lookup --------------------

func DigHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		writeJSON(w, http.StatusBadRequest, "`host` query parameter required")
		return
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	var ipStrs []string
	for _, ip := range ips {
		ipStrs = append(ipStrs, ip.String())
	}

	cname, _ := net.LookupCNAME(host)

	json.NewEncoder(w).Encode(DNSResult{
		Host:  host,
		IPs:   ipStrs,
		CNAME: cname,
	})
}

// -------------------- TCP --------------------
func TCPHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	port := r.URL.Query().Get("port")
	if host == "" || port == "" {
		writeJSON(w, http.StatusBadRequest, "`host` and `port` query parameters required")
		return
	}

	addr := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	result := TCPResult{Host: host, Port: port}

	if err != nil {
		result.Open = false
		result.Error = err.Error()
	} else {
		result.Open = true
		conn.Close()
	}

	json.NewEncoder(w).Encode(result)
}

// -------------------- Other --------------------

func CloudHandler(w http.ResponseWriter, r *http.Request) {
	info := host.DetectCloud()
	writeJSON(w, http.StatusOK, info)
}

func KubernetesHandler(w http.ResponseWriter, r *http.Request) {
	info := host.DetectKubernetes()
	writeJSON(w, http.StatusOK, info)
}

// Standard JSON
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var toWrite any

	// Handle error vs success automatically
	if status >= 400 {
		switch val := v.(type) {
		case string:
			toWrite = map[string]string{"error": val}
		default:
			toWrite = map[string]any{"error": val}
		}
	} else {
		switch val := v.(type) {
		case string:
			toWrite = map[string]string{"message": val}
		default:
			toWrite = val
		}
	}

	if err := json.NewEncoder(w).Encode(toWrite); err != nil {
		// fallback if encoding fails
		http.Error(w, `{"error":"failed to encode JSON"}`, http.StatusInternalServerError)
	}
}
