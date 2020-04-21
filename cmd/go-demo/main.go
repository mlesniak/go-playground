package main

import (
	"net/http"

	// "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func main() {
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
