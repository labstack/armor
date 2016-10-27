package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Secure struct {
		Base                    `json:",squash"`
		middleware.SecureConfig `json:",squash"`
	}
)

func (*Secure) Init() (err error) {
	return
}

func (*Secure) Priority() int {
	return 1
}

func (s *Secure) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return s.Middleware(next)
}
