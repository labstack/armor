package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	BodyLimit struct {
		Base                       `yaml:",squash"`
		middleware.BodyLimitConfig `yaml:",squash"`
	}
)

func (b *BodyLimit) Initialize() error {
	b.Middleware = middleware.BodyLimitWithConfig(b.BodyLimitConfig)
	return nil
}

func (b *BodyLimit) Update(p Plugin) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.BodyLimitConfig = p.(*BodyLimit).BodyLimitConfig
	b.Initialize()
}

func (*BodyLimit) Priority() int {
	return 1
}

func (b *BodyLimit) Process(next echo.HandlerFunc) echo.HandlerFunc {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.Middleware(next)
}
