package plugin

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func (r *Rewrite) Process(next echo.HandlerFunc) echo.HandlerFunc {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.Middleware(next)
}
