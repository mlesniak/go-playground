package main

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	logger "github.com/mlesniak/go-demo/pkg/log"
)

var log = logger.New()

type request struct {
	Number int `json:"number"`
}

type response struct {
	Number int `json:"number"`
}

func main() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middlewares.
	e.Use(middleware.JWT([]byte("secret")))

	// Endpoints.
	addVersionEndpoint(e)
	addAPIEndpoint(e)

	log.Info("Application started")
	log.Info(e.Start(":8080"))
}

func addVersionEndpoint(e *echo.Echo) {
	e.GET("/api/version", func(c echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		log = log.WithField("user", claims["user"])

		log.Info("Version info requested")
		commit := os.Getenv("COMMIT")
		if commit == "" {
			commit = "<No COMMIT environment variable set>"
		}
		return c.JSON(http.StatusOK, struct {
			Version string `json:"version"`
		}{commit})
	})
}

func addAPIEndpoint(e *echo.Echo) {
	e.POST("/api", func(c echo.Context) error {
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
