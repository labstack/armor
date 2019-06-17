package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/armor"
	"github.com/labstack/armor/plugin"
	"github.com/labstack/armor/store"
	"github.com/labstack/armor/util"
	"github.com/labstack/echo/v4"
)

func decodePath(c echo.Context) string {
	return strings.Replace(c.Param("path"), "~", "/", 1)
}

func lookupPlugin(name string, host *armor.Host, path *armor.Path) (p plugin.Plugin) {
	plugins := []plugin.Plugin{}

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
		if p.Name() == name {
			return p
		}
	}
	return nil
}

func (h *handler) addPlugin(c echo.Context) (err error) {
	p := new(store.Plugin)
	if err = c.Bind(p); err != nil {
		return
	}
	p.ID = util.ID()
	p.Host = c.Param("host")
	p.Path = decodePath(c)
	p.Source = store.API
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if err = h.store.AddPlugin(p); err != nil {
		return err
	}
	h.cluster.UserEvent(armor.EventPluginLoad, []byte(p.ID), true)
	return c.JSON(http.StatusCreated, p)

	// hostName := c.Param("host")
	// pathName := h.getPath(c)
	// rawPlugin := armor.RawPlugin{}

	// if err = c.Bind(&rawPlugin); err != nil {
	// 	return
	// }
	// name := rawPlugin["name"].(string)

	// if hostName == "" && pathName == "" {
	// Global

	// } else if hostName != "" && pathName == "" {
	// host := h.armor.Hosts[hostName]
	// if host == nil {
	// 	host = new(armor.Host)
	// 	host.Init(hostName, h.armor)
	// }
	// } else if hostName != "" && pathName != "" {
	// host := h.armor.FindHost(hostName)
	// if host == nil {
	// 	host = h.armor.AddHost(hostName)
	// }
	// path := host.FindPath(pathName)
	// if path == nil {
	// 	path = host.AddPath(pathName)
	// }
	// p, err := plugin.Decode(rawPlugin, h.armor, host.Echo)
	// if err != nil {
	// 	return err
	// }
	// h.cluster.UserEvent(cluster.EventPluginAdd, []byte(name), true)
	// if err = h.store.AddPlugin(name, hostName, pathName, rawPlugin); err != nil {
	// 	return err
	// }
	// path.AddPlugin(p)
	// }

}

func (h *handler) findPlugin(c echo.Context) (err error) {
	id := c.Param("id")
	p, err := h.store.FindPlugin(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, p)
}

func (h *handler) findPlugins(c echo.Context) (err error) {
	plugins, err := h.store.FindPlugins()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, plugins)
}

// TODO: Support saving by name
// update plugins set config = jsonb_set(config, '{root}', '"/tmp/roots"') where id = 'l7J3';
func (h *handler) savePlugin(c echo.Context) (err error) {
	p := new(store.Plugin)
	if err = c.Bind(p); err != nil {
		return
	}
	p.Host = c.Param("host")
	p.Path = decodePath(c)
	p.Source = store.API
	p.UpdatedAt = time.Now()
	err = h.store.UpdatePlugin(p)
	if err != nil {
		return
	}
	h.cluster.UserEvent(armor.EventPluginUpdate, []byte(p.ID), true)
	return c.NoContent(http.StatusNoContent)
}

// TODO: Support saving by name
// func (h *handler) updatePluginConfig(c echo.Context) (err error) {
// 	id := c.Param("id")
// 	config := echo.Map{}
// 	if err = c.Bind(&config); err != nil {
// 		return
// 	}

// 	// Find from DB
// 	pi, err := h.store.FindPlugin(id)
// 	if err != nil {
// 		return
// 	}

// 	for k, v := range config {
// 		jsonb_set(data, '{name}', '"my-other-name"')
// 	}
// 	fmt.Printf("%+v, %s", config, id)
// 	// p1, _ := plugin.Decode(update.Raw, nil, nil)

// 	// // Find from DB
// 	// pluginDB, err := h.store.FindPlugin(update.ID)
// 	// if err != nil {
// 	// 	return
// 	// }
// 	// p2, _ := plugin.Decode(update.Raw, nil, nil)

// 	// host := h.armor.FindHost(pluginDB.Host)
// 	// path := host.FindPath(pluginDB.Path)
// 	// p := h.lookupPlugin(pluginDB.Name, host, path)

// 	// switch t := p.(type) {
// 	// case *plugin.Proxy:
// 	// 	if update
// 	// 	t.Targets =
// 	// }

// 	// p.ID = c.Param("id")
// 	// host := h.armor.FindHost(p.Host)
// 	// path := host.FindPath(p.Path)
// 	// p =
// 	// p := h.lookupPlugin(name, host, path)
// 	// proxy := p.(*plugin.Proxy)
// 	// return proxy.AddTarget(c)
// 	return
// }

// func (h *handler) removePlugin(c echo.Context) (err error) {
// 	host := h.armor.FindHost(c.Param("host"))
// 	path := host.FindPath(h.getPath(c))
// 	p := h.lookupPlugin(c.Param("plugin"), host, path)
// 	proxy := p.(*plugin.Proxy)
// 	return proxy.AddTarget(c)
// }

// func (h *handler) addProxyTarget(c echo.Context) (err error) {
// 	host := h.armor.FindHost(c.Param("host"))
// 	path := host.FindPath(h.getPath(c))
// 	p := h.lookupPlugin(c.Param("plugin"), host, path)
// 	proxy := p.(*plugin.Proxy)
// 	return proxy.AddTarget(c)
// }

// func (h *handler) removeProxyTarget(c echo.Context) (err error) {
// 	host := h.armor.FindHost(c.Param("host"))
// 	path := host.FindPath(h.getPath(c))
// 	p := h.lookupPlugin(c.Param("plugin"), host, path)
// 	proxy := p.(*plugin.Proxy)
// 	return proxy.RemoveTarget(c)
// }
