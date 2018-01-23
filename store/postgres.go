package store

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/armor"
)

type (
	Postgres struct {
		db *sqlx.DB
	}
)

func (p *Postgres) AddPlugin(plugin armor.RawPlugin) error {
	name := plugin["name"].(string)
	now := time.Now()
	config, err := json.Marshal(plugin)
	if err != nil {
		return err
	}
	query := `insert into plugins (name, config, enabled, created_at, updated_at)
		values ($1, $2, $3, $4, $5)`
	_, err = p.db.Exec(query, name, config, true, now, now)
	return err
}
