package plugin

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	RedirectConfig struct {
		From string `yaml:"from"`
		To   string `yaml:"to"`
		Code int    `yaml:"code"`
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
	t := NewTemplate(r.To)
	// Defaults
	if r.Code == 0 {
		r.Code = http.StatusMovedPermanently
	}
	r.Echo.GET(r.From, func(c echo.Context) error {
		to, err := t.Execute(c)
		if err != nil {
			return err
		}
		return c.Redirect(r.Code, to)
	})
}

func (r *Redirect) Update(p Plugin) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.RedirectConfig = p.(*Redirect).RedirectConfig
	r.Initialize()
}

func (*Redirect) Priority() int {
	return -1
}

func (r *Redirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
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

func (*HTTPSRedirect) Priority() int {
	return -1
}

func (r *HTTPSRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
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

func (*HTTPSWWWRedirect) Priority() int {
	return -1
}

func (r *HTTPSWWWRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
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

func (*HTTPSNonWWWRedirect) Priority() int {
	return -1
}

func (r *HTTPSNonWWWRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
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

func (*WWWRedirect) Priority() int {
	return -1
}

func (r *WWWRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
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

func (*NonWWWRedirect) Priority() int {
	return -1
}

func (r *NonWWWRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return r.Middleware(next)
}
