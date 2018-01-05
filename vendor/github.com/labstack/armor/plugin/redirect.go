package plugin

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Redirect struct {
		Base `yaml:",squash"`
		From string `yaml:"from"`
		To   string `yaml:"to"`
		Code int    `yaml:"code"`
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

func (r *Redirect) Init() (err error) {
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
	return
}

func (*Redirect) Priority() int {
	return -1
}

func (r *Redirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func (r *HTTPSRedirect) Init() (err error) {
	r.Middleware = middleware.HTTPSRedirectWithConfig(r.RedirectConfig)
	return
}

func (*HTTPSRedirect) Priority() int {
	return -1
}

func (r *HTTPSRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return r.Middleware(next)
}

func (r *HTTPSWWWRedirect) Init() (err error) {
	r.Middleware = middleware.HTTPSWWWRedirectWithConfig(r.RedirectConfig)
	return
}

func (*HTTPSWWWRedirect) Priority() int {
	return -1
}

func (r *HTTPSWWWRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return r.Middleware(next)
}

func (r *HTTPSNonWWWRedirect) Init() (err error) {
	e := NewExpression(r.Skip)
	r.RedirectConfig.Skipper = func(c echo.Context) bool {
		skip, err := e.Evaluate(c)
		if err != nil {
			return false
		}
		return skip.(bool)
	}
	r.Middleware = middleware.HTTPSNonWWWRedirectWithConfig(r.RedirectConfig)
	return
}

func (*HTTPSNonWWWRedirect) Priority() int {
	return -1
}

func (r *HTTPSNonWWWRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return r.Middleware(next)
}

func (r *WWWRedirect) Init() (err error) {
	r.Middleware = middleware.WWWRedirectWithConfig(r.RedirectConfig)
	return
}

func (*WWWRedirect) Priority() int {
	return -1
}

func (r *WWWRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return r.Middleware(next)
}

func (r *NonWWWRedirect) Init() (err error) {
	r.Middleware = middleware.NonWWWRedirectWithConfig(r.RedirectConfig)
	return
}

func (*NonWWWRedirect) Priority() int {
	return -1
}

func (r *NonWWWRedirect) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return r.Middleware(next)
}
