package context

import (
	"github.com/labstack/echo/v4"
)

// Authentication contains all user information of an authenticated user.
type Authentication struct {
	Username string
	Roles    []string
}


// CustomContext is our application specific context.
type CustomContext struct {
	echo.Context

	Authentication Authentication
}

// Usernama returns the Username of the current logged in user.
func (c *CustomContext) Username() string {
	// TODO Handle the case that we are unauthorized yet.
	return c.Authentication.Username
}
