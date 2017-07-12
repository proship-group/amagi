package helpers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

// GinJSONResponse gin json response
func GinJSONResponse(c *gin.Context, resp interface{}) error {
	c.JSON(http.StatusOK, resp)
	return nil
}
