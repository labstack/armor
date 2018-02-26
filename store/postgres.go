package store

import (
	"github.com/jmoiron/sqlx"
)

type (
	Postgres struct {
		*Base
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
	pg = &Postgres{Base: &Base{db: sqlx.MustConnect("postgres", uri)}}
	pg.db.MustExec(postgresSchema)
	return
}
