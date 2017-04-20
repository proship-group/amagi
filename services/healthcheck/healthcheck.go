package healthcheck

import (
	"net/http"

	"github.com/b-eee/amagi/api/httpd"
)

// GetHealthCheck app healthcheck
func GetHealthCheck(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {
	ResponseLogMessage("GetHealthCheck")
	httpd.RespondToJSON(w, Healthy())
}

// NonNextGetHealthCheck non next app healthcheck
func NonNextGetHealthCheck(w http.ResponseWriter, r *http.Request) {
	ResponseLogMessage("NonNextGetHealthCheck")
	httpd.RespondToJSON(w, Healthy())
}
