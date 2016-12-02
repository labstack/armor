package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	AddTrailingSlash struct {
		Base                           `json:",squash"`
		middleware.TrailingSlashConfig `json:",squash"`
	}

	RemoveTrailingSlash struct {
		Base                           `json:",squash"`
		middleware.TrailingSlashConfig `json:",squash"`
	}
)

func (s *AddTrailingSlash) Init() (err error) {
	s.Middleware = middleware.AddTrailingSlashWithConfig(s.TrailingSlashConfig)
	return
}

func (*AddTrailingSlash) Priority() int {
	return -1
}

func (s *AddTrailingSlash) Execute(next echo.HandlerFunc) echo.HandlerFunc {
	return s.Middleware(next)
}

func (s *RemoveTrailingSlash) Init() (err error) {
	s.Middleware = middleware.RemoveTrailingSlashWithConfig(s.TrailingSlashConfig)
	return
}

func (*RemoveTrailingSlash) Priority() int {
	return -1
}

func (s *RemoveTrailingSlash) Execute(next echo.HandlerFunc) echo.HandlerFunc {
	return s.Middleware(next)
}
