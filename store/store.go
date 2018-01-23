package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/labstack/armor"
)

func New(a *armor.Armor) (armor.Store, error) {
	if a.Postgres != nil {
		db, err := sqlx.Connect("postgres", a.Postgres.URI)
		if err != nil {
			return nil, err
		}
		return &Postgres{
			db: db,
		}, nil
	}
	return nil, nil
}
