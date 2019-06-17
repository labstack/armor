package api

import (
	"github.com/labstack/armor"
	"github.com/labstack/armor/store"
	"github.com/labstack/echo/v4"
)

type (
	handler struct {
		armor   *armor.Armor
		store   store.Store
		cluster *armor.Cluster
	}
)

func (h *handler) getPath(c echo.Context) string {
	return "/" + c.Param("path")
}

func Init(a *armor.Armor, e *echo.Echo) error {
	h := &handler{armor: a, store: a.Store, cluster: a.Cluster}

	// Nodes
	nodes := e.Group("/nodes")
	nodes.GET("", h.nodes)

	// Plugins
	plugins := e.Group("/plugins")
	plugins.GET("/:id", h.findPlugin)
	plugins.GET("", h.findPlugins)

	// Hosts
	hosts := e.Group("/hosts/:host")
	hostPlugins := hosts.Group("/plugins")
	hostPlugins.POST("", h.addPlugin)
	// hostPlugins.GET("/:id", h.findPlugin)
	// hostPlugins.GET("", h.findPlugins)
	hostPlugins.PUT("/:id", h.savePlugin)

	// Paths
	paths := hosts.Group("/paths/:path")
	pathPlugins := paths.Group("/plugins")
	pathPlugins.POST("", h.addPlugin)
	// pathPlugins.GET("/:id", h.findPlugin)
	// pathPlugins.GET("", h.findPlugins)
	pathPlugins.PUT("/:id", h.savePlugin)

	// plugins.PATCH("/:id/config", h.updatePluginConfig)
	// plugins.DELETE("/:id", h.removePlugin)
	// plugins.POST("/targets", h.addProxyTarget)
	// plugins.DELETE("/targets/:target", h.removeProxyTarget)

	return e.Start(a.Admin.Address)
}
