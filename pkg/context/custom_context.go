package context

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

// injectRequestID adds request-specific information to the log.
func (c *CustomContext) injectRequestID(prev zerolog.Logger) zerolog.Logger {
	id := c.Response().Header().Get(echo.HeaderXRequestID)
	return prev.With().Str("requestId", id).Logger()
}

// injectRequestID adds request-specific information to the log.
func (c *CustomContext) injectUsername(prev zerolog.Logger) zerolog.Logger {
	if c.Authentication.Username != "" {
		return prev.With().Str("username", c.Authentication.Username).Logger()
	}

	return prev
}

// Log returns a logger with additional log fields based on the data from the context.
func (c *CustomContext) Log() zerolog.Logger {
	l := log.Logger
	l = c.injectRequestID(l)
	l = c.injectUsername(l)
	return l
}