package armor

import (
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
)

type (
	Armor struct {
		Address      string           `yaml:"address"`
		TLS          *TLS             `yaml:"tls"`
		ReadTimeout  time.Duration    `yaml:"read_timeout"`
		WriteTimeout time.Duration    `yaml:"write_timeout"`
		RawPlugins   []RawPlugin      `yaml:"plugins"`
		Hosts        map[string]*Host `yaml:"hosts"`
		Plugins      []Plugin         `yaml:"-"`
		Logger       *log.Logger      `yaml:"-"`
		Colorer      *color.Color     `yaml:"-"`
	}

	TLS struct {
		Address  string `yaml:"address"`
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
		Auto     bool   `yaml:"auto"`
		CacheDir string `yaml:"cache_dir"`
	}

	Host struct {
		Name       string           `yaml:"-"`
		CertFile   string           `yaml:"cert_file"`
		KeyFile    string           `yaml:"key_file"`
		RawPlugins []RawPlugin      `yaml:"plugins"`
		Paths      map[string]*Path `yaml:"paths"`
		Plugins    []Plugin         `yaml:"-"`
		Echo       *echo.Echo       `yaml:"-"`
	}

	Path struct {
		RawPlugins []RawPlugin `yaml:"plugins"`
		Plugins    []Plugin    `yaml:"-"`
	}

	RawPlugin map[string]interface{}

	Plugin interface {
		Name() string
		Init() error
		Process(echo.HandlerFunc) echo.HandlerFunc
		Priority() int
	}
)

const (
	Version = "0.3.4"
	Website = "https://armor.labstack.com"
)
