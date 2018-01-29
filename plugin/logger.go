package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Logger struct {
		Base                    `yaml:",squash"`
		middleware.LoggerConfig `yaml:",squash"`
	}
)

func (l *Logger) Initialize() error {
	l.Middleware = middleware.LoggerWithConfig(l.LoggerConfig)
	return nil
}

func (l *Logger) Update(p Plugin) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.LoggerConfig = p.(*Logger).LoggerConfig
	l.Initialize()
}

func (*Logger) Priority() int {
	return -1
}

func (l *Logger) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return l.Middleware(next)
}
