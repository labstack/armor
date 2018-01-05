package armor

import (
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
)

type (
	Armor struct {
		mutex        sync.RWMutex
		Address      string        `yaml:"address"`
		TLS          *TLS          `yaml:"tls"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
		RawPlugins   []RawPlugin   `yaml:"plugins"`
		Hosts        Hosts         `yaml:"hosts"`
		Plugins      []Plugin      `yaml:"-"`
		Echo         *echo.Echo    `yaml:"-"`
		Logger       *log.Logger   `yaml:"-"`
		Colorer      *color.Color  `yaml:"-"`
	}

	TLS struct {
		Address  string `yaml:"address"`
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
		Auto     bool   `yaml:"auto"`
		CacheDir string `yaml:"cache_dir"`
	}

	Host struct {
		mutex      sync.RWMutex
		Name       string      `yaml:"-"`
		CertFile   string      `yaml:"cert_file"`
		KeyFile    string      `yaml:"key_file"`
		RawPlugins []RawPlugin `yaml:"plugins"`
		Paths      Paths       `yaml:"paths"`
		Plugins    []Plugin    `yaml:"-"`
		Echo       *echo.Echo  `yaml:"-"`
	}

	Path struct {
		mutex      sync.RWMutex
		Name       string      `yaml:"-"`
		RawPlugins []RawPlugin `yaml:"plugins"`
		Plugins    []Plugin    `yaml:"-"`
		Group      *echo.Group `yaml:"-"`
	}

	Hosts map[string]*Host

	Paths map[string]*Path

	RawPlugin map[string]interface{}

	Plugin interface {
		Name() string
		Init() error
		Process(echo.HandlerFunc) echo.HandlerFunc
		Priority() int
	}
)

const (
	Version = "0.3.7"
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

func (a *Armor) AddPlugin(p Plugin) {
	if p.Priority() < 0 {
		a.Echo.Pre(p.Process)
	} else {
		a.Echo.Use(p.Process)
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.Plugins = append(a.Plugins, p)
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

func (h *Host) AddPlugin(p Plugin) {
	if p.Priority() < 0 {
		h.Echo.Pre(p.Process)
	} else {
		h.Echo.Use(p.Process)
	}
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.Plugins = append(h.Plugins, p)
}

func (p *Path) AddPlugin(plugin Plugin) {
	p.Group.Use(plugin.Process)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Plugins = append(p.Plugins, plugin)
}
