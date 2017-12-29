package plugin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Rewrite struct {
		Base                     `json:",squash"`
		middleware.RewriteConfig `json:",squash"`
	}
)

func (r *Rewrite) Init() (err error) {
	r.Middleware = middleware.RewriteWithConfig(r.RewriteConfig)
	return
}

func (*Rewrite) Priority() int {
	return 1
}

func (r *Rewrite) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return r.Middleware(next)
}
