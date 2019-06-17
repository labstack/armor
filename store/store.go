package store

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/labstack/armor/plugin"
	_ "github.com/lib/pq"
)

type (
	Store interface {
		AddPlugin(*Plugin) error
		FindPlugin(string) (*Plugin, error)
		FindPlugins() ([]*Plugin, error)
		UpdatePlugin(*Plugin) error
		DeleteBySource(source string) error
		Close() error
	}

	Plugin struct {
		ID        string           `json:"id" db:"id" storm:"id"`
		Name      string           `json:"name" db:"name"`
		Order     int              `json:"order" db:"order"`
		Host      string           `json:"host" db:"host"`
		Path      string           `json:"path" db:"path"`
		Config    types.JSONText   `json:"config" db:"config"`
		Source    string           `json:"source,omitempty" db:"source"`
		CreatedAt time.Time        `json:"created_at" db:"created_at"`
		UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
		Raw       plugin.RawPlugin `json:"-" db:"-"`
		Unique    string           `json:"-" db:"-" storm:"unique"`
	}
)

const (
	API       = "api"
	Consul    = "consul"
	ETCD      = "etcd"
	File      = "file"
	Zookeeper = "zookeeper"
)

func decodeRawPlugin(plugins []*Plugin) (err error) {
	for _, p := range plugins {
		p.Raw = plugin.RawPlugin{
			"name":  p.Name,
			"order": p.Order,
		}
		if err = json.Unmarshal(p.Config, &p.Raw); err != nil {
			return
		}
	}
	return
}
