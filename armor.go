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
		Plugins      []Plugin         `yaml:"plugins"`
		Hosts        map[string]*Host `yaml:"hosts"`
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
		Name     string           `yaml:"-"`
		CertFile string           `yaml:"cert_file"`
		KeyFile  string           `yaml:"key_file"`
		Plugins  []Plugin         `yaml:"plugins"`
		Paths    map[string]*Path `yaml:"paths"`
		Echo     *echo.Echo       `yaml:"-"`
	}

	Path struct {
		Plugins []Plugin `yaml:"plugins"`
	}

	Plugin map[string]interface{}
)

const (
	Version = "0.2.12"
	Website = "https://armor.labstack.com"
)
