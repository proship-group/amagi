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
		c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return err
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	return nil
}

// GinHTTPOk gin generic http response
func GinHTTPOk(c *gin.Context, resp gin.H) error {
	c.JSON(http.StatusOK, resp)
	return nil
}

// GinHTTPAnonOk gin generic http response writer that accepts interface response
func GinHTTPAnonOk(c *gin.Context, data interface{}) error {
	c.JSON(http.StatusOK, data)
	return nil
}

// GinJSONResponse gin json response
func GinJSONResponse(c *gin.Context, resp interface{}) error {
	c.JSON(http.StatusOK, resp)
	return nil
}

// GinJSONStatusResponse gin json response
func GinJSONStatusResponse(c *gin.Context, status int, resp interface{}) error {
	c.JSON(status, resp)
	return nil
}

// GinErrorResponse gin err response
func GinErrorResponse(c *gin.Context, status int, resp interface{}) {
	c.AbortWithStatusJSON(status, resp)
	return
}
