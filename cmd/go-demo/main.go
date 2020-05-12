package main

import (
	"os"
	"time"

	ctx "context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mlesniak/go-demo/pkg/authentication"
	"github.com/mlesniak/go-demo/pkg/context"
	"github.com/mlesniak/go-demo/pkg/demo"
	"github.com/mlesniak/go-demo/pkg/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	c, _ := ctx.WithTimeout(ctx.Background(), 10*time.Second)
	options := options.Client().
		SetAuth(options.Credential{
			Username: "admin",
			Password: "admin",
		}).
		ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(c, options)
	if err != nil {
		panic(err)
	}
	api := client.Database("dev").Collection("api")

	c, _ = ctx.WithTimeout(ctx.Background(), 5*time.Second)
	_, err = api.InsertOne(c, bson.M{"name": "pi", "value": 3.14159})
	if err != nil {
		panic(err)
	}

}

func xmain() {
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

		LoginURL:   "/api/login",
		LogoutURL:  "/api/logout",
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
