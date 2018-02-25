package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Secure struct {
		Base                    `yaml:",squash"`
		middleware.SecureConfig `yaml:",squash"`
	}
)

func (s *Secure) Initialize() {
	s.Middleware = middleware.SecureWithConfig(s.SecureConfig)
}

func (s *Secure) Update(p Plugin) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.SecureConfig = p.(*Secure).SecureConfig
	s.Initialize()
}

func (*Secure) Priority() int {
	return 1
}

func (s *Secure) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return s.Middleware(next)
}
