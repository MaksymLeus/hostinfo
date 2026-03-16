package handlers

import (
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Cloud(c echo.Context) error {
	info := DetectCloud()
	return c.JSON(http.StatusOK, info)
}

func DetectCloud() CloudInfo {
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
