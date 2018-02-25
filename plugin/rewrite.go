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

func (r *Rewrite) Initialize() {
	r.Middleware = middleware.RewriteWithConfig(r.RewriteConfig)
}

func (r *Rewrite) Update(p Plugin) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.RewriteConfig = p.(*Rewrite).RewriteConfig
	r.Initialize()
}

func (*Rewrite) Priority() int {
	return 1
}

func (r *Rewrite) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return r.Middleware(next)
}
