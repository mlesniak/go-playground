package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
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
	serve()
}

type foo struct {
	Version string `json:"version"`
}

func serve() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, Message())
	})

	// Endpoint to check correct deployment.
	e.GET("/api/version", func(c echo.Context) error {
		log.Info("Version info requested")
		commit := os.Getenv("COMMIT")
		if commit == "" {
			commit = "<No COMMIT environment variable set>"
		}
		return c.JSON(http.StatusOK, foo{commit})
	})

	e.POST("/api", func(c echo.Context) error {
		var json request
		err := c.Bind(&json)
		if err != nil {
			log.WithField("error", err).Warn("Unable to parse json")
			return c.String(http.StatusBadRequest, "Unable to parse request")
		}
		log.WithField("number", json.Number).Info("Request received")
		resp := response{json.Number+1}
		return c.JSON(http.StatusOK, resp)
	})

	log.Info("Application started")
	log.Info(e.Start(":8080"))
}

// Message returns a greeting string.
func Message() string {
	return "OK"
}

