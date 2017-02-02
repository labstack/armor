package plugin

import (
	"fmt"
	"net/http"

	"github.com/labstack/cubex/middleware/go/cube"
	"github.com/labstack/echo"
)

type (
	Cube struct {
		cube.Cube `json:",squash"`
		Base      `json:",squash"`
		Path      string `json:"path"`
		Key       string `json:"key"`
	}
)

func (c *Cube) Init() (err error) {
	c.Cube = *cube.New(c.Echo.Server)
	fmt.Println(c.Echo)
	c.Skipper = func(r *http.Request) bool {
		return r.URL.Path == c.Path
	}
	c.Echo.GET(c.Path, func(ctx echo.Context) error {
		if ctx.Request().Header.Get("X-Cube-Key") != c.Key {
			return echo.ErrUnauthorized
		}
		return ctx.JSON(http.StatusOK, c.Data())
	})
	return
}

func (*Cube) Priority() int {
	return 1
}

func (c *Cube) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return c.Cube.Middleware(next)
}
