package plugin

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	Proxy struct {
		Base                   `yaml:",squash"`
		middleware.ProxyConfig `yaml:",squash"`
		Balance                string    `yaml:"balance"`
		Targets                []*Target `yaml:"targets"`
	}

	Target struct {
		Name string `yaml:"name"`
		URL  string `yaml:"url"`
	}
)

func (t Target) ProxyTarget() (target *middleware.ProxyTarget, err error) {
	u, err := url.Parse(t.URL)
	if err != nil {
		return nil, fmt.Errorf("not able to parse proxy: url=%s, error=%v", t.URL, err)
	}
	return &middleware.ProxyTarget{
		Name: t.Name,
		URL:  u,
	}, nil
}

func (p *Proxy) Initialize() {
	// Targets
	targets := make([]*middleware.ProxyTarget, len(p.Targets))
	for i, t := range p.Targets {
		pg, err := t.ProxyTarget()
		if err != nil {
			panic(err)
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
}

func (p *Proxy) Update(plugin Plugin) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.ProxyConfig = plugin.(*Proxy).ProxyConfig
	p.Initialize()
}

func (p *Proxy) Process(next echo.HandlerFunc) echo.HandlerFunc {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
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
