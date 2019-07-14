package plugin

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type (
	RedirectConfig struct {
		template *Template
		From     string `yaml:"from"`
		To       string `yaml:"to"`
		Code     int    `yaml:"code"`
	}

	Redirect struct {
		Base           `yaml:",squash"`
		RedirectConfig `yaml:",squash"`
	}

	HTTPSRedirect struct {
		Base                      `yaml:",squash"`
		middleware.RedirectConfig `yaml:",squash"`
	}

	HTTPSWWWRedirect struct {
		Base                      `yaml:",squash"`
		middleware.RedirectConfig `yaml:",squash"`
	}

	HTTPSNonWWWRedirect struct {
		Base                      `yaml:",squash"`
		middleware.RedirectConfig `yaml:",squash"`
	}

	WWWRedirect struct {
		Base                      `yaml:",squash"`
		middleware.RedirectConfig `yaml:",squash"`
	}

	NonWWWRedirect struct {
		Base                      `yaml:",squash"`
		middleware.RedirectConfig `yaml:",squash"`
	}
)

func (r *Redirect) Initialize() {
	r.template = NewTemplate(r.To)
	// Defaults
	if r.Code == 0 {
		r.Code = http.StatusMovedPermanently
	}
}

func (r *Redirect) Update(p Plugin) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.RedirectConfig = p.(*Redirect).RedirectConfig
	r.Initialize()
}

func (r *Redirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return func(c echo.Context) error {
		if c.Request().URL.Path == r.From {
			to, err := r.template.Execute(c)
			if err != nil {
				return err
			}
			return c.Redirect(r.Code, to)
		}
		return next(c)
	}
}

func (r *HTTPSRedirect) Initialize() {
	r.Middleware = middleware.HTTPSRedirectWithConfig(r.RedirectConfig)
}

func (r *HTTPSRedirect) Update(p Plugin) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.RedirectConfig = p.(*HTTPSRedirect).RedirectConfig
	r.Initialize()
}

func (r *HTTPSRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.Middleware(next)
}

func (r *HTTPSWWWRedirect) Initialize() {
	r.Middleware = middleware.HTTPSWWWRedirectWithConfig(r.RedirectConfig)
}

func (r *HTTPSWWWRedirect) Update(p Plugin) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.RedirectConfig = p.(*HTTPSWWWRedirect).RedirectConfig
	r.Initialize()
}

func (r *HTTPSWWWRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.Middleware(next)
}

func (r *HTTPSNonWWWRedirect) Initialize() {
	e := NewExpression(r.Skip)
	r.RedirectConfig.Skipper = func(c echo.Context) bool {
		skip, err := e.Evaluate(c)
		if err != nil {
			return false
		}
		return skip.(bool)
	}
	r.Middleware = middleware.HTTPSNonWWWRedirectWithConfig(r.RedirectConfig)
}

func (r *HTTPSNonWWWRedirect) Update(p Plugin) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.RedirectConfig = p.(*HTTPSNonWWWRedirect).RedirectConfig
	r.Initialize()
}

func (r *HTTPSNonWWWRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.Middleware(next)
}

func (r *WWWRedirect) Initialize() {
	r.Middleware = middleware.WWWRedirectWithConfig(r.RedirectConfig)
}

func (r *WWWRedirect) Update(p Plugin) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.RedirectConfig = p.(*WWWRedirect).RedirectConfig
	r.Initialize()
}

func (r *WWWRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.Middleware(next)
}

func (r *NonWWWRedirect) Initialize() {
	r.Middleware = middleware.NonWWWRedirectWithConfig(r.RedirectConfig)
}

func (r *NonWWWRedirect) Update(p Plugin) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.RedirectConfig = p.(*NonWWWRedirect).RedirectConfig
	r.Initialize()
}

func (r *NonWWWRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.Middleware(next)
}
