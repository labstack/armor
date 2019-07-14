package armor

import (
	"crypto/tls"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/serf/serf"

	"github.com/labstack/armor/plugin"
	"github.com/labstack/armor/store"
	"github.com/labstack/armor/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
)

type (
	Armor struct {
		mutex         sync.RWMutex
		Name          string             `json:"name"`
		Address       string             `json:"address"`
		Port          string             `json:"-"`
		TLS           *TLS               `json:"tls"`
		Admin         *Admin             `json:"admin"`
		Storm         *Storm             `json:"storm"`
		Postgres      *Postgres          `json:"postgres"`
		Cluster       *Cluster           `json:"cluster"`
		ReadTimeout   time.Duration      `json:"read_timeout"`
		WriteTimeout  time.Duration      `json:"write_timeout"`
		RawPlugins    []plugin.RawPlugin `json:"plugins"`
		Hosts         Hosts              `json:"hosts"`
		RootDir       string             `json:"-"`
		Store         store.Store        `json:"-"`
		Plugins       []plugin.Plugin    `json:"-"`
		Echo          *echo.Echo         `json:"-"`
		Logger        *log.Logger        `json:"-"`
		Colorer       *color.Color       `json:"-"`
		DefaultConfig bool               `json:"-"`
	}

	TLS struct {
		Address      string `json:"address"`
		Port         string `json:"-"`
		CertFile     string `json:"cert_file"`
		KeyFile      string `json:"key_file"`
		Auto         bool   `json:"auto"`
		CacheDir     string `json:"cache_dir"`
		Email        string `json:"email"`
		DirectoryURL string `json:"directory_url"`
		Secured      bool   `json:"secured"`
	}

	Admin struct {
		Address string `json:"address"`
	}

	Storm struct {
		URI string `json:"uri"`
	}

	Postgres struct {
		URI string `json:"uri"`
	}

	Cluster struct {
		*serf.Serf
		Address string   `json:"address"`
		Peers   []string `json:"peers"`
	}

	Host struct {
		mutex       sync.RWMutex
		initialized bool
		Name        string             `json:"-"`
		CertFile    string             `json:"cert_file"`
		KeyFile     string             `json:"key_file"`
		RawPlugins  []plugin.RawPlugin `json:"plugins"`
		Paths       Paths              `json:"paths"`
		Plugins     []plugin.Plugin    `json:"-"`
		Group       *echo.Group        `json:"-"`
		ClientCAs   []string           `json:"client_ca"`
		TLSConfig   *tls.Config        `json:"-"`
	}

	Path struct {
		mutex       sync.RWMutex
		initialized bool
		Name        string             `json:"-"`
		RawPlugins  []plugin.RawPlugin `json:"plugins"`
		Plugins     []plugin.Plugin    `json:"-"`
		Group       *echo.Group        `json:"-"`
	}

	Hosts map[string]*Host

	Paths map[string]*Path
)

const (
	Version = "0.4.14"
	Website = "https://armor.labstack.com"
)

var (
	prePlugins = map[string]bool{
		plugin.PluginLogger:              true,
		plugin.PluginRedirect:            true,
		plugin.PluginHTTPSRedirect:       true,
		plugin.PluginHTTPSWWWRedirect:    true,
		plugin.PluginHTTPSNonWWWRedirect: true,
		plugin.PluginWWWRedirect:         true,
		plugin.PluginAddTrailingSlash:    true,
		plugin.PluginRemoveTrailingSlash: true,
		plugin.PluginNonWWWRedirect:      true,
		plugin.PluginRewrite:             true,
	}
)

func (a *Armor) FindHost(name string, add bool) (h *Host) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	h = a.Hosts[name]

	// Host lookup
	if h == nil && !add {
		return
	}

	// Add host
	if h == nil {
		h = new(Host)
		a.Hosts[name] = h
	}

	// Initialize host
	if !h.initialized {
		h.Name = name
		h.Paths = make(Paths)
		h.Group = a.Echo.Host(net.JoinHostPort(name, a.Port))
		routers := a.Echo.Routers()
		routers[net.JoinHostPort(name, a.TLS.Port)] = routers[name]
		h.initialized = true
	}

	return
}

func (a *Armor) AddPlugin(p plugin.Plugin) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if p.Order() < 0 {
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
		host := a.FindHost(p.Host, true)
		p := plugin.Decode(p.Raw, a.Echo, a.Logger)
		p.Initialize()
		if update {
			host.UpdatePlugin(p)
		} else {
			host.AddPlugin(p)
		}
	} else if p.Host != "" && p.Path != "" {
		// Path level
		host := a.FindHost(p.Host, true)
		path := host.FindPath(p.Path)
		p := plugin.Decode(p.Raw, a.Echo, a.Logger)
		p.Initialize()
		if update {
			path.UpdatePlugin(p)
		} else {
			path.AddPlugin(p)
		}
	}
}

func (a *Armor) SavePlugins() {
	plugins := []*store.Plugin{}

	// Global plugins
	for _, rp := range a.RawPlugins {
		plugins = append(plugins, &store.Plugin{
			Name:   rp.Name(),
			Config: rp.JSON(),
		})
	}

	for hn, host := range a.Hosts {
		// Host plugins
		for _, rp := range host.RawPlugins {
			plugins = append(plugins, &store.Plugin{
				Name:   rp.Name(),
				Host:   hn,
				Config: rp.JSON(),
			})
		}

		for pn, path := range host.Paths {
			// Path plugins
			for _, rp := range path.RawPlugins {
				plugins = append(plugins, &store.Plugin{
					Name:   rp.Name(),
					Host:   hn,
					Path:   pn,
					Config: rp.JSON(),
				})
			}
		}
	}

	// Delete
	if err := a.Store.DeleteBySource("file"); err != nil {
		panic(err)
	}

	// Save
	i, j := -50, 0
	for _, p := range plugins {
		p.Source = store.File
		p.ID = util.ID()
		now := time.Now()
		p.CreatedAt = now
		p.UpdatedAt = now
		if _, ok := prePlugins[p.Name]; ok {
			i++
			p.Order = i
		} else {
			j++
			p.Order = j
		}
		if err := a.Store.AddPlugin(p); err != nil {
			panic(err)
		}
	}
}

func (h *Host) FindPath(name string) (p *Path) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	p = h.Paths[name]

	// Add path
	if p == nil {
		p = new(Path)
		h.Paths[name] = p
	}

	// Initialize path
	if !p.initialized {
		p.Name = name
		p.Group = h.Group.Group(name)
		p.initialized = true
	}

	return
}

func (h *Host) AddPlugin(p plugin.Plugin) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.Group.Use(p.Process)
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
