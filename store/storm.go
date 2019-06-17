package store

import (
	"fmt"

	"github.com/asdine/storm/q"

	"github.com/asdine/storm"
)

type (
	Storm struct {
		*storm.DB
	}
)

func NewStorm(uri string) (s *Storm, err error) {
	s = new(Storm)
	s.DB, err = storm.Open(uri)
	return
}

func (s *Storm) AddPlugin(p *Plugin) error {
	p.Unique = fmt.Sprintf("%s:%s:%s", p.Name, p.Host, p.Path)
	return s.Save(p)
}

func (s *Storm) FindPlugin(id string) (p *Plugin, err error) {
	p = new(Plugin)
	err = s.One("ID", id, p)
	return
}

func (s *Storm) FindPlugins() (plugins []*Plugin, err error) {
	plugins = []*Plugin{}
	if err = s.Select().OrderBy("Order").Find(&plugins); err != nil {
		return
	}
	return plugins, decodeRawPlugin(plugins)
}

func (s *Storm) UpdatePlugin(p *Plugin) error {
	return s.Update(p)
}

func (s *Storm) DeleteBySource(source string) (err error) {
	query := s.Select(q.Eq("Source", source))
	err = query.Delete(new(Plugin))
	if err != nil && err != storm.ErrNotFound {
		return
	}
	return nil
}

func (s *Storm) Close() error {
	return s.DB.Close()
}
