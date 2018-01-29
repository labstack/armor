package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Static struct {
		Base                    `yaml:",squash"`
		middleware.StaticConfig `yaml:",squash"`
	}
)

func (s *Static) Initialize() error {
	s.Middleware = middleware.StaticWithConfig(s.StaticConfig)
	return nil
}

func (s *Static) Update(p Plugin) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.StaticConfig = p.(*Static).StaticConfig
	s.Initialize()
}

func (*Static) Priority() int {
	return 1
}

func (s *Static) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return s.Middleware(next)
}
