package context

import "github.com/labstack/echo/v4"

// CreateCustomContext returns a default middleware which creates and injects the custom context.
func CreateCustomContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &CustomContext{
			Context: c,
			Authentication: Authentication{}}
		return next(cc)
	}
}
