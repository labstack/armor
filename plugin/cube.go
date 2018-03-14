package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/labstack-go/cube"
	labstack "github.com/labstack/labstack-go/echo"
)

type (
	Cube struct {
		Base       `yaml:",squash"`
		CubeConfig `yaml:",squash"`
	}

	CubeConfig struct {
		cube.Options
		APIKey string `yaml:"api_key"`
	}
)

func (c *Cube) Initialize() {
	c.Middleware = labstack.CubeWithOptions(c.APIKey, c.Options)
}

func (*Cube) Priority() int {
	return 1
}

func (c *Cube) Update(p Plugin) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.CubeConfig = p.(*Cube).CubeConfig
	c.Initialize()
}

func (c *Cube) Process(next echo.HandlerFunc) echo.HandlerFunc {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Middleware(next)
}
