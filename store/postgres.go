package store

import (
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/armor/plugin"
)

type (
	Postgres struct {
		*sqlx.DB
	}
)

const (
	postgresSchema = `
		create table if not exists plugins (
			id text primary key,
			name text not null,
			host text not null,
			path text not null,
			config jsonb not null,
			source text not null,
			created_at timestamptz not null,
			updated_at timestamptz not null,
			unique (name, host, path)
		);
	`
)

func NewPostgres(uri string) (pg *Postgres) {
	pg = new(Postgres)
	pg.DB = sqlx.MustConnect("postgres", uri)
	pg.MustExec(postgresSchema)
	return
}

func (pg *Postgres) AddPlugin(p *Plugin) (err error) {
	query := `insert into plugins (id, name, host, path, config, source, created_at,
		updated_at) values (:id, :name, :host, :path, :config, :source, :created_at,
		:updated_at)`
	_, err = pg.NamedExec(query, p)
	return
}

func (pg *Postgres) FindPlugin(id string) (p *Plugin, err error) {
	query := `select * from plugins where id = $1`
	p = new(Plugin)
	if err = pg.Get(p, query, id); err != nil {
		if err == sql.ErrNoRows {
		}
	}
	p.Raw = plugin.RawPlugin{
		"name":  p.Name,
		"order": p.Order,
	}
	err = json.Unmarshal(p.Config, &p.Raw)
	return
}

func (pg *Postgres) FindPlugins() (plugins []*Plugin, err error) {
	query := `select * from plugins`
	plugins = []*Plugin{}

	if err = pg.Select(&plugins, query); err != nil {
		if err == sql.ErrNoRows {
			// return nil, api.ErrEmailNotFound
		}
	}

	return plugins, decodeRawPlugin(plugins)
}

func (pg *Postgres) UpdatePlugin(p *Plugin) (err error) {
	query := `update plugins set config = :config, updated_at = :updated_at
		where name = :id and host = :host and path = :path`
	_, err = pg.NamedExec(query, p)
	return
}

func (pg *Postgres) DeleteBySource(source string) (err error) {
	query := `delete from plugins where source = $1`
	_, err = pg.Exec(query, source)
	return
}

func (pg *Postgres) Close() error {
	return pg.DB.Close()
}
