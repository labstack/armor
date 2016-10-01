package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"os"

	"github.com/labstack/armor"
	"github.com/labstack/armor/http"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
)

const (
	banner = `
 _______  ______    __   __  _______  ______
|   _   ||    _ |  |  |_|  ||       ||    _ |
|  |_|  ||   | ||  |       ||   _   ||   | ||
|       ||   |_||_ |       ||  | |  ||   |_||_
|       ||    __  ||       ||  |_|  ||    __  |
|   _   ||   |  | || ||_|| ||       ||   |  | |
|__| |__||___|  |_||_|   |_||_______||___|  |_|

                                      %s

Simple HTTP server, supports HTTP/2 and auto TLS
      %s
___________________O/___________________________
                   O\
`
	defaultConfig = `{
    "address": ":8080",
    "plugins": {
      "logger": {},
      "static": {
        "browse": true,
        "root": "."
      }
    }
  }`
)

func main() {
	// Initialize
	logger := log.New("armor")
	a := &armor.Armor{
		Logger: logger,
	}

	// Global flags
	c := flag.String("c", "", "config file")
	p := flag.String("p", "", "listen port")
	v := flag.Bool("v", false, "print the version")

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

	if *v {
		color.Printf("armor %s", color.Red("v"+armor.Version))
		os.Exit(0)
	}

	// Load config
	data, err := ioutil.ReadFile(*c)
	if err != nil {
		// Use default config
		data = []byte(defaultConfig)
	}
	if err = json.Unmarshal(data, a); err != nil {
		logger.Fatal(err)
	}

	// Flags should override
	if *p != "" {
		a.Address = net.JoinHostPort("", *p)
	}
	color.Printf(banner+"\n", color.Red("v"+armor.Version), color.Blue("https://github.com/labstack/armor"))
	http.Start(a)
}
