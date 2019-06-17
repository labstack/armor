package plugin

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func (s *Secure) Process(next echo.HandlerFunc) echo.HandlerFunc {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Middleware(next)
}
