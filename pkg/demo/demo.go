package demo

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mlesniak/go-demo/pkg/context"
	"github.com/mlesniak/go-demo/pkg/database"
	"github.com/mlesniak/go-demo/pkg/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddEndpoint adds the endpoint for testing purposes.
// We mix everything here, since it's just a simple demo endpoint.
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

		// Check if we have something in the database.
		res := database.Get().FindOne(database.DefaultContext(), bson.M{"number": json.Number})
		if !errors.Is(res.Err(), mongo.ErrNoDocuments) {
			type result struct {
				Number int `json:"number"`
				Value  int `json:"value"`
			}
			var row result
			res.Decode(&row)
			log.Info().
				Int("number", row.Number).
				Int("value", row.Value).
				Msg("Found value in database")
			resp := numberResponse{row.Value}
			return c.JSON(http.StatusOK, resp)
		}

		resp := numberResponse{json.Number + 1}
		return c.JSON(http.StatusOK, resp)
	})
}
