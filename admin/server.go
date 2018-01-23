package admin

import (
	"github.com/labstack/armor"
	"github.com/labstack/echo"
)

func Start(a *armor.Armor) {
	e := echo.New()
	e.HideBanner = true
	// e.Use(middleware.BasicAuth(func(usr, pwd string, _ echo.Context) (bool, error) {
	// 	return usr == "admin" && pwd == "L@B$t@ck0709", nil
	// }))
	h := &handler{armor: a, store: a.Store}

	// Nodes
	nodes := e.Group("/nodes")
	nodes.GET("", h.nodes)

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
	pathPlugins := paths.Group("/plugins")
	pathPlugins.POST("", h.addPlugin)
	pathPlugins.GET("/:plugin", h.findPlugin)
	pathPlugins.PUT("/:plugin", h.updatePlugin)
	pathPlugins.DELETE("/:plugin", h.removePlugin)
	pathPlugins.POST("/targets", h.addProxyTarget)
	pathPlugins.DELETE("/targets/:target", h.removeProxyTarget)

	e.Start(a.Admin.Address)
}
