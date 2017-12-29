package plugin

import "github.com/labstack/echo"

type (
	File struct {
		Base `yaml:",squash"`
		Path string `yaml:"path"`
	}
)

func (f *File) Init() (err error) {
	return
}

func (*File) Priority() int {
	return 1
}

func (f *File) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.File(f.Path)
	}
}
