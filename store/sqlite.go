package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type (
	Sqlite struct {
		*Base
	}
)

const (
	sqliteSchema = `
		create table if not exists plugins (
			id text primary key,
			name text not null,
			host text not null,
			path text not null,
			config blob not null,
			created_at timestamp not null,
			updated_at timestamp not null,
			unique (name, host, path)
		);
	`
)

func NewSqlite(uri string) (s *Sqlite) {
	s = &Sqlite{Base: &Base{db: sqlx.MustConnect("sqlite3", uri)}}
	s.db.MustExec(sqliteSchema)
	return
}
