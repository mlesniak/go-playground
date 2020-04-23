package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func main() {
	initializeLogging()
	serve()
}

func serve() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.GET("/", func(c echo.Context) error {
		log.Info("Demo logging context")
		return c.String(http.StatusOK, Message())
	})

	log.Info("Application started")
	log.Info(e.Start(":8080"))
}

// Message returns a greeting string.
func Message() string {
	return "Hello, world"
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
