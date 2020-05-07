package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/mlesniak/go-demo/pkg/authentication"
	"github.com/mlesniak/go-demo/pkg/context"
	"github.com/mlesniak/go-demo/pkg/demo"
	"github.com/mlesniak/go-demo/pkg/version"
)

func main() {
	e := newEchoServer()

	// Middlewares.
	// TODO Use RequestId and add to custom context
	e.Use(context.CreateCustomContext)
	e.Use(authentication.KeycloakWithConfig(e, authentication.KeycloakConfig{
		Protocol: "http",
		Hostname: "localhost",
		Port: "8081",
		Realm: "mlesniak",

		LoginURL: "/api/login",
		LogoutURL: "/api/logout",
		IgnoredURL: []string{
			"/api/login",
			"/api/version",
		},
	}))

	// Endpoints.
	version.AddVersionEndpoint(e)
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Info().Str("port", port).Msg("Application started")
	err := e.Start(":" + port)
	if err != nil {
		log.Panic().Msg(err.Error())
	}
}
