package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// Containers godoc
// @Summary List local Docker containers
// @Description Returns a list of running Docker containers if the Docker socket is mounted
// @Tags runtime
// @Accept json
// @Produce json
// @Success 200 {array} ContainerInfo
// @Router /containers [get]
func Containers(c echo.Context) error {
	containers := getLocalContainers()
	return c.JSON(http.StatusOK, containers)
}

type dockerContainer struct {
	ID      string `json:"Id"`
	Names   []string
	Image   string
	State   string
	Status  string
	Created int64
	Ports   []struct {
		PrivatePort uint16
		PublicPort  uint16
		Type        string
	}
}

func getLocalDockerClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}
}

func getLocalContainers() []ContainerInfo {
	var result []ContainerInfo

	if _, err := os.Stat("/var/run/docker.sock"); err != nil {
		return result // Docker socket not mounted/available
	}

	client := getLocalDockerClient()
	reqURL := "http://unix/containers/json?all=true"

	resp, err := client.Get(reqURL)
	if err != nil {
		return result
	}
	defer resp.Body.Close()

	var rawContainers []dockerContainer
	if err := json.NewDecoder(resp.Body).Decode(&rawContainers); err != nil {
		return result
	}

	for _, rc := range rawContainers {
		name := ""
		if len(rc.Names) > 0 {
			name = rc.Names[0]
			// Trim leading slash if present
			if len(name) > 0 && name[0] == '/' {
				name = name[1:]
			}
		}

		var ports []string
		for _, p := range rc.Ports {
			if p.PublicPort > 0 {
				portStr := fmt.Sprintf("%d:%d", p.PrivatePort, p.PublicPort)
				ports = append(ports, portStr)
			}
		}

		result = append(result, ContainerInfo{
			ID:      rc.ID[:12], // Short ID
			Name:    name,
			Image:   rc.Image,
			State:   rc.State,
			Status:  rc.Status,
			Created: rc.Created,
			Ports:   ports,
		})
	}

	fmt.Println("Detected containers:", len(result))
	return result
}
