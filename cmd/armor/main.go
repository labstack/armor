package main

import (
	"flag"
	"io/ioutil"
	stdLog "log"
	"net"
	"os"

	"github.com/go-yaml/yaml"

	"github.com/labstack/armor"
	"github.com/labstack/armor/http"
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
    address: ":8080"
    plugins:
    - name: logger
    - name: static
      browse: true
      root: "."
  `
)

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
	if *version {
		color.Printf("armor %s\n", color.Red("v"+armor.Version))
		os.Exit(0)
	}

	// Load config
	data, err := ioutil.ReadFile(*config)
	if err != nil {
		// Use default config
		data = []byte(defaultConfig)
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

	// Initialize and load the plugins
	h := http.Init(a)
	h.LoadPlugins()

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
