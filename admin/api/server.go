package api

import (
	"github.com/labstack/armor"
	"github.com/labstack/echo"
)

func Start(a *armor.Armor) {
	e := echo.New()
	h := &handler{armor: a}

	// Global

	// Hosts
	hosts := e.Group("/hosts/:host")

	// Host plugins
	hostPlugins := hosts.Group("/plugins/:plugin")
	hostPlugins.POST("", h.addPlugin)
	hostPlugins.GET("/:plugin", h.findPlugin)
	hostPlugins.PUT("/:plugin", h.updatePlugin)
	hostPlugins.DELETE("/:plugin", h.removePlugin)
	hostPlugins.PATCH("/targets", h.addProxyTarget)
	hostPlugins.DELETE("/targets/:target", h.removeProxyTarget)

	// Paths
	paths := hosts.Group("/paths/:path")

	// Path plugins
	pathPlugins := paths.Group("/plugins/:plugin")
	pathPlugins.POST("", h.addPlugin)
	pathPlugins.GET("/:plugin", h.findPlugin)
	pathPlugins.PUT("/:plugin", h.updatePlugin)
	pathPlugins.DELETE("/:plugin", h.removePlugin)
	pathPlugins.POST("/targets", h.addProxyTarget)
	pathPlugins.DELETE("/targets/:target", h.removeProxyTarget)

	e.Start(":8081")
}
