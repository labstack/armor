package plugin

import "github.com/labstack/echo/v4"

type (
	FileConfig struct {
		Path string `yaml:"path"`
	}

	File struct {
		Base       `yaml:",squash"`
		FileConfig `yaml:",squash"`
	}
)

func (f *File) Initialize() {
}

func (f *File) Update(p Plugin) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.FileConfig = p.(*File).FileConfig
	f.Initialize()
}

func (f *File) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		f.mutex.RLock()
		defer f.mutex.RUnlock()
		return c.File(f.Path)
	}
}
