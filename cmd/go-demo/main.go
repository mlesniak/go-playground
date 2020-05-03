package main

import (
	"net/http"
	"os"

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
	addJWTMiddleware(e)

	// Endpoints.
	addVersionEndpoint(e)
	addAPIEndpoint(e)

	log.Info("Application started")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Info(e.Start(":" + port))
}

func addVersionEndpoint(e *echo.Echo) {
	e.GET("/api/version", func(c echo.Context) error {
		log := logger.AddUser(log, c)

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
		log := logger.AddUser(log, c)

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

func addJWTMiddleware(e *echo.Echo) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("NO JWT_SECRET set. Aborting.")
	}
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(secret),
		Skipper: func(c echo.Context) bool {
			// List of urls to ignore for authentication.
			ignoredURL := []string{
				"/api/version",
			}

			path := c.Request().URL.Path
			for _, v := range ignoredURL {
				if v == path {
					return true
				}
			}
			return false
		},
	}))
}
