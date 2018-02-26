package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	stdLog "log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/lib/pq"
	sqlite3 "github.com/mattn/go-sqlite3"

	"github.com/labstack/armor"
	"github.com/labstack/armor/admin"
	"github.com/labstack/armor/cluster"
	"github.com/labstack/armor/http"
	"github.com/labstack/armor/store"
	"github.com/labstack/armor/util"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
)

const (

	// http://patorjk.com/software/taag/#p=display&f=Small%20Slant&t=Armor
	banner = `
   ___                     
  / _ | ______ _  ___  ____
 / __ |/ __/  ' \/ _ \/ __/
/_/ |_/_/ /_/_/_/\___/_/    %s

Uncomplicated, modern HTTP server
%s
________________________O/_______
                        O\
`
	defaultConfig = `
    name: armor
    address: :8080
    admin:
      address: 127.0.0.1:8081
    cluster:
      address: :8082
      peers:
        - 127.0.0.1:8082
    sqlite:
      uri: %s
    plugins:
      -
        name: logger
      -
        name: static
        browse: true
        root: .
  `
)

func savePlugins(a *armor.Armor) {
	plugins := []*store.Plugin{}

	// Global plugins
	for _, rp := range a.RawPlugins {
		plugins = append(plugins, &store.Plugin{
			Name:   rp.Name(),
			Config: rp.Bytes(),
		})
	}

	for hn, host := range a.Hosts {
		// Host plugins
		for _, rp := range host.RawPlugins {
			plugins = append(plugins, &store.Plugin{
				Name:   rp.Name(),
				Host:   hn,
				Config: rp.Bytes(),
			})
		}

		for pn, path := range host.Paths {
			// Path plugins
			for _, rp := range path.RawPlugins {
				plugins = append(plugins, &store.Plugin{
					Name:   rp.Name(),
					Host:   hn,
					Path:   pn,
					Config: rp.Bytes(),
				})
			}
		}
	}

	// Save
	for _, p := range plugins {
		now := time.Now()
		p.ID = util.ID()
		p.CreatedAt = now
		p.UpdatedAt = now

		if err := a.Store.AddPlugin(p); err != nil {
			switch e := err.(type) {
			case sqlite3.Error:
				if e.Code != sqlite3.ErrConstraint {
					panic(err)
				}
			case *pq.Error:
				if e.Code != "23505" {
					panic(err)
				}
			}
		}
	}
}

func main() {
	// Initialize
	logger := log.New("armor")
	colorer := color.New()
	logger.SetLevel(log.INFO)
	colorer.SetOutput(logger.Output())
	stdLog.SetOutput(logger.Output())
	a := &armor.Armor{
		Logger:  logger,
		Colorer: colorer,
	}

	// Global flags
	config := flag.String("c", "", "config file")
	port := flag.String("p", "", "listen port")
	version := flag.Bool("v", false, "armor version")

	// daemon := flag.Bool("d", false, "run in daemon mode")
	// -daemon
	// -p [http port]
	// -w [--www]
	// -v [--version]
	// -h [--help]
	// --pid
	// Commands
	// - stop
	// - restart
	// - reload
	// port := flag.String("p", "", "the port to bind to")
	// directory?
	flag.Parse()
	// if *config == nil {
	// 	config = "config.yaml"
	// }
	if *version {
		color.Printf("armor %s\n", color.Red("v"+armor.Version))
		os.Exit(0)
	}

	// Config - start
	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}
	// Load
	data, err := ioutil.ReadFile(*config)
	if err != nil {
		// Use default config
		data = []byte(fmt.Sprintf(defaultConfig, filepath.Join(wd, "armor.db")))
	}
	if err = yaml.Unmarshal(data, a); err != nil {
		logger.Fatalf("armor: not able to parse the config file, error=%v", err)
	}

	// Flags should override
	if *port != "" {
		a.Address = net.JoinHostPort("", *port)
	}

	// Defaults
	if a.Address == "" {
		a.Address = ":80"
	}
	if a.Hosts == nil {
		a.Hosts = make(armor.Hosts)
	}
	// Config - end

	// Init http
	h := http.Init(a)

	// Store
	if a.SQLite != nil {
		a.Store = store.NewSqlite(a.SQLite.URI)
	}
	if a.Postgres != nil {
		a.Store = store.NewPostgres(a.Postgres.URI)
	}
	savePlugins(a)

	// Start cluster
	if a.Cluster != nil {
		go cluster.Start(a)
	}

	// Start admin
	if a.Admin != nil {
		go admin.Start(a)
	}

	// Start server - start
	colorer.Printf(banner, colorer.Red("v"+armor.Version), colorer.Blue(armor.Website))
	if a.TLS != nil {
		go func() {
			logger.Fatal(h.StartTLS())
		}()
	}
	logger.Fatal(h.Start())
	// Start server - end
}
