package demo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mlesniak/go-demo/pkg/authentication"
	logger "github.com/mlesniak/go-demo/pkg/log"
)

var log = logger.New()

// AddEndpoint adds the endpoint for testing purposes.
func AddEndpoint(e *echo.Echo) {
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
			// TODO Use error object
			return c.String(http.StatusBadRequest, "Unable to parse request")
		}
		log.WithField("number", json.Number).Info("Request received")
		resp := response{json.Number + 1}
		return c.JSON(http.StatusOK, resp)
	})
}
