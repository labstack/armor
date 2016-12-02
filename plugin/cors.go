package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	CORS struct {
		Base                  `json:",squash"`
		middleware.CORSConfig `json:",squash"`
	}
)

func (c *CORS) Init() (err error) {
	c.Middleware = middleware.CORSWithConfig(c.CORSConfig)
	return
}

func (*CORS) Priority() int {
	return 1
}

func (c *CORS) Execute(next echo.HandlerFunc) echo.HandlerFunc {
	return c.Middleware(next)
}
