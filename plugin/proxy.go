package plugin

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Proxy struct {
		Base                   `json:",squash"`
		middleware.ProxyConfig `json:",squash"`
		Balance                string    `json:"balance"`
		Targets                []*Target `json:"targets"`
	}

	Target struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
)

func (t Target) ProxyTarget() (target *middleware.ProxyTarget, err error) {
	u, err := url.Parse(t.URL)
	if err != nil {
		return nil, fmt.Errorf("not able to parse proxy url=%s, error=%v", t.URL, err)
	}
	return &middleware.ProxyTarget{
		Name: t.Name,
		URL:  u,
	}, nil
}

func (p *Proxy) Init() (err error) {
	// Targets
	targets := make([]*middleware.ProxyTarget, len(p.Targets))
	for i, t := range p.Targets {
		pg, err := t.ProxyTarget()
		if err != nil {
			return err
		}
		targets[i] = pg
	}

	// Balancer
	switch p.Balance {
	case "round-robin":
		p.Balancer = middleware.NewRoundRobinBalancer(targets)
	default: // Random
		p.Balancer = middleware.NewRandomBalancer(targets)
	}

	// Need to be initialied in the end to reflect config changes.
	p.Middleware = middleware.ProxyWithConfig(p.ProxyConfig)

	return
}

func (*Proxy) Priority() int {
	return 1
}

func (p *Proxy) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return p.Middleware(next)
}

func (p *Proxy) AddTarget(c echo.Context) (err error) {
	t := new(Target)
	if err = c.Bind(t); err != nil {
		return
	}
	pt, err := t.ProxyTarget()
	if err != nil {
		return
	}
	p.Balancer.AddTarget(pt)
	return c.NoContent(http.StatusOK)
}

func (p *Proxy) RemoveTarget(c echo.Context) error {
	if p.Balancer.RemoveTarget(c.Param("target")) {
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusNotFound)
}
