package healthcheck

import (
	"net/http"

	"github.com/b-eee/amagi/api/httpd"
)

// GetHealthCheck app healthcheck
func GetHealthCheck(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {
	httpd.RespondToJSON(w, Healthy("GetHealthCheck", ResponseLogMessage))
}

// NonNextGetHealthCheck non next app healthcheck
func NonNextGetHealthCheck(w http.ResponseWriter, r *http.Request) {
	httpd.RespondToJSON(w, Healthy("NonNextGetHealthCheck", ResponseLogMessage))
}
