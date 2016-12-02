package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Logger struct {
		Base                    `json:",squash"`
		middleware.LoggerConfig `json:",squash"`
	}
)

func (l *Logger) Init() (err error) {
	l.Middleware = middleware.LoggerWithConfig(l.LoggerConfig)
	return
}

func (*Logger) Priority() int {
	return 1
}

func (l *Logger) Execute(next echo.HandlerFunc) echo.HandlerFunc {
	return l.Middleware(next)
}
