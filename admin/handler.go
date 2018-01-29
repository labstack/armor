package admin

import (
	"net/http"
	"time"

	"github.com/labstack/armor"
	"github.com/labstack/armor/cluster"
	"github.com/labstack/armor/plugin"
	"github.com/labstack/armor/store"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/random"
)

type (
	handler struct {
		armor   *armor.Armor
		store   store.Store
		cluster *armor.Cluster
	}

	Node struct {
		Name string `json:"name"`
	}
)

func (h *handler) getPath(c echo.Context) string {
	return "/" + c.Param("path")
}

func (h *handler) nodes(c echo.Context) error {
	cluster := h.armor.Cluster
	nodes := []*Node{}

	for _, m := range cluster.Members() {
		nodes = append(nodes, &Node{
			Name: m.Name,
		})
	}

	return c.JSON(http.StatusOK, nodes)
}

func (h *handler) lookupPlugin(name string, host *armor.Host, path *armor.Path) (p plugin.Plugin) {
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
	if err = c.Bind(&p); err != nil {
		return
	}

	p.ID = random.String(8)
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if err = h.store.AddPlugin(p); err != nil {
		return err
	}

	h.cluster.UserEvent(cluster.EventPluginAdd, []byte(p.ID), true)

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

func (h *handler) updatePlugin(c echo.Context) (err error) {
	p := new(store.Plugin)
	if err = c.Bind(&p); err != nil {
		return
	}
	p.ID = c.Param("id")
	p.UpdatedAt = time.Now()
	err = h.store.UpdatePlugin(p)
	h.cluster.UserEvent(cluster.EventPluginUpdate, []byte(p.ID), true)
	return
}

func (h *handler) removePlugin(c echo.Context) (err error) {
	return nil
}

func (h *handler) addProxyTarget(c echo.Context) (err error) {
	host := h.armor.FindHost(c.Param("host"))
	path := host.FindPath(h.getPath(c))
	p := h.lookupPlugin(c.Param("plugin"), host, path)
	proxy := p.(*plugin.Proxy)
	return proxy.AddTarget(c)
}

func (h *handler) removeProxyTarget(c echo.Context) (err error) {
	host := h.armor.FindHost(c.Param("host"))
	path := host.FindPath(h.getPath(c))
	p := h.lookupPlugin(c.Param("plugin"), host, path)
	proxy := p.(*plugin.Proxy)
	return proxy.RemoveTarget(c)
}
