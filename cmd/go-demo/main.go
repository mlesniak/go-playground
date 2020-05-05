package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/mlesniak/go-demo/pkg/authentication"
	logger "github.com/mlesniak/go-demo/pkg/log"
	"github.com/mlesniak/go-demo/pkg/version"
)

var log = logger.New()

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middlewares.
	e.Use(authentication.KeycloakWithConfig(authentication.KeycloakConfig{
		IgnoredURL: []string{
			"/api/login",
			"/api/version",
		},
	}))

	// Endpoints.
	version.AddVersionEndpoint(e)
	authentication.AddAuthenticationEndpoints(e)
	addAPIEndpoint(e)

	log.Info("Application started")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Info(e.Start(":" + port))
}

func addAPIEndpoint(e *echo.Echo) {
	type request struct {
		Number int `json:"number"`
	}

	type response struct {
		Number int `json:"number"`
	}

	e.POST("/api", func(c echo.Context) error {
		log := authentication.AddUser(log, c)

		var json request
		err := c.Bind(&json)
		if err != nil {
			log.WithField("error", err).Warn("Unable to parse json")
			return c.String(http.StatusBadRequest, "Unable to parse request")
		}
		log.WithField("number", json.Number).Info("Request received")
		resp := response{json.Number + 1}
		return c.JSON(http.StatusOK, resp)
	})
}
