package api

import (
	"github.com/labstack/armor"
	"github.com/labstack/armor/plugin"
	"github.com/labstack/echo"
)

type (
	handler struct {
		armor *armor.Armor
	}
)

func (h *handler) getPath(c echo.Context) string {
	return "/" + c.Param("path")
}

func (h *handler) lookupPlugin(host *armor.Host, path *armor.Path, plugin string) (p armor.Plugin) {
	plugins := []armor.Plugin{}

	// Global
	if host == nil && path == nil {
	} else if host != nil && path == nil {
		// Host
		plugins = host.Plugins
	} else if host != nil && path != nil {
		// Path
		plugins = path.Plugins
	}

	for _, p := range plugins {
		if p.Name() == plugin {
			return p
		}
	}
	return nil
}

func (h *handler) addPlugin(c echo.Context) (err error) {
	// host := c.Param("host")
	// path := c.Param("path")

	return nil
}

func (h *handler) findPlugin(c echo.Context) (err error) {
	// host := c.Param("host")
	// path := c.Param("path")

	return nil
}

func (h *handler) updatePlugin(c echo.Context) (err error) {
	// host := c.Param("host")
	// path := c.Param("path")

	return nil
}

func (h *handler) removePlugin(c echo.Context) (err error) {
	// host := c.Param("host")
	// path := c.Param("path")

	return nil
}

func (h *handler) addProxyTarget(c echo.Context) (err error) {
	host := h.armor.Hosts[c.Param("host")]
	path := host.Paths[h.getPath(c)]
	p := h.lookupPlugin(host, path, c.Param("plugin"))
	proxy := p.(*plugin.Proxy)
	return proxy.AddTarget(c)
}

func (h *handler) removeProxyTarget(c echo.Context) (err error) {
	host := h.armor.Hosts[c.Param("host")]
	path := host.Paths[h.getPath(c)]
	p := h.lookupPlugin(host, path, c.Param("plugin"))
	proxy := p.(*plugin.Proxy)
	return proxy.RemoveTarget(c)
}
