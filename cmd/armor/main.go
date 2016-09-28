package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/labstack/armor"
	"github.com/labstack/armor/http"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
)

const (
	version = "0.1.0"
	banner  = `
 _______  ______    __   __  _______  ______
|   _   ||    _ |  |  |_|  ||       ||    _ |
|  |_|  ||   | ||  |       ||   _   ||   | ||
|       ||   |_||_ |       ||  | |  ||   |_||_
|       ||    __  ||       ||  |_|  ||    __  |
|   _   ||   |  | || ||_|| ||       ||   |  | |
|__| |__||___|  |_||_|   |_||_______||___|  |_|

                                      %s

              Simple HTTP Server
      https://github.com/labstack/armor
_________________ O/___________________________
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
	config := flag.String("c", "", "armor config file")
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
	// version := flag.String("v", "", "print version")
	// directory?
	flag.Parse()

	// Commands
	// TODO:

	// Load config
	data, err := ioutil.ReadFile(*config)
	if err != nil {
		// Use default config
		data = []byte(defaultConfig)
	}
	if err = json.Unmarshal(data, a); err != nil {
		logger.Fatal(err)
	}

	// Flags should override
	// if *port != "" {
	//   a.HTTP.Address = net.JoinHostPort("", *port)
	// }
	color.Printf(banner+"\n", color.Red("v"+version))
	http.Start(a)
}
