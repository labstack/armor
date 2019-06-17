package plugin

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	CORS struct {
		Base                  `yaml:",squash"`
		middleware.CORSConfig `yaml:",squash"`
	}
)

func (c *CORS) Initialize() {
	c.Middleware = middleware.CORSWithConfig(c.CORSConfig)
}

func (c *CORS) Update(p Plugin) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.CORSConfig = p.(*CORS).CORSConfig
	c.Initialize()
}

func (c *CORS) Process(next echo.HandlerFunc) echo.HandlerFunc {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Middleware(next)
}
