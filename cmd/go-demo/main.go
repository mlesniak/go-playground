package main

import (
	"net/http"
	"os"

	// "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func main() {
	initializeLogging()
	serve()
}

func serve() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		log.Info("Demo logging context")
		return c.String(http.StatusOK, Message())
	})

	e.Logger.Info(e.Start(":8080"))
}

// Message returns a greeting string.
func Message() string {
	return "Hello, world"
}

func initializeLogging() {
	// On non-local environments, log to a file using json.
	// (To be done)

	log.SetFormatter(&logrus.JSONFormatter{})
	const logfile = "main.log.json"
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Unable to create logfile?")
	}

	log.SetOutput(file)
}
