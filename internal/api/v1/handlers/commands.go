package handlers

import (
	"io"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	probing "github.com/prometheus-community/pro-bing"
)

// -------------------- Ping --------------------
// Ping godoc
// @Summary Ping a host
// @Description Sends ICMP pings to the specified host and returns statistics
// @Tags commands
// @Accept json
// @Produce json
// @Param host query string true "Host to ping"
// @Success 200 {object} PingResult
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /ping [post]
func Ping(c echo.Context) error {
	hostParam := c.QueryParam("host")
	if hostParam == "" {
		return c.JSON(http.StatusBadRequest, "`host` query parameter required")
	}

	pinger, err := probing.NewPinger(hostParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	pinger.Count = 3
	pinger.Timeout = 5 * time.Second
	pinger.SetPrivileged(false)

	if err := pinger.Run(); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	stats := pinger.Statistics()
	return c.JSON(http.StatusOK, PingResult{
		Host:        hostParam,
		PacketsSent: stats.PacketsSent,
		PacketsRecv: stats.PacketsRecv,
		Loss:        stats.PacketLoss,
		MinRTT:      stats.MinRtt.String(),
		MaxRTT:      stats.MaxRtt.String(),
		AvgRTT:      stats.AvgRtt.String(),
	})
}

// -------------------- Curl / HTTP GET --------------------
// Curl godoc
// @Summary Perform an HTTP GET request
// @Description Performs an HTTP GET request to the specified URL and returns the response
// @Tags commands
// @Accept json
// @Produce json
// @Param url query string true "URL to fetch"
// @Success 200 {object} CurlResult
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /curl [post]
func Curl(c echo.Context) error {
	url := c.QueryParam("url")
	if url == "" {
		return c.JSON(http.StatusBadRequest, "`url` query parameter required")
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return c.JSON(http.StatusOK, CurlResult{
		URL:        url,
		StatusCode: resp.StatusCode,
		Body:       string(body),
	})
}

// -------------------- Dig / DNS lookup --------------------
// Dig godoc
// @Summary Perform a DNS lookup
// @Description Looks up IP addresses and CNAME for the specified host
// @Tags commands
// @Accept json
// @Produce json
// @Param host query string true "Host for DNS lookup"
// @Success 200 {object} DNSResult
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /dig [post]
func Dig(c echo.Context) error {
	hostParam := c.QueryParam("host")
	if hostParam == "" {
		return c.JSON(http.StatusBadRequest, "`host` query parameter required")
	}

	ips, err := net.LookupIP(hostParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	ipStrs := make([]string, len(ips))
	for i, ip := range ips {
		ipStrs[i] = ip.String()
	}

	cname, _ := net.LookupCNAME(hostParam)

	return c.JSON(http.StatusOK, DNSResult{
		Host:  hostParam,
		IPs:   ipStrs,
		CNAME: cname,
	})
}

// -------------------- TCP --------------------
// TCP godoc
// @Summary Check TCP connection
// @Description Tests a TCP connection to a host and port
// @Tags commands
// @Accept json
// @Produce json
// @Param host query string true "Host to connect to"
// @Param port query string true "Port to connect to"
// @Success 200 {object} TCPResult
// @Failure 400 {string} string "Bad Request"
// @Router /tcp [post]
func TCP(c echo.Context) error {
	hostParam := c.QueryParam("host")
	port := c.QueryParam("port")
	if hostParam == "" || port == "" {
		return c.JSON(http.StatusBadRequest, "`host` and `port` query parameters required")
	}

	addr := net.JoinHostPort(hostParam, port)
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	result := TCPResult{Host: hostParam, Port: port}

	if err != nil {
		result.Open = false
		result.Error = err.Error()
	} else {
		result.Open = true
		conn.Close()
	}

	return c.JSON(http.StatusOK, result)
}
