package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/mlesniak/go-demo/pkg/authentication"
	"github.com/mlesniak/go-demo/pkg/context"
	"github.com/mlesniak/go-demo/pkg/demo"
	logger "github.com/mlesniak/go-demo/pkg/log"
	"github.com/mlesniak/go-demo/pkg/version"
)

var log = logger.New()

func main() {
	e := newEchoServer()

	// Middlewares.
	// TODO Use RequestId and add to custom context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &context.CustomContext{c}
			return next(cc)
		}
	})
	e.Use(authentication.KeycloakWithConfig(authentication.KeycloakConfig{
		IgnoredURL: []string{
			"/api/login",
			"/api/version",
		},
	}))

	// Endpoints.
	version.AddVersionEndpoint(e)
	authentication.AddAuthenticationEndpoints(e)
	demo.AddEndpoint(e)

	start(e)
}

func newEchoServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	return e
}

func start(e *echo.Echo) {
	log.Info("Application started")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Info(e.Start(":" + port))
}
