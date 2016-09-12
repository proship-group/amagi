package pubnub

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/b-eee/amagi/api/slack"
)

// ChanName channel name constructor
func ChanName(names ...string) string {
	names = append(names, SetHost())

	return strings.Join(names, "_")
}

// SetHost set/get host
func SetHost() string {
	var host string
	switch os.Getenv("ENV") {
	case "local":
		host = strings.Join([]string{getCurrentIP()}, "_")
	case "localhost":
		host = strings.Join([]string{getCurrentIP()}, "_")
	case "docker":
		host = strings.Join([]string{dockerHostIP()}, "_")
	default:
		host = dockerHostIP()

	}

	return host
}

func dockerHostIP() string {
	str := "localhost"

	if host := os.Getenv("DOCKERHOST"); host != "" {
		return host
	}

	// ip, err := net.ResolveIPAddr("ip", "dockerhost")
	// if err == nil {
	// 	// fmt.Printf("cant resolve ip dockerhost %v\n", err)
	// 	// return "dockerhost"
	// 	return ip.String()
	// }
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

func formatHostName(message, hostname string, host *slack.Host) string {
	return fmt.Sprintf("%v | %v", slack.ColorizedHost(host), message)
}
