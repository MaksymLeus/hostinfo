package handlers

import (
	"net"
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
