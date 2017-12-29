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

func (s *Secure) Init() (err error) {
	s.Middleware = middleware.SecureWithConfig(s.SecureConfig)
	return
}

func (*Secure) Priority() int {
	return 1
}

func (s *Secure) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return s.Middleware(next)
}
