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
	hostName := c.Param("host")
	pathName := h.getPath(c)
	rawPlugin := armor.RawPlugin{}

	if err = c.Bind(&rawPlugin); err != nil {
		return
	}

	if hostName == "" && pathName == "" {
		// Global

	} else if hostName != "" && pathName == "" {
		// host := h.armor.Hosts[hostName]
		// if host == nil {
		// 	host = new(armor.Host)
		// 	host.Init(hostName, h.armor)
		// }
	} else if hostName != "" && pathName != "" {
		host := h.armor.FindHost(hostName)
		if host == nil {
			host = h.armor.AddHost(hostName)
		}
		path := host.FindPath(pathName)
		if path == nil {
			path = host.AddPath(pathName)
		}
		p, err := plugin.Decode(rawPlugin, h.armor, host.Echo)
		if err != nil {
			return err
		}
		path.AddPlugin(p)
	}

	return nil
}

func (h *handler) findPlugin(c echo.Context) (err error) {
	return nil
}

func (h *handler) updatePlugin(c echo.Context) (err error) {
	return nil
}

func (h *handler) removePlugin(c echo.Context) (err error) {
	return nil
}

func (h *handler) addProxyTarget(c echo.Context) (err error) {
	host := h.armor.FindHost(c.Param("host"))
	path := host.FindPath(h.getPath(c))
	p := h.lookupPlugin(host, path, c.Param("plugin"))
	proxy := p.(*plugin.Proxy)
	return proxy.AddTarget(c)
}

func (h *handler) removeProxyTarget(c echo.Context) (err error) {
	host := h.armor.FindHost(c.Param("host"))
	path := host.FindPath(h.getPath(c))
	p := h.lookupPlugin(host, path, c.Param("plugin"))
	proxy := p.(*plugin.Proxy)
	return proxy.RemoveTarget(c)
}
