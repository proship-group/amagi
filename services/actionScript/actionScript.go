package actionScript

import (
	"fmt"
	"net/http"

	"github.com/b-eee/amagi/api/httpd"

	utils "github.com/b-eee/amagi"
)

// TryActionScript try action script
func TryActionScript(w http.ResponseWriter, r *http.Request, _ http.HandlerFunc) {
	var s Script
	if err := httpd.DecodePostRequest(r.Body, &s); err != nil {
		utils.Info(fmt.Sprintf("error TryActionScript DecodePostRequest %v", err))
		return
	}

	if err := s.TryScript(); err != nil {
		httpd.HTTPError(w, err)
		return
	}

	httpd.RespondNilError(w)
}
