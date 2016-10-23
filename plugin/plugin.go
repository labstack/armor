package plugin

import (
	"fmt"

	"github.com/labstack/armor"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/mitchellh/mapstructure"
)

type (
	Plugin interface {
		Name() string
		Initialize() error
		Process(echo.HandlerFunc) echo.HandlerFunc
		Priority() int
		Terminate()
	}

	// Base defines the base struct for plugins.
	Base struct {
		name       string
		Middleware echo.MiddlewareFunc `json:"-"`
		Host       string              `json:"-"`
		Path       string              `json:"-"`
		Armor      *armor.Armor        `json:"-"`
		Logger     *log.Logger         `json:"-"`
	}
)

func (b *Base) Name() string {
	return b.name
}

// Decode searches the plugin by name, decodes the provided map into plugin and
// calls Plugin#Initialize().
func Decode(name string, i interface{}, host string, path string, a *armor.Armor) (p Plugin, err error) {
	base := Base{
		name:   name,
		Host:   host,
		Path:   path,
		Armor:  a,
		Logger: a.Logger,
	}
	if p = Lookup(base); p == nil {
		return p, fmt.Errorf("plugin=%s not found", name)
	}
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  p,
	})
	if err = dec.Decode(i); err != nil {
		return
	}
	return p, p.Initialize()
}

// Lookup returns a plugin by name.
func Lookup(base Base) (p Plugin) {
	switch base.Name() {
	case "body-limit":
		p = &BodyLimit{Base: base}
	case "logger":
		p = &Logger{Base: base}
	case "redirect":
		p = &Redirect{Base: base}
	case "https-redirect":
		p = &HTTPSRedirect{Base: base}
	case "https-www-redirect":
		p = &HTTPSWWWRedirect{Base: base}
	case "https-non-www-redirect":
		p = &HTTPSNonWWWRedirect{Base: base}
	case "www-redirect":
		p = &WWWRedirect{Base: base}
	case "non-www-redirect":
		p = &NonWWWRedirect{Base: base}
	case "add-trailing-slash":
		p = &AddTrailingSlash{Base: base}
	case "remove-trailing-slash":
		p = &RemoveTrailingSlash{Base: base}
	case "cors":
		p = &CORS{Base: base}
	case "gzip":
		p = &Gzip{Base: base}
	case "header":
		p = &Header{Base: base}
	case "proxy":
		p = &Proxy{Base: base}
	case "static":
		p = &Static{Base: base}
	case "nats":
		// p = &NATS{Base: base}
	}
	return
}
