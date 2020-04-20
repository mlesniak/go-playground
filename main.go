package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, Message())
	})
	e.Logger.Fatal(e.Start(":8080"))
}

// Message returns a greeting string.
func Message() string {
	return "Hello, world"
}