package plugin

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

type (
	Static struct {
		Base `json:",squash"`

		// Root directory from where the static content is served.
		// Required.
		Root string `json:"root"`

		// Index file for serving a directory.
		// Optional. Default value "index.html".
		Index string `json:"index"`

		// Enable HTML5 mode by forwarding all not-found requests to root so that
		// SPA (single-page application) can handle the routing.
		// Optional. Default value false.
		HTML5 bool `json:"html5"`

		// Enable directory browsing.
		// Optional. Default value false.
		Browse bool `json:"browse"`

		// Templates
		indexTemplate *Template
	}
)

func (s *Static) Init() (err error) {
	// Defaults
	if s.Root == "" {
		// TODO:
	}
	if s.Index == "" {
		s.Index = "index.html"
	}

	s.indexTemplate = NewTemplate(s.Index)

	return
}

func (*Static) Priority() int {
	return 1
}

func (s *Static) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fs := http.Dir(s.Root)
		name := c.Param("*")
		file, err := fs.Open(name)

		if err != nil {
			if os.IsNotExist(err) {
				if s.HTML5 {
					return s.serveIndex(c, fs)
				}
				return echo.ErrNotFound
			}
			return err
		}

		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			return err
		}
		if fi.IsDir() {
			if s.Browse {
				return s.listDir(file, c.Response())
			}
			return s.serveIndex(c, fs)
		}
		http.ServeContent(c.Response(), c.Request(), fi.Name(), fi.ModTime(), file)
		return nil
	}
}

func (s *Static) serveIndex(c echo.Context, fs http.Dir) error {
	i, err := s.indexTemplate.Execute(c)
	if err != nil {
		return err
	}
	file, err := fs.Open(i)
	if err != nil {
		if os.IsNotExist(err) {
			return echo.ErrNotFound
		}
		return err
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	http.ServeContent(c.Response(), c.Request(), i, fi.ModTime(), file)
	return nil
}

func (s *Static) listDir(dir http.File, res *echo.Response) error {
	dirs, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	// Create a directory index
	res.Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	if _, err = fmt.Fprintf(res, "<pre>\n"); err != nil {
		return err
	}
	for _, d := range dirs {
		name := d.Name()
		color := "#212121"
		if d.IsDir() {
			color = "#e91e63"
			name += "/"
		}
		if _, err = fmt.Fprintf(res, "<a href=\"%s\" style=\"color: %s;\">%s</a>\n", name, color, name); err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(res, "</pre>\n")
	return err
}
