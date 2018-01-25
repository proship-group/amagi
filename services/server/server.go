package server

import (
	"fmt"
	"os"

	utils "github.com/b-eee/amagi"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/contrib/sentry"
	"github.com/gin-gonic/gin"
)

type (
	// StartUp server startup struct
	StartUp struct {
		GinEngine       *gin.Engine
		StartupPackacge interface{}
	}
)

var (
	ginEngine      *gin.Engine
	startUpPackage interface{}
)

// NewWebHost new gin webhost
func NewWebHost() *StartUp {
	s := StartUp{
		GinEngine: gin.New(),
	}

	// s.GinEngine.Use(gin.Logger())
	return &s
}

// UseStartUp use startup package
func (s *StartUp) UseStartUp(packageName interface{}) *StartUp {
	s.StartupPackacge = packageName

	return s
}

// UseExternalAPIRoutes set external API routes
func (s *StartUp) UseExternalAPIRoutes(startAPPAPIRoutes func(*gin.Engine) error) *StartUp {
	startAPPAPIRoutes(s.GinEngine)

	return s
}

// UseGinRecovery use gin http recovery
func (s *StartUp) UseGinRecovery() *StartUp {
	(*s).GinEngine.Use(gin.Recovery())

	return s
}

// UseGinLogger use gin default logger
func (s *StartUp) UseGinLogger() *StartUp {
	// (*s.GinEngine).Use(gin.Logger())
	(*s.GinEngine).Use(gin.Logger())
	return s
}

// UseMiddlewares add custom middlewares to gin
func (s *StartUp) UseMiddlewares(middlewares ...gin.HandlerFunc) *StartUp {
	for _, middleware := range middlewares {
		(*s).GinEngine.Use(middleware)
	}

	return s
}

// GetGinEngine get current gin engine
func (s *StartUp) GetGinEngine() *gin.Engine {
	return ginEngine
}

// Run run the server stack
func (s *StartUp) Run(port string) *StartUp {
	(*s).GinEngine.Run(port)

	return s
}

// DatabaseStartups initialize database
func (s *StartUp) DatabaseStartups(dbs ...func()) *StartUp {
	for _, startdb := range dbs {
		startdb()
	}

	return s
}

// UseSentry use sentry middleware
func (s *StartUp) UseSentry() *StartUp {
	if ok := os.Getenv("SENTRY_DSN"); len(ok) != 0 {
		utils.Info(fmt.Sprintf("starting sentry.."))
		(*s).GinEngine.Use(sentry.Recovery(raven.DefaultClient, false))
	}

	return s
}

// ConfigureServices configure services from func
func (s *StartUp) ConfigureServices(services func()) *StartUp {
	services()

	return s
}
