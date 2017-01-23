package http

import (
	"crypto/tls"
	"net/http"
	"path/filepath"
	"time"

	"github.com/labstack/armor"
	"github.com/labstack/armor/plugin"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/acme/autocert"
)

type (
	HTTP struct {
		armor  *armor.Armor
		echo   *echo.Echo
		logger *log.Logger
	}
)

func Init(a *armor.Armor) (h *HTTP) {
	h = &HTTP{
		armor:  a,
		echo:   echo.New(),
		logger: a.Logger,
	}
	e := h.echo
	e.Logger = h.logger

	// Internal
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, "armor/"+armor.Version)
			return next(c)
		}
	})

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

	return
}

func (h *HTTP) Start() error {
	a := h.armor
	e := h.echo
	s := &http.Server{
		Addr:         a.Address,
		ReadTimeout:  a.ReadTimeout * time.Second,
		WriteTimeout: a.WriteTimeout * time.Second,
	}
	return e.StartServer(s)
}

func (h *HTTP) StartTLS() error {
	a := h.armor
	e := h.echo
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
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		if a.TLS.CacheDir == "" {
			a.TLS.CacheDir = filepath.Join(home, ".armor", "cache")
		}
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

func (h *HTTP) LoadPlugins() {
	a := h.armor
	e := h.echo

	// Global plugins
	for _, pi := range a.Plugins {
		p, err := plugin.Decode(pi, a, e)
		if err != nil {
			h.logger.Fatal(err)
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
		host.Echo.Logger = a.Logger

		// Host plugins
		for _, pi := range host.Plugins {
			p, err := plugin.Decode(pi, a, host.Echo)
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

			// Path plugins
			for _, pi := range path.Plugins {
				p, err := plugin.Decode(pi, a, host.Echo)
				if err != nil {
					h.logger.Fatal(err)
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
}
