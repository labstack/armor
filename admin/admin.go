package admin

import (
	"github.com/labstack/armor"
	"github.com/labstack/armor/admin/api"
	"github.com/labstack/echo/v4"
)

func loadPlugins(a *armor.Armor) (err error) {
	plugins, err := a.Store.FindPlugins()
	if err != nil {
		return
	}
	for _, p := range plugins {
		a.LoadPlugin(p, false)
	}
	return
}

func Start(a *armor.Armor) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	if !a.DefaultConfig {
		a.Colorer.Printf("⇨ admin server started on %s\n", a.Colorer.Green(a.Admin.Address))
	}
	// e.Use(middleware.BasicAuth(func(usr, pwd string, _ echo.Context) (bool, error) {
	// 	return usr == "admin" && pwd == "L@B$t@ck0709", nil
	// }))

	// Load plugins
	if err := loadPlugins(a); err != nil {
		a.Logger.Fatal(err)
	}

	// API
	if err := api.Init(a, e); err != nil {
		a.Logger.Fatal(err)
	}
}
