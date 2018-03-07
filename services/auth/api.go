package auth

import (
	"fmt"
	"time"

	"github.com/b-eee/amagi/api/httpd"
	"github.com/b-eee/amagi/helpers"

	"github.com/gin-gonic/gin"

	utils "github.com/b-eee/amagi"
	jwt "github.com/dgrijalva/jwt-go"
)

type (
	// ExportedHandler handler interface
	ExportedHandler struct {
		Handler func(*gin.Context)
		Method  string
		Path    string
	}

	// Container auth type container
	Container struct {
		Container        interface{}
		PrefixedPath     []string
		LoginHandler     func(gin.H) error
		MapClaims        func() jwt.MapClaims
		ExportedHandlers []ExportedHandler
	}
)

// ImportAuthAPIs generate and import auth paths login/logout/session handlers to caller
func (cn *Container) ImportAuthAPIs(route *gin.Engine) []ExportedHandler {
	cn.ExportedHandlers = []ExportedHandler{}
	for _, prefixedPath := range cn.PrefixedPath {
		cn.ExportedHandlers =
			append(cn.ExportedHandlers, ExportedHandler{Path: fmt.Sprintf("%s/%s", prefixedPath, "/login"), Method: "POST", Handler: cn.Login})
	}

	return cn.ExportedHandlers
}

func (cn *Container) Login(c *gin.Context) {
	if err := httpd.DecodePostRequest(c.Request.Body, cn.Container); err != nil {
		helpers.GinHTTPError(c, utils.Error(fmt.Sprintf("error in Decoding Login %v", err)))
		return
	}

	tokenResp := cn.CreateTokenEndpoint(c)
	// LoginHandler is a func that user setups
	if err := cn.LoginHandler(tokenResp); err != nil {
		helpers.GinHTTPError(c, utils.Error(fmt.Sprintf("error on LoginHandler %v", err)))
		return
	}

	helpers.GinJSONResponse(c, tokenResp)
}

// Logout logout api
func (cn *Container) Logout(c *gin.Context) {}

// AuthenticateCurrentPath authenticate current API path
func (cn *Container) AuthenticateCurrentPath(c *gin.Context) {}

// CreateTokenEndpoint create/generate and sign the jwt token
func (cn *Container) CreateTokenEndpoint(c *gin.Context) gin.H {
	mapClaims := cn.MapClaims()

	mapClaims["exp"] = time.Now().Add(time.Second * time.Duration(sessionExpireOffSet())).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		helpers.GinHTTPError(c, fmt.Errorf("error on creating token"))
		return gin.H{}
	}

	// cn.Container.
	// 	helpers.GinHTTPAnonOk(c, gin.H{"token": tokenString})
	return gin.H{"token": tokenString}
}
