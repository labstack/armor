package plugin

import "github.com/labstack/echo"

type (
	File struct {
		Base         `json:",squash"`
		Path         string `json:"path"`
		pathTemplate *Template
	}
)

func (f *File) Init() (err error) {
	f.pathTemplate = NewTemplate(f.Path)
	return
}

func (*File) Priority() int {
	return 1
}

func (f *File) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		p, err := f.pathTemplate.Execute(c)
		if err != nil {
			return err
		}
		return c.File(p)
	}
}
