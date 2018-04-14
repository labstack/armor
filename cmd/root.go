package cmd

import (
	"fmt"
	"io/ioutil"
	stdLog "log"
	"net"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/labstack/armor"
	"github.com/labstack/armor/admin"
	"github.com/labstack/armor/store"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
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
    address: %s 
    plugins:
      - name: logger
      - name: static
        browse: true
        root: %s
  `
)

var (
	configFile string
	port       string
	root       string
	expose     bool
	rootCmd    = &cobra.Command{
		Use:   "armor",
		Short: "Armor is an uncomplicated, modern HTTP server",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is $HOME/.tunnel.yaml)")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "port to listen on")
	rootCmd.PersistentFlags().StringVarP(&root, "root", "", ".", "root directory to serve static content")
	rootCmd.PersistentFlags().BoolVar(&expose, "expose", false, "securely expose server to internet")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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

	// Home dir
	homeDir, err := homedir.Dir()
	if err != nil {
		logger.Fatalf("Failed to find the home directory %v", err)
	}
	a.HomeDir = filepath.Join(homeDir, ".armor")
	if err = os.MkdirAll(a.HomeDir, 0755); err != nil {
		logger.Fatalf("Failed to create config directory %v", err)
	}

	// Config
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		a.DefaultConfig = true
		// Use default config
		data = []byte(fmt.Sprintf(defaultConfig, net.JoinHostPort("", port), root))
	}
	if err = yaml.Unmarshal(data, a); err != nil {
		logger.Fatalf("Failed to parse the config file %v", err)
	}

	// Defaults
	if a.Address == "" {
		a.Address = ":80"
	}
	if a.Storm == nil {
		a.Storm = &armor.Storm{
			URI: filepath.Join(a.HomeDir, "storm.db"),
		}
	}
	if a.Admin == nil {
		a.Admin = &armor.Admin{
			Address: "localhost:8081",
		}
	}
	if a.Cluster == nil {
		a.Cluster = &armor.Cluster{
			Address: ":8082",
			Peers:   []string{"localhost:8082"},
		}
	}
	if a.Hosts == nil {
		a.Hosts = make(armor.Hosts)
	}

	// HTTP
	h := a.NewHTTP()

	// Store
	if a.Postgres != nil {
		a.Store = store.NewPostgres(a.Postgres.URI)
	} else {
		if a.Store, err = store.NewStorm(a.Storm.URI); err != nil {
			logger.Fatalf("Failed to connect to storm %v", err)
		}
	}
	defer a.Store.Close()
	a.SavePlugins()

	// Start cluster
	go a.StartCluster()

	// Start admin
	go admin.Start(a)

	// Create tunnel
	if expose {
		go h.CreateTunnel()
	}

	// Start server
	colorer.Printf(banner, colorer.Red("v"+armor.Version), colorer.Blue(armor.Website))
	if a.TLS != nil {
		go func() {
			logger.Fatal(h.StartTLS())
		}()
	}
	logger.Fatal(h.Start())
}
