package auth

import "github.com/gin-gonic/gin"

type (
	// RouteHandler main route handler interface
	RouteHandler map[string]Handler

	// Handler api route handler and method
	Handler struct {
		Method string
		Func   func(*gin.Context)
	}
)
