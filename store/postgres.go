package store

import (
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/armor/plugin"
	"github.com/lib/pq"
)

type (
	Postgres struct {
		db *sqlx.DB
	}
)

func (pg *Postgres) AddPlugin(plugin *Plugin) (err error) {
	query := `insert into plugins (id, name, host, path, config, created_at, updated_at)
		values (:id, :name, :host, :path, :config, :created_at, :updated_at)`
	if _, err = pg.db.NamedExec(query, plugin); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				// return api.ErrEmailAlreadyRegistered
			}
		}
	}
	return
}

func (pg *Postgres) FindPlugin(id string) (p *Plugin, err error) {
	query := `select * from plugins where id = $1`
	p = new(Plugin)
	if err = pg.db.Get(p, query, id); err != nil {
		if err == sql.ErrNoRows {
			// return nil, api.ErrEmailNotFound
		}
	}
	p.Raw = plugin.RawPlugin{
		"name": p.Name,
	}
	err = json.Unmarshal(p.Config, &p.Raw)
	return
}

func (pg *Postgres) UpdatePlugin(p *Plugin) (err error) {
	query := `update plugins set name = :name, host = :host, path = :path,
		config = :config, updated_at = :updated_at where id = :id`
	_, err = pg.db.NamedExec(query, p)
	return
}
