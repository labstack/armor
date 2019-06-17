package plugin

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	Logger struct {
		Base                    `yaml:",squash"`
		middleware.LoggerConfig `yaml:",squash"`
	}
)

func (l *Logger) Initialize() {
	l.Middleware = middleware.LoggerWithConfig(l.LoggerConfig)
}

func (l *Logger) Update(p Plugin) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.LoggerConfig = p.(*Logger).LoggerConfig
	l.Initialize()
}

func (l *Logger) Process(next echo.HandlerFunc) echo.HandlerFunc {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.Middleware(next)
}
