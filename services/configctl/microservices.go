package configctl

import (
	"fmt"
	"os"
)

var (
	// ConfigCtlPort config controller port
	ConfigCtlPort = 8083
)

// ConfCtLHost config controller host name
func ConfCtLHost() string {
	var host string

	env := os.Getenv("ENV")
	switch env {
	case "local":
		host = "localhost"
	default:
		host = "beee-configctl"
	}

	return fmt.Sprintf("%v://%v:%v", ConfCtLProtocol(env), host, ConfigCtlPort)
}

// ConfCtLProtocol config controller protocol
func ConfCtLProtocol(env string) string {
	var protocol string
	switch env {
	case "local":
		protocol = "http"

	default:
		protocol = "http"
	}

	return protocol
}

// ImporterURL importer URL
func ImporterURL(port int) string {
	var importerURL string
	switch os.Getenv("ENV") {
	case "local":
		importerURL = fmt.Sprintf("localhost:%v", port)
	case "docker":
		importerURL = fmt.Sprintf("importer:%v", port)
	default:
		importerURL = fmt.Sprintf("beee-importer:%v", port)
	}
	return importerURL
}
