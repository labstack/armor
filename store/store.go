package store

import (
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
		UpdatePlugin(*Plugin) error
	}

	Plugin struct {
		ID        string           `json:"id" db:"id"`
		Name      string           `json:"name" db:"name"`
		Host      string           `json:"host" db:"host"`
		Path      string           `json:"path" db:"path"`
		Config    types.JSONText   `json:"config" db:"config"`
		CreatedAt time.Time        `json:"created_at" db:"created_at"`
		UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
		Raw       plugin.RawPlugin `json:"-" db:"-"`
	}
)

func New(uri string) (Store, error) {
	db, err := sqlx.Connect("postgres", uri)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		db: db,
	}, nil
}
