package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/mlesniak/go-demo/pkg/authentication"
	"github.com/mlesniak/go-demo/pkg/demo"
	logger "github.com/mlesniak/go-demo/pkg/log"
	"github.com/mlesniak/go-demo/pkg/version"
)

var log = logger.New()

func main() {
	e := newEchoServer()

	// Middlewares.
	// TODO Use RequestId and add to custom context
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

