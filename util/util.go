package util

import (
	"net"
	"strings"

	"github.com/labstack/gommon/random"
)

func ID() string {
	return random.String(8, random.Lowercase, random.Numeric)
}

// StripPort strips port from the host.
// https://golang.org/pkg/net/url/#URL.Hostname
func StripPort(host string) string {
	colon := strings.IndexByte(host, ':')
	if colon == -1 {
		return host
	}
	if i := strings.IndexByte(host, ']'); i != -1 {
		return strings.TrimPrefix(host[:i], "[")
	}
	return host[:colon]
}

func PrivateIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip := ipNet.IP
				_, cidr24BitBlock, _ := net.ParseCIDR("10.0.0.0/8")
				_, cidr20BitBlock, _ := net.ParseCIDR("172.16.0.0/12")
				_, cidr16BitBlock, _ := net.ParseCIDR("192.168.0.0/16")
				private := cidr24BitBlock.Contains(ip) || cidr20BitBlock.Contains(ip) || cidr16BitBlock.Contains(ip)
				if private {
					return ip.String()
				}
			}
		}
	}
	return ""
}
