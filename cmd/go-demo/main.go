package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type request struct {
	Number int `json:"number"`
}

type response struct {
	Number int `json:"number"`
}

// Set global field values for log entry. According to https://github.com/sirupsen/logrus and
// e.g. https://github.com/sirupsen/logrus/pull/653#issuecomment-339763585, we might use our
// own log package in the future.
var log *logrus.Entry
func init() {
	commit := os.Getenv("COMMIT")
	if commit == "" {
		logrus.Warn("Environment variable 'commit' not set, unable to record hash in every log")
		log = logrus.NewEntry(logrus.New())
		return
	}

	log = logrus.WithFields(logrus.Fields{"commit":commit})
}

func main() {
	initializeLogging()
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

func initializeLogging() {
	// On local environment ignore all file logging.
	env := os.Getenv("ENVIRONMENT")
	if env == "local" {
		log.Info("Local environment, using solely console output")
		return
	}

	// Set to JSON.
	formatter := &logrus.JSONFormatter{
	  	FieldMap: logrus.FieldMap{
			 logrus.FieldKeyTime:  "@timestamp",
			 logrus.FieldKeyLevel: "level",
			 logrus.FieldKeyMsg:   "message",
	   },
	}
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(formatter)
}
