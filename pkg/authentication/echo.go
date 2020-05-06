package authentication

import (
	"github.com/labstack/echo/v4"
	"github.com/mlesniak/go-demo/pkg/context"
	"github.com/sirupsen/logrus"
)

// AddUser adds a user field to log if it was pre-filled.
func AddUser(log *logrus.Entry, cc echo.Context) *logrus.Entry {
	c := cc.(*context.CustomContext)

	username := c.Username()
	if username != "" {
		return log.WithField("username", username)
	}

	return log
}
