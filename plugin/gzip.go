package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Gzip struct {
		Base                  `json:",squash"`
		middleware.GzipConfig `json:",squash"`
	}
)

func (g *Gzip) Init() (err error) {
	g.Middleware = middleware.GzipWithConfig(g.GzipConfig)
	return
}

func (*Gzip) Priority() int {
	return 1
}

func (g *Gzip) Execute(next echo.HandlerFunc) echo.HandlerFunc {
	return g.Middleware(next)
}
