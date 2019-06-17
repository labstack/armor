package plugin

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	Static struct {
		Base                    `yaml:",squash"`
		middleware.StaticConfig `yaml:",squash"`
	}
)

func (s *Static) Initialize() {
	s.Middleware = middleware.StaticWithConfig(s.StaticConfig)
}

func (s *Static) Update(p Plugin) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.StaticConfig = p.(*Static).StaticConfig
	s.Initialize()
}

func (s *Static) Process(next echo.HandlerFunc) echo.HandlerFunc {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Middleware(next)
}
