package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type request struct {
	Number int `json:"number"`
}

type response struct {
	Number int `json:"number"`
}

func main() {
	initializeLogging()
	serve()
}

func serve() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, Message())
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

func initializeLogging() {
	// On local environment ignore all file logging.
	env := os.Getenv("ENVIRONMENT")
	if env == "local" {
		log.Info("Local environment, using solely console output")
		return
	}

	// Set to JSON.
	formatter := &log.JSONFormatter{
	  	FieldMap: log.FieldMap{
			 log.FieldKeyTime:  "@timestamp",
			 log.FieldKeyLevel: "level",
			 log.FieldKeyMsg:   "message",
	   },
	}
	log.SetOutput(os.Stdout)
	log.SetFormatter(formatter)
}
