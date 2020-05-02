package log

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// AddUser adds a user field to log if it was pre-filled from the JWT middleware.
func AddUser(log *logrus.Entry, c echo.Context) *logrus.Entry {
	u := c.Get("user")
	if u != nil {
		token := u.(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		return log.WithField("user", claims["user"])
	}

	return log
}
