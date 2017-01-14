package http

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/labstack/armor"
	"github.com/labstack/armor/plugin"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/acme/autocert"
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
	e.HideBanner = true

	// Internal
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, "armor/"+armor.Version)
			return next(c)
		}
	})

	// Global plugins
	for _, pi := range a.Plugins {
		p, err := plugin.Decode(pi, a)
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
		host.Name = hn
		host.Echo = echo.New()

		for _, pi := range host.Plugins {
			p, err := plugin.Decode(pi, a)
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
			g := host.Echo.Group(pn)

			for _, pi := range path.Plugins {
				p, err := plugin.Decode(pi, a)
				if err != nil {
					h.logger.Error(err)
				}
				g.Use(p.Process)
			}

			// NOP handlers to trigger plugins
			g.Any("", echo.NotFoundHandler)
			if pn == "/" {
				g.Any("*", echo.NotFoundHandler)
			} else {
				g.Any("/*", echo.NotFoundHandler)
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

	if a.TLS != nil {
		go h.logger.Fatal(h.startTLS(a, e))
	}
	h.logger.Fatal(h.start(a, e))
}

func (h *HTTP) startTLS(a *armor.Armor, e *echo.Echo) error {
	s := &http.Server{
		Addr:         a.TLS.Address,
		TLSConfig:    new(tls.Config),
		ReadTimeout:  a.ReadTimeout * time.Second,
		WriteTimeout: a.WriteTimeout * time.Second,
	}

	if a.TLS.Auto {
		hosts := []string{}
		for host := range a.Hosts {
			hosts = append(hosts, host)
		}
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(hosts...) // Added security
		e.AutoTLSManager.Cache = autocert.DirCache(a.TLS.CacheDir)
	}

	for name, host := range a.Hosts {
		if host.CertFile == "" || host.KeyFile == "" {
			continue
		}
		cert, err := tls.LoadX509KeyPair(host.CertFile, host.KeyFile)
		if err != nil {
			h.logger.Fatal(err)
		}
		s.TLSConfig.NameToCertificate[name] = &cert
	}

	s.TLSConfig.GetCertificate = func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		if cert, ok := s.TLSConfig.NameToCertificate[clientHello.ServerName]; ok {
			// Use provided certificate
			return cert, nil
		} else if a.TLS.Auto {
			return e.AutoTLSManager.GetCertificate(clientHello)
		}
		return nil, nil // No certificate
	}

	return e.StartServer(s)
}

func (h *HTTP) start(a *armor.Armor, e *echo.Echo) error {
	s := &http.Server{
		Addr:         a.Address,
		ReadTimeout:  a.ReadTimeout * time.Second,
		WriteTimeout: a.WriteTimeout * time.Second,
	}
	args := []interface{}{a.Colorer.Blue(armor.Website), a.Colorer.Red("v" + armor.Version), "http", a.Colorer.Green(s.Addr)}
	if a.TLS != nil {
		args[2] = "https"
	}
	a.Colorer.Printf(armor.Banner, args...)
	return e.StartServer(s)
}
