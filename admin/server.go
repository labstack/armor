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
	h := &handler{armor: a, store: a.Store, cluster: a.Cluster}

	// Nodes
	nodes := e.Group("/nodes")
	nodes.GET("", h.nodes)

	// Plugins
	plugins := e.Group("/plugins")
	plugins.POST("", h.addPlugin)
	plugins.GET("/:id", h.findPlugin)
	plugins.PUT("/:id", h.updatePlugin)
	plugins.DELETE("/:id", h.removePlugin)
	plugins.POST("/targets", h.addProxyTarget)
	plugins.DELETE("/targets/:target", h.removeProxyTarget)

	e.Start(a.Admin.Address)
}
