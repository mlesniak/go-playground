package authentication

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// AddUser adds a user field to log if it was pre-filled.
func AddUser(log *logrus.Entry, c echo.Context) *logrus.Entry {
	u := c.Get(ContextName)
	if t, ok := u.(Authentication); ok {
		return log.WithField("username", t.Username)
	}

	return log
}
