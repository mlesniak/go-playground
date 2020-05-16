package version

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// AddVersionEndpoint adds the version endpoint
func AddVersionEndpoint(e *echo.Echo) {
	e.GET("/api/version", func(c echo.Context) error {
		log.Info().Msg("Version requested")
		commit := os.Getenv("COMMIT")
		if commit == "" {
			commit = "<No COMMIT environment variable set>"
		}
		return c.JSON(http.StatusOK, struct {
			Version string `json:"version"`
		}{commit})
	})
}