package store

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
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
	}

	Base struct {
		db *sqlx.DB
	}

	Plugin struct {
		ID        string           `json:"id" db:"id"`
		Name      string           `json:"name" db:"name"`
		Host      string           `json:"host" db:"host"`
		Path      string           `json:"path" db:"path"`
		Config    types.JSONText   `json:"config" db:"config"`
		Source    string           `json:"source,omitempty" db:"source"`
		CreatedAt time.Time        `json:"created_at" db:"created_at"`
		UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
		Raw       plugin.RawPlugin `json:"-" db:"-"`
	}
)

const (
	API       = "api"
	Consul    = "consul"
	ETCD      = "etcd"
	File      = "file"
	Zookeeper = "zookeeper"
)

func (b *Base) AddPlugin(plugin *Plugin) (err error) {
	query := `insert into plugins (id, name, host, path, config, source, created_at,
		updated_at) values (:id, :name, :host, :path, :config, :source, :created_at,
		:updated_at)`
	_, err = b.db.NamedExec(query, plugin)
	return
}

func (s *Base) FindPlugin(id string) (p *Plugin, err error) {
	query := `select * from plugins where id = $1`
	p = new(Plugin)
	if err = s.db.Get(p, query, id); err != nil {
		if err == sql.ErrNoRows {
		}
	}
	p.Raw = plugin.RawPlugin{
		"name": p.Name,
	}
	err = json.Unmarshal(p.Config, &p.Raw)
	return
}

func (b *Base) FindPlugins() (plugins []*Plugin, err error) {
	query := `select * from plugins`
	plugins = []*Plugin{}

	if err = b.db.Select(&plugins, query); err != nil {
		if err == sql.ErrNoRows {
			// return nil, api.ErrEmailNotFound
		}
	}
	for _, p := range plugins {
		p.Raw = plugin.RawPlugin{
			"name": p.Name,
		}
		if err = json.Unmarshal(p.Config, &p.Raw); err != nil {
			return
		}
	}

	return
}

func (b *Base) UpdatePlugin(p *Plugin) (err error) {
	query := `update plugins set config = :config, updated_at = :updated_at
		where name = :id and host = :host and path = :path`
	_, err = b.db.NamedExec(query, p)
	return
}
