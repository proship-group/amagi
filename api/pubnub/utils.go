package pubnub

import (
	"net"
	"os"
	"strings"

	"github.com/b-eee/amagi/api/slack"
)

// ChanName channel name constructor
func ChanName(names ...string) string {
	names = append(names, setHost())

	return strings.Join(names, "_")
}

func setHost() string {
	var host string
	switch os.Getenv("ENV") {
	case "local":
		host = strings.Join([]string{getCurrentIP()}, "_")
	case "docker":
		host = strings.Join([]string{dockerHostIP()}, "_")
	default:
		host = slack.HostName()

	}

	return host
}

func dockerHostIP() string {
	ip, err := net.ResolveIPAddr("ip", "dockerhost")
	if err != nil {
		return "dockerhost"
	}
	return ip.String()
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
