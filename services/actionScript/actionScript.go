package actionScript

import (
	"fmt"

	"github.com/b-eee/amagi/api/httpd"
	"github.com/gin-gonic/gin"

	utils "github.com/b-eee/amagi"
)

// TryActionScript try action script
func TryActionScript(c *gin.Context) {
	var s Script
	if err := httpd.DecodePostRequest(c.Request.Body, &s); err != nil {
		utils.Info(fmt.Sprintf("error TryActionScript DecodePostRequest %v", err))
		return
	}

	s.ReplaceEnvVars(map[string]string{})

	if err := s.ExecuteScript(); err != nil {
		httpd.HTTPError(c.Writer, err)
		return
	}

	httpd.RespondNilError(c.Writer)
}
