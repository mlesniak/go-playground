package context

import (
	"github.com/labstack/echo/v4"
	"github.com/mlesniak/go-demo/pkg/authentication"
)

// CustomContext is our application specific context.
type CustomContext struct {
	echo.Context

	Authentication *authentication.Authentication
}

// Usernama returns the Username of the current logged in user.
func (c *CustomContext) Username() string {
	// TODO Handle the case that we are unauthorized yet.
	return c.Authentication.Username
}
