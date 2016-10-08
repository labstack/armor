package armor

import (
	"crypto/tls"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type (
	Armor struct {
		Address      string                 `json:"address"`
		TLS          *TLS                   `json:"tls"`
		ReadTimeout  time.Duration          `json:"read_timeout"`
		WriteTimeout time.Duration          `json:"write_timeout"`
		Plugins      map[string]interface{} `json:"plugins"`
		Hosts        map[string]*Host       `json:"hosts"`
		Logger       *log.Logger            `json:"-"`
	}

	TLS struct {
		Address      string                      `json:"address"`
		CertFile     string                      `json:"cert_file"`
		KeyFile      string                      `json:"key_file"`
		Auto         bool                        `json:"auto"`
		CacheFile    string                      `json:"cache_file"`
		Certificates map[string]*tls.Certificate `json:"-"`
	}

	Host struct {
		CertFile string                 `json:"cert_file"`
		KeyFile  string                 `json:"key_file"`
		Plugins  map[string]interface{} `json:"plugins"`
		Paths    map[string]*Path       `json:"paths"`
		Echo     *echo.Echo             `json:"-"`
	}

	Path struct {
		Plugins map[string]interface{} `json:"plugins"`
	}
)

const (
	Version = "0.1.3"
)
