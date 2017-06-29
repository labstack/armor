package plugin

import (
	"github.com/labstack/cube/echo"
	"github.com/labstack/echo"
)

type (
	Cube struct {
		Base        `json:",squash"`
		cube.Config `json:",squash"`
	}
)

func (c *Cube) Init() (err error) {
	c.Middleware = cube.MiddlewareWithConfig(c.Config)
	return
}

func (*Cube) Priority() int {
	return -1
}

func (c *Cube) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return c.Middleware(next)
}
