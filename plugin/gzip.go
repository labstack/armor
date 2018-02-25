package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Gzip struct {
		Base                  `yaml:",squash"`
		middleware.GzipConfig `yaml:",squash"`
	}
)

func (g *Gzip) Initialize() {
	g.Middleware = middleware.GzipWithConfig(g.GzipConfig)
}

func (g *Gzip) Update(p Plugin) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.GzipConfig = p.(*Gzip).GzipConfig
	g.Initialize()
}

func (*Gzip) Priority() int {
	return 1
}

func (g *Gzip) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return g.Middleware(next)
}
