package plugin

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	AddTrailingSlash struct {
		Base                           `yaml:",squash"`
		middleware.TrailingSlashConfig `yaml:",squash"`
	}

	RemoveTrailingSlash struct {
		Base                           `yaml:",squash"`
		middleware.TrailingSlashConfig `yaml:",squash"`
	}
)

func (s *AddTrailingSlash) Initialize() {
	s.Middleware = middleware.AddTrailingSlashWithConfig(s.TrailingSlashConfig)
}

func (s *AddTrailingSlash) Update(p Plugin) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TrailingSlashConfig = p.(*AddTrailingSlash).TrailingSlashConfig
	s.Initialize()
}

func (s *AddTrailingSlash) Process(next echo.HandlerFunc) echo.HandlerFunc {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Middleware(next)
}

func (s *RemoveTrailingSlash) Initialize() {
	s.Middleware = middleware.RemoveTrailingSlashWithConfig(s.TrailingSlashConfig)
}

func (s *RemoveTrailingSlash) Update(p Plugin) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TrailingSlashConfig = p.(*RemoveTrailingSlash).TrailingSlashConfig
	s.Initialize()
}

func (s *RemoveTrailingSlash) Process(next echo.HandlerFunc) echo.HandlerFunc {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Middleware(next)
}
