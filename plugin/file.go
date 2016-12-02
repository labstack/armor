package plugin

import "github.com/labstack/echo"

type (
	File struct {
		Base `json:",squash"`
		Path string `json:"path"`
	}
)

func (f *File) Init() (err error) {
	return
}

func (*File) Priority() int {
	return 1
}

func (f *File) Execute(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.File(f.Path)
	}
}
