package http

import (
	"crypto/tls"
	"net/http"
	"path/filepath"
	"time"

	"github.com/labstack/armor"
	"github.com/labstack/armor/plugin"
	"github.com/labstack/armor/util"
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
	e := echo.New()
	a.Echo = e
	h = &HTTP{
		armor:  a,
		echo:   e,
		logger: a.Logger,
	}
	e.HideBanner = true
	e.HidePort = true
	e.Server = &http.Server{
		Addr:         a.Address,
		ReadTimeout:  a.ReadTimeout * time.Second,
		WriteTimeout: a.WriteTimeout * time.Second,
	}
	if a.TLS != nil {
		e.TLSServer = &http.Server{
			Addr:         a.TLS.Address,
			TLSConfig:    new(tls.Config),
			ReadTimeout:  a.ReadTimeout * time.Second,
			WriteTimeout: a.WriteTimeout * time.Second,
		}
	}
	e.Logger = h.logger

	// Internal
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Before(func() {
				c.Response().Header().Set(echo.HeaderServer, "armor/"+armor.Version)
			})
			return next(c)
		}
	})

	// Route all requests
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := a.FindHost(util.StripPort(req.Host))
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
	a.Colorer.Printf("⇨ http server started on %s\n", a.Colorer.Green(a.Address))
	return e.StartServer(e.Server)
}

func (h *HTTP) StartTLS() error {
	a := h.armor
	e := h.echo
	s := e.TLSServer

	// Enable HTTP/2
	s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")

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

	// Load certificates - start
	// Global
	if a.TLS.CertFile != "" && a.TLS.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(a.TLS.CertFile, a.TLS.KeyFile)
		if err != nil {
			h.logger.Fatal(err)
		}
		s.TLSConfig.Certificates = append(s.TLSConfig.Certificates, cert)
	}
	// Host
	for _, host := range a.Hosts {
		if host.CertFile == "" || host.KeyFile == "" {
			continue
		}
		cert, err := tls.LoadX509KeyPair(host.CertFile, host.KeyFile)
		if err != nil {
			h.logger.Fatal(err)
		}
		s.TLSConfig.Certificates = append(s.TLSConfig.Certificates, cert)
	}
	s.TLSConfig.BuildNameToCertificate()
	// Load certificates - end

	s.TLSConfig.GetCertificate = func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		if cert, ok := s.TLSConfig.NameToCertificate[clientHello.ServerName]; ok {
			// Use provided certificate
			return cert, nil
		} else if a.TLS.Auto {
			return e.AutoTLSManager.GetCertificate(clientHello)
		}
		return nil, nil // No certificate
	}

	a.Colorer.Printf("⇨ https server started on %s\n", a.Colorer.Green(a.TLS.Address))
	return e.StartServer(s)
}

func (h *HTTP) LoadPlugins() {
	a := h.armor
	e := h.echo

	// Global plugins
	for _, rp := range a.RawPlugins {
		p, err := plugin.Decode(rp, a, e)
		if err != nil {
			h.logger.Fatal(err)
		}
		a.AddPlugin(p)
	}

	// Hosts
	for name, host := range a.Hosts {
		a.AddHost(name)

		// Host plugins
		for _, rp := range host.RawPlugins {
			p, err := plugin.Decode(rp, a, host.Echo)
			if err != nil {
				h.logger.Error(err)
			}
			host.AddPlugin(p)
		}

		// Paths
		for name, path := range host.Paths {
			host.AddPath(name)

			// Path plugins
			for _, rp := range path.RawPlugins {
				p, err := plugin.Decode(rp, a, host.Echo)
				if err != nil {
					h.logger.Fatal(err)
				}
				path.AddPlugin(p)
			}
		}
	}
}
