package plugin

import (
	"net/http"

	"github.com/labstack/cubex/go/cube"
	"github.com/labstack/echo"
)

type (
	Cube struct {
		middleware *cube.Cube
		Base       `json:",squash"`
		Path       string `json:"path"`
		Key        string `json:"key"`
	}
)

func (c *Cube) Init() (err error) {
	s := c.Echo.Server
	if c.Armor.TLS != nil {
		s = c.Echo.TLSServer
	}
	m := cube.New(s)
	c.middleware = m
	c.Middleware = m.Middleware
	c.Echo.GET(c.Path, func(ctx echo.Context) error {
		if ctx.Request().Header.Get("X-Cube-Key") != c.Key {
			return echo.ErrUnauthorized
		}
		return ctx.JSON(http.StatusOK, m.Data())
	})
	return
}

func (*Cube) Priority() int {
	return -1
}

func (c *Cube) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return c.Middleware(next)
}
