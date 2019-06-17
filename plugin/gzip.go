package plugin

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func (g *Gzip) Process(next echo.HandlerFunc) echo.HandlerFunc {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.Middleware(next)
}
