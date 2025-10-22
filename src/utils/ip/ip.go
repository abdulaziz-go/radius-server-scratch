package ip

import (
	"net"
)


func DetectIPVersion(ipStr string) string {
	if ipStr == "" {
		return ""
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}

	if ip.To4() != nil {
		return "v4"
	}

	return "v6"
}
