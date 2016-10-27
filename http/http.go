package http

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/labstack/armor"
	"github.com/labstack/armor/plugin"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type (
	HTTP struct {
		logger *log.Logger
	}
)

func Start(a *armor.Armor) {
	h := &HTTP{
		logger: a.Logger,
	}
	e := echo.New()
	e.Logger = h.logger

	// Internal
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, "armor/"+armor.Version)
			return next(c)
		}
	})

	// Global plugins
	for name, pg := range a.Plugins {
		p, err := plugin.Decode(name, pg, "", "", a)
		if err != nil {
			h.logger.Error(err)
		}
		if p.Priority() < 0 {
			e.Pre(p.Process)
		} else {
			e.Use(p.Process)
		}
	}

	// Hosts
	for hn, host := range a.Hosts {
		host.Echo = echo.New()
		for name, pg := range host.Plugins {
			p, err := plugin.Decode(name, pg, hn, "", a)
			if err != nil {
				h.logger.Error(err)
			}
			if p.Priority() < 0 {
				host.Echo.Pre(p.Process)
			} else {
				host.Echo.Use(p.Process)
			}
		}

		// Paths
		for pn, path := range host.Paths {
			if pn == "/" {
				pn = ""
			}
			g := host.Echo.Group(pn)
			for name, pg := range path.Plugins {
				p, err := plugin.Decode(name, pg, hn, pn, a)
				if err != nil {
					h.logger.Error(err)
				}
				g.Use(p.Process)
			}
		}
	}

	// Route all requests
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := a.Hosts[req.Host]
		if host == nil {
			return echo.ErrNotFound
		}
		host.Echo.ServeHTTP(res, req)
		return
	})

	e.Server = &http.Server{
		ReadTimeout:  a.ReadTimeout * time.Second,
		WriteTimeout: a.WriteTimeout * time.Second,
		Handler:      e,
	}

	if a.TLS != nil {
		go func() {
			if err := h.startTLS(a, e); err != nil {
				h.logger.Fatal(err)
			}
		}()
	}
	if err := e.Start(a.Address); err != nil {
		h.logger.Fatal(err)
	}
}

func (h *HTTP) startTLS(a *armor.Armor, e *echo.Echo) error {
	if a.TLS.Auto {
		e.TLSCacheFile = a.TLS.CacheFile
		for host := range a.Hosts {
			e.TLSHosts = append(e.TLSHosts, host)
		}
		return e.StartAutoTLS()
	}
	a.TLS.Certificates = make(map[string]*tls.Certificate)
	e.TLSConfig.GetCertificate = func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		if a.TLS.Auto {
		}
		return a.TLS.Certificates[clientHello.ServerName], nil
	}
	for name, host := range a.Hosts {
		if host.CertFile == "" || host.KeyFile == "" {
			continue
		}
		cert, err := tls.LoadX509KeyPair(host.CertFile, host.KeyFile)
		if err != nil {
			h.logger.Fatal(err)
		}
		a.TLS.Certificates[name] = &cert
	}
	return e.StartTLS(a.TLS.Address, a.TLS.CertFile, a.TLS.KeyFile)
}
