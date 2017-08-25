package net

import (
	"net"
	"os"
)

func dockerHostIP() string {
	str := "localhost"

	if host := os.Getenv("DOCKERHOST"); host != "" {
		return host
	}

	return str
}

func getCurrentIP() string {
	locIP := "local"
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			locIP = ipv4.String()
		}
	}
	return locIP
}

// GetCurrentHostIP get current host IP for log_stream channel name
func GetCurrentHostIP() string {
	var host string
	switch os.Getenv("ENV") {
	case "local":
		host = getCurrentIP()
	case "localhost":
		host = getCurrentIP()
	case "docker":
		host = dockerHostIP()
	default:
		host = "localhost"
	}

	return host
}

// AppHostName get app hostname
func AppHostName() string {
	host, err := os.Hostname()
	if err != nil {
		return "hostname_not_found"
	}

	return host
}
