package version

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	logger "github.com/mlesniak/go-demo/pkg/log"
)

var log = logger.New()

// AddVersionEndpoint adds the version endpoint
func AddVersionEndpoint(e *echo.Echo) {
	e.GET("/api/version", func(c echo.Context) error {
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