package helpers

import (
	"net/http"
	"os"

	"gopkg.in/gin-gonic/gin.v1"
)

var (
	// SendHTTPErrorEnv send http error env flag
	SendHTTPErrorEnv = "SEND_HTTP_ERROR"
)

// GinHTTPError gin framework for using http error interface handler
func GinHTTPError(c *gin.Context, err error) error {
	if flag := os.Getenv(SendHTTPErrorEnv); flag == "false" {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 500})
		return err
	}

	c.JSON(http.StatusInternalServerError, gin.H{"status": 500, "err": err.Error()})
	return nil
}

// GinHTTPOk gin generic http response
func GinHTTPOk(c *gin.Context, resp gin.H) error {
	c.JSON(http.StatusOK, resp)
	return nil
}
