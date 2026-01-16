package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

type HostInfo struct {
	Hostname   string
	IPs        []string
	MACs       []string
	OS         string
	Arch       string
	GoVersion  string
	StartTime  string
	Now        string
	Env        map[string]string
	Cloud      CloudInfo
	Kubernetes KubernetesInfo
}

type CloudInfo struct {
	Provider string
	Region   string
	Zone     string
	Instance string
	Extra    map[string]string
}

type KubernetesInfo struct {
	Enabled        bool
	PodName        string
	PodNamespace   string
	PodIP          string
	NodeName       string
	ServiceAccount string
	Container      string
}

var startTime = time.Now()

func main() {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))

	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		info := HostInfo{
			Hostname:  getHostname(),
			IPs:       getIPs(),
			MACs:      getMACs(),
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
			GoVersion: runtime.Version(),
			StartTime: startTime.Format(time.RFC3339),
			Now:       time.Now().Format(time.RFC3339),
			Env:       getEnv(),
			Cloud:     detectCloud(),
		}

		_ = tmpl.Execute(w, info)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server listening on :%s --> http://localhost:%s", port, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))

}

func getHostname() string {
	h, _ := os.Hostname()
	return h
}

func getIPs() []string {
	var ips []string
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}

func getMACs() []string {
	var macs []string
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		if i.HardwareAddr != nil {
			macs = append(macs, i.HardwareAddr.String())
		}
	}
	return macs
}

func getEnv() map[string]string {
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
	return env
}

func detectCloud() CloudInfo {
	if aws := detectAWS(); aws.Provider != "" {
		return aws
	}
	if gcp := detectGCP(); gcp.Provider != "" {
		return gcp
	}
	if azure := detectAzure(); azure.Provider != "" {
		return azure
	}
	return CloudInfo{Provider: "local"}
}

func detectAWS() CloudInfo {
	client := http.Client{Timeout: 500 * time.Millisecond}

	req, _ := http.NewRequest("GET",
		"http://169.254.169.254/latest/meta-data/instance-id", nil)

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return CloudInfo{}
	}
	defer resp.Body.Close()

	id, _ := io.ReadAll(resp.Body)

	region := awsMeta("placement/region")
	zone := awsMeta("placement/availability-zone")

	return CloudInfo{
		Provider: "aws",
		Region:   region,
		Zone:     zone,
		Instance: string(id),
		Extra: map[string]string{
			"AMI":  awsMeta("ami-id"),
			"Type": awsMeta("instance-type"),
		},
	}
}

func awsMeta(path string) string {
	client := http.Client{Timeout: 300 * time.Millisecond}
	resp, err := client.Get("http://169.254.169.254/latest/meta-data/" + path)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b)
}

func detectGCP() CloudInfo {
	client := http.Client{Timeout: 500 * time.Millisecond}

	req, _ := http.NewRequest("GET",
		"http://metadata.google.internal/computeMetadata/v1/instance/id", nil)
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return CloudInfo{}
	}
	defer resp.Body.Close()

	id, _ := io.ReadAll(resp.Body)

	return CloudInfo{
		Provider: "gcp",
		Region:   gcpMeta("instance/region"),
		Zone:     gcpMeta("instance/zone"),
		Instance: string(id),
		Extra: map[string]string{
			"Machine": gcpMeta("instance/machine-type"),
			"Project": gcpMeta("project/project-id"),
		},
	}
}

func gcpMeta(path string) string {
	client := http.Client{Timeout: 300 * time.Millisecond}
	req, _ := http.NewRequest("GET",
		"http://metadata.google.internal/computeMetadata/v1/"+path, nil)
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b)
}

func detectAzure() CloudInfo {
	client := http.Client{Timeout: 500 * time.Millisecond}

	req, _ := http.NewRequest("GET",
		"http://169.254.169.254/metadata/instance?api-version=2021-02-01", nil)
	req.Header.Set("Metadata", "true")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return CloudInfo{}
	}
	defer resp.Body.Close()

	return CloudInfo{
		Provider: "azure",
		Extra: map[string]string{
			"VM": "Azure VM detected",
		},
	}
}

func detectKubernetes() KubernetesInfo {
	// Always present in k8s
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		return KubernetesInfo{Enabled: false}
	}

	return KubernetesInfo{
		Enabled:        false, //true, set to false to avoid showing k8s info in the UI until all is verified
		PodName:        "POD_NAME",
		PodNamespace:   "POD_NAMESPACE",
		PodIP:          "POD_IP",
		NodeName:       "NODE_NAME",
		ServiceAccount: "SERVICE_ACCOUNT",
		Container:      "CONTAINER_NAME",
	}
}
