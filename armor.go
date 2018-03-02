package armor

import (
	"sync"
	"time"

	"github.com/hashicorp/serf/serf"

	"github.com/labstack/armor/plugin"
	"github.com/labstack/armor/store"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
)

type (
	Armor struct {
		mutex        sync.RWMutex
		Name         string             `yaml:"name"`
		Address      string             `yaml:"address"`
		TLS          *TLS               `yaml:"tls"`
		Admin        *Admin             `yaml:"admin"`
		Postgres     *Postgres          `yaml:"postgres"`
		SQLite       *SQLite            `yaml:"sqlite"`
		Cluster      *Cluster           `yaml:"cluster"`
		ReadTimeout  time.Duration      `yaml:"read_timeout"`
		WriteTimeout time.Duration      `yaml:"write_timeout"`
		RawPlugins   []plugin.RawPlugin `yaml:"plugins"`
		Hosts        Hosts              `yaml:"hosts"`
		Store        store.Store        `yaml:"-"`
		Plugins      []plugin.Plugin    `yaml:"-"`
		Echo         *echo.Echo         `yaml:"-"`
		Logger       *log.Logger        `yaml:"-"`
		Colorer      *color.Color       `yaml:"-"`
	}

	TLS struct {
		Address      string `yaml:"address"`
		CertFile     string `yaml:"cert_file"`
		KeyFile      string `yaml:"key_file"`
		Auto         bool   `yaml:"auto"`
		CacheDir     string `yaml:"cache_dir"`
		Email        string `yaml:"email"`
		DirectoryURL string `yaml:"directory_url"`
	}

	Admin struct {
		Address string `yaml:"address"`
	}

	Postgres struct {
		URI string `yaml:"uri"`
	}

	SQLite struct {
		URI string `yaml:"uri"`
	}

	Cluster struct {
		*serf.Serf
		Address string   `yaml:"address"`
		Peers   []string `yaml:"peers"`
	}

	Host struct {
		mutex      sync.RWMutex
		Name       string             `yaml:"-"`
		CertFile   string             `yaml:"cert_file"`
		KeyFile    string             `yaml:"key_file"`
		RawPlugins []plugin.RawPlugin `yaml:"plugins"`
		Paths      Paths              `yaml:"paths"`
		Plugins    []plugin.Plugin    `yaml:"-"`
		Echo       *echo.Echo         `yaml:"-"`
	}

	Path struct {
		mutex      sync.RWMutex
		Name       string             `yaml:"-"`
		RawPlugins []plugin.RawPlugin `yaml:"plugins"`
		Plugins    []plugin.Plugin    `yaml:"-"`
		Group      *echo.Group        `yaml:"-"`
	}

	Hosts map[string]*Host

	Paths map[string]*Path
)

const (
	Version = "0.4.0-dev"
	Website = "https://armor.labstack.com"
)

func (a *Armor) AddHost(name string) *Host {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	h, ok := a.Hosts[name]
	if !ok {
		h = &Host{Paths: make(Paths)}
		a.Hosts[name] = h
	}
	h.Echo = echo.New()
	h.Name = name
	h.Echo.Logger = a.Logger
	return h
}

func (a *Armor) FindHost(name string) *Host {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.Hosts[name]
}

func (a *Armor) AddPlugin(p plugin.Plugin) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if p.Priority() < 0 {
		a.Echo.Pre(p.Process)
	} else {
		a.Echo.Use(p.Process)
	}
	a.Plugins = append(a.Plugins, p)
}

func (a *Armor) UpdatePlugin(plugin plugin.Plugin) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	for _, p := range a.Plugins {
		if p.Name() == plugin.Name() {
			p.Update(plugin)
		}
	}
}

func (a *Armor) LoadPlugin(p *store.Plugin, update bool) {
	if p.Host == "" && p.Path == "" {
		// Global level
		p := plugin.Decode(p.Raw, a.Echo, a.Logger)
		p.Initialize()
		if update {
			a.UpdatePlugin(p)
		} else {
			a.AddPlugin(p)
		}
	} else if p.Host != "" && p.Path == "" {
		// Host level
		host := a.FindHost(p.Host)
		if host == nil {
			host = a.AddHost(p.Host)
		}
		p := plugin.Decode(p.Raw, host.Echo, a.Logger)
		p.Initialize()
		if update {
			host.UpdatePlugin(p)
		} else {
			host.AddPlugin(p)
		}
	} else if p.Host != "" && p.Path != "" {
		// Path level
		host := a.FindHost(p.Host)
		if host == nil {
			host = a.AddHost(p.Host)
		}
		path := host.FindPath(p.Path)
		if path == nil {
			path = host.AddPath(p.Path)
		}
		p := plugin.Decode(p.Raw, host.Echo, a.Logger)
		p.Initialize()
		if update {
			path.UpdatePlugin(p)
		} else {
			path.AddPlugin(p)
		}
	}
}

func (h *Host) AddPath(name string) *Path {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	p, ok := h.Paths[name]
	if !ok {
		p = new(Path)
		h.Paths[name] = p
	}
	p.Name = name
	p.Group = h.Echo.Group(name)
	return p
}

func (h *Host) FindPath(name string) *Path {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.Paths[name]
}

func (h *Host) AddPlugin(p plugin.Plugin) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if p.Priority() < 0 {
		h.Echo.Pre(p.Process)
	} else {
		h.Echo.Use(p.Process)
	}
	h.Plugins = append(h.Plugins, p)
}

func (h *Host) UpdatePlugin(plugin plugin.Plugin) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	for _, p := range h.Plugins {
		if p.Name() == plugin.Name() {
			p.Update(plugin)
		}
	}
}

func (p *Path) AddPlugin(plugin plugin.Plugin) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Group.Use(plugin.Process)
	p.Plugins = append(p.Plugins, plugin)
}

func (p *Path) UpdatePlugin(plugin plugin.Plugin) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	for _, p := range p.Plugins {
		if p.Name() == plugin.Name() {
			p.Update(plugin)
		}
	}
}
