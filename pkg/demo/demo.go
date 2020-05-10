package demo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mlesniak/go-demo/pkg/context"
	"github.com/mlesniak/go-demo/pkg/response"
)

// AddEndpoint adds the endpoint for testing purposes.
func AddEndpoint(e *echo.Echo) {
	type request struct {
		Number int `json:"number"`
	}

	type numberResponse struct {
		Number int `json:"number"`
	}

	e.POST("/api", func(cc echo.Context) error {
		c := cc.(*context.CustomContext)
		log := c.Log()

		var json request
		err := c.Bind(&json)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.NewError("Unable to parse request"))
		}
		log.Info().Int("number", json.Number).Msg("Computing result")
		resp := numberResponse{json.Number + 1}
		return c.JSON(http.StatusOK, resp)
	})
}
