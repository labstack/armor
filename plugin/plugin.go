package plugin

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
	"sync"

	"github.com/labstack/armor"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/mitchellh/mapstructure"
	"github.com/valyala/fasttemplate"
)

type (
	Plugin interface {
		Name() string
		Init() error
		Process(echo.HandlerFunc) echo.HandlerFunc
		Priority() int
	}

	// Base defines the base struct for plugins.
	Base struct {
		name       string
		Middleware echo.MiddlewareFunc `json:"-"`
		Armor      *armor.Armor        `json:"-"`
		Logger     *log.Logger         `json:"-"`
	}

	Template struct {
		template *fasttemplate.Template
	}
)

var (
	bufferPool sync.Pool
)

// Decode searches the plugin by name, decodes the provided map into plugin and
// calls Plugin#Init().
func Decode(name string, i interface{}, a *armor.Armor) (p Plugin, err error) {
	base := Base{
		name:   name,
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
	return p, p.Init()
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
	case "secure":
		p = &Secure{Base: base}
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
	case "file":
		p = &File{Base: base}
	case "nats":
		// p = &NATS{Base: base}
	}
	return
}

func (b *Base) Name() string {
	return b.name
}

func NewTemplate(src string) *Template {
	return &Template{
		template: fasttemplate.New(src, "${", "}"),
	}
}

func (t *Template) Execute(c echo.Context) (string, error) {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	_, err := t.template.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
		switch tag {
		case "scheme":
			return buf.Write([]byte(c.Scheme()))
		case "method":
			return buf.Write([]byte(c.Request().Method))
		case "path":
			return buf.Write([]byte(c.Request().URL.Path))
		case "uri":
			return buf.Write([]byte(c.Request().RequestURI))
		case "dir":
			return buf.Write([]byte(path.Dir(c.Param("*"))))
		default:
			switch {
			case strings.HasPrefix(tag, "header:"):
				return buf.Write([]byte(c.Request().Header.Get(tag[7:])))
			case strings.HasPrefix(tag, "path:"):
				return buf.Write([]byte(c.Param(tag[5:])))
			case strings.HasPrefix(tag, "query:"):
				return buf.Write([]byte(c.QueryParam(tag[6:])))
			case strings.HasPrefix(tag, "form:"):
				return buf.Write([]byte(c.FormValue(tag[5:])))
			default:
				// TODO:
			}
		}

		return 0, nil
	})

	return buf.String(), err
}

func init() {
	bufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}
