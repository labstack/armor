package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	BodyLimit struct {
		Base                       `json:",squash"`
		middleware.BodyLimitConfig `json:",squash"`
	}
)

func (l *BodyLimit) Initialize() (err error) {
	l.Middleware = middleware.BodyLimitWithConfig(l.BodyLimitConfig)
	return
}

func (*BodyLimit) Priority() int {
	return 1
}

func (l *BodyLimit) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return l.Middleware(next)
}

func (*BodyLimit) Terminate() {
}
