package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Static struct {
		Base                    `json:",squash"`
		middleware.StaticConfig `json:",squash"`
	}
)

func (s *Static) Initialize() (err error) {
	s.Prefix = s.Path
	s.Middleware = middleware.StaticWithConfig(s.StaticConfig)
	return
}

func (*Static) Priority() int {
	return 1
}

func (s *Static) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return s.Middleware(next)
}

func (*Static) Terminate() {
}
