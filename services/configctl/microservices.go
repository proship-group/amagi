package configctl

import (
	"fmt"

	"github.com/b-eee/amagi/services/externalSvc"
)

// ConfCtLHost config controller host name
func ConfCtLHost() string {
	return fmt.Sprintf("%v://%v:%v", ConfCtLProtocol(env), externalSvc.EnvConfigctlHost, externalSvc.EnvConfigctlPort)
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
func ImporterURL() string {
	return fmt.Sprintf("%s:%s", externalSvc.EnvImporterHost, externalSvc.EnvImporterPort)
}

// ApicoreURL apicore cluster/dev URL
func ApicoreURL() string {
	return fmt.Sprintf("%s:%s", externalSvc.EnvApicoreHost, externalSvc.EnvApicorePort)
}
