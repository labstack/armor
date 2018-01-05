package plugin

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/Knetic/govaluate"
	"github.com/labstack/armor"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/mitchellh/mapstructure"
	"github.com/valyala/fasttemplate"
)

type (
	// Base defines the base struct for plugins.
	Base struct {
		name       string
		Skip       string              `yaml:"skip"`
		Middleware echo.MiddlewareFunc `yaml:"-"`
		Armor      *armor.Armor        `yaml:"-"`
		Echo       *echo.Echo          `yaml:"-"`
		Logger     *log.Logger         `yaml:"-"`
	}

	Template struct {
		*fasttemplate.Template
	}

	Expression struct {
		*fasttemplate.Template
	}
)

var (
	bufferPool sync.Pool
)

// Initialize
func init() {
	bufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}

// lookup returns a plugin by name.
func lookup(base Base) (p armor.Plugin) {
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
	case "rewrite":
		p = &Rewrite{Base: base}
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

// Decode searches the plugin by name, decodes the provided map into plugin and
// calls Plugin#Init().
func Decode(rp armor.RawPlugin, a *armor.Armor, e *echo.Echo) (p armor.Plugin, err error) {
	name := rp["name"].(string)
	base := Base{
		name:   name,
		Skip:   "false",
		Armor:  a,
		Echo:   e,
		Logger: a.Logger,
	}
	if p = lookup(base); p == nil {
		return p, fmt.Errorf("plugin=%s not found", name)
	}
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  p,
	})
	if err = dec.Decode(rp); err != nil {
		return
	}
	return p, p.Init()
}

func (b *Base) Name() string {
	return b.name
}

func NewTemplate(t string) *Template {
	return &Template{Template: fasttemplate.New(t, "${", "}")}
}

func (t *Template) Execute(c echo.Context) (string, error) {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)
	_, err := t.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
		mapTag(buf, c, tag)
		return 0, nil
	})
	return buf.String(), err
}

func NewExpression(t string) *Expression {
	return &Expression{fasttemplate.New(t, "${", "}")}
}

func (e *Expression) Evaluate(c echo.Context) (interface{}, error) {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	if _, err := e.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
		buf.WriteString("'")
		mapTag(buf, c, tag)
		buf.WriteString("'")
		return 0, nil
	}); err != nil {
		return nil, err
	}

	expr, err := govaluate.NewEvaluableExpression(buf.String())
	if err != nil {
		return nil, err
	}
	return expr.Evaluate(nil)
}

func mapTag(b *bytes.Buffer, c echo.Context, t string) {
	switch t {
	case "scheme":
		b.WriteString(c.Scheme())
	case "method":
		b.WriteString(c.Request().Method)
	case "uri":
		b.WriteString(c.Request().RequestURI)
	case "path":
		b.WriteString(c.Request().URL.Path)
	default:
		switch {
		case strings.HasPrefix(t, "header:"):
			b.WriteString(c.Request().Header.Get(t[7:]))
		case strings.HasPrefix(t, "path:"):
			b.WriteString(c.Param(t[5:]))
		case strings.HasPrefix(t, "query:"):
			b.WriteString(c.QueryParam(t[6:]))
		case strings.HasPrefix(t, "form:"):
			b.WriteString(c.FormValue(t[5:]))
		}
	}
}
