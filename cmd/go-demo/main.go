package main

import (
	"os"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
	"github.com/mlesniak/go-demo/pkg/authentication"
	"github.com/mlesniak/go-demo/pkg/context"
	"github.com/mlesniak/go-demo/pkg/demo"
	"github.com/mlesniak/go-demo/pkg/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	configureLogging()
	e := newEchoServer()

	// Middlewares.
	e.Use(middleware.RequestID())
	e.Use(context.CreateCustomContext)
	e.Use(authentication.KeycloakWithConfig(e, authentication.KeycloakConfig{
		Protocol: os.Getenv("KEYCLOAK_PROTOCOL"),
		Hostname: os.Getenv("KEYCLOAK_HOST"),
		Port:     os.Getenv("KEYCLOAK_PORT"),
		Realm:    os.Getenv("KEYCLOAK_REALM"),
		Client:   os.Getenv("KEYCLOAK_CLIENT"),

		LoginURL:  "/api/login",
		LogoutURL: "/api/logout",
		RefreshURL: "/api/refresh",
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

func configureLogging() {
	// On local environment ignore all file logging.
	env := os.Getenv("ENVIRONMENT")
	if env == "local" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Info().Msg("Local environment, using solely console output")
	} else {
		// Logging to json by default.
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
		zerolog.TimestampFieldName = "@timestamp"
	}

	// Add commit hash to each entry.
	commit := os.Getenv("COMMIT")
	if commit == "" {
		log.Warn().Msg("COMMIT environment variable not set, won't record commit hash on each request")
	} else {
		log.Logger = log.With().Str("commit", commit).Logger()
	}
}
