package uptime

import (
	"net/http"

	"github.com/b-eee/amagi/api/httpd"
)

// HealthCheck app healthcheck
func HealthCheck(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {
	w.WriteHeader(200)
	httpd.RespondToJSON(w, Healthy())
}
