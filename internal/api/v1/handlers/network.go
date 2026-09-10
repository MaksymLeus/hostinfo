package handlers

import (
	"net"

	gopsnet "github.com/shirou/gopsutil/v3/net"
)

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

// GetNetForMetrics gathers detailed interface I/O stats
func getNetForMetrics() []NetInfo {
	var netInfos []NetInfo

	// Get per-interface I/O counters
	ioCounters, err := gopsnet.IOCounters(true)
	if err != nil {
		return netInfos
	}

	// Get system network interfaces for IP mapping
	ifaces, err := net.Interfaces()
	if err != nil {
		return netInfos
	}

	ifaceMap := make(map[string][]string)
	for _, i := range ifaces {
		var addrs []string
		if addrsList, err := i.Addrs(); err == nil {
			for _, a := range addrsList {
				if ipnet, ok := a.(*net.IPNet); ok {
					addrs = append(addrs, ipnet.IP.String())
				}
			}
		}
		ifaceMap[i.Name] = addrs
	}

	for _, io := range ioCounters {
		netInfos = append(netInfos, NetInfo{
			Name:        io.Name,
			Addrs:       ifaceMap[io.Name],
			BytesSent:   io.BytesSent,
			BytesRecv:   io.BytesRecv,
			PacketsSent: io.PacketsSent,
			PacketsRecv: io.PacketsRecv,
			ErrorsIn:    io.Errin,
			ErrorsOut:   io.Errout,
		})
	}

	return netInfos
}
