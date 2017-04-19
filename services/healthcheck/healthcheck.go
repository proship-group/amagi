package healthcheck

import (
	"net/http"

	"github.com/b-eee/amagi/api/httpd"
)

// GetHealthCheck app healthcheck
func GetHealthCheck(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {
	w.WriteHeader(200)
	httpd.RespondToJSON(w, Healthy())
}
