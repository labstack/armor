package plugin

import (
	"fmt"
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
		Name string `json:"name,omitempty"`
		URL  string `json:"url"`
	}
)

func (p *Proxy) Init() (err error) {
	// Targets
	targets := make([]*middleware.ProxyTarget, len(p.Targets))
	for i, t := range p.Targets {
		u, err := url.Parse(t.URL)
		if err != nil {
			return fmt.Errorf("not able to parse proxy url=%s, error=%v", t.URL, err)
		}
		targets[i] = &middleware.ProxyTarget{
			URL: u,
		}
	}

	// Balancer
	switch p.Balance {
	case "round-robin":
		p.Balancer = &middleware.RoundRobinBalancer{Targets: targets}
	default: // Random
		p.Balancer = &middleware.RandomBalancer{Targets: targets}
	}

	// Need to be initialied in the end to reflect config changes.
	p.Middleware = middleware.Proxy(p.ProxyConfig)
	return
}

func (*Proxy) Priority() int {
	return 1
}

func (p *Proxy) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return p.Middleware(next)
}
