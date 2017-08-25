package healthcheck

import (
	"github.com/b-eee/amagi/api/httpd"
	"github.com/gin-gonic/gin"
)

// GetHealthCheck app healthcheck
func GetHealthCheck(c *gin.Context) {
	httpd.RespondToJSON(c.Writer, Healthy("GetHealthCheck", ResponseLogMessage))
}

// NonNextGetHealthCheck non next app healthcheck
func NonNextGetHealthCheck(c *gin.Context) {
	httpd.RespondToJSON(c.Writer, Healthy("NonNextGetHealthCheck", ResponseLogMessage))
}
