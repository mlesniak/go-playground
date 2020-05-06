package context

import "github.com/labstack/echo"

// CustomContext is our application specific context.
type CustomContext struct {
	echo.Context
}

func (c *CustomContext) Foo() {
	println("foo")
}

func (c *CustomContext) Bar() {
	println("bar")
}