package plugin

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	BodyLimit struct {
		Base                       `yaml:",squash"`
		middleware.BodyLimitConfig `yaml:",squash"`
	}
)

func (b *BodyLimit) Initialize() {
	b.Middleware = middleware.BodyLimitWithConfig(b.BodyLimitConfig)
}

func (b *BodyLimit) Update(p Plugin) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.BodyLimitConfig = p.(*BodyLimit).BodyLimitConfig
	b.Initialize()
}

func (b *BodyLimit) Process(next echo.HandlerFunc) echo.HandlerFunc {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.Middleware(next)
}
