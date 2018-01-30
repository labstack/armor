package cluster

import (
	"log"
	"net"
	"strconv"

	"github.com/hashicorp/logutils"
	"github.com/hashicorp/serf/serf"
	"github.com/labstack/armor"
	"github.com/labstack/armor/plugin"
)

// Events
const (
	EventPluginAdd    = "1"
	EventPluginUpdate = "2"
)

func Start(a *armor.Armor) {
	conf := serf.DefaultConfig()
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   a.Logger.Output(),
	}
	conf.MemberlistConfig.Logger = log.New(filter, "", log.LstdFlags)
	conf.LogOutput = filter
	conf.NodeName = a.Name
	host, port, err := net.SplitHostPort(a.Cluster.Address)
	if err != nil {
		a.Logger.Fatal(err)
	}
	if host == "" {
		host = "127.0.0.1"
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		a.Logger.Fatal(err)
	}
	conf.MemberlistConfig.BindAddr = host
	conf.MemberlistConfig.BindPort = p
	a.Cluster.Serf, err = serf.Create(conf)
	if err != nil {
		a.Logger.Fatal(err)
	}
	a.Cluster.Join(a.Cluster.Peers, true)
	ch := make(chan serf.Event, 64)
	conf.EventCh = ch
	for {
		select {
		case e := <-ch:
			switch t := e.(type) {
			case serf.UserEvent:
				switch t.Name {
				case EventPluginAdd, EventPluginUpdate:
					id := string(t.Payload)
					p, err := a.Store.FindPlugin(id)
					if err != nil {
						a.Logger.Error(err)
					}

					if p.Host == "" && p.Path == "" {
						// Global level
					} else if p.Host != "" && p.Path == "" {
						// Host level
					} else if p.Host != "" && p.Path != "" {
						// Path level
						host := a.FindHost(p.Host)
						if host == nil {
							host = a.AddHost(p.Host)
						}
						path := host.FindPath(p.Path)
						if path == nil {
							path = host.AddPath(p.Path)
						}
						p, err := plugin.Decode(p.Raw, host.Echo, a.Logger)
						if err != nil {
							a.Logger.Error(err)
						}
						if err = p.Initialize(); err != nil {
							a.Logger.Fatal(err)
						}
						if t.Name == EventPluginAdd {
							path.AddPlugin(p)
						} else {
							path.UpdatePlugin(p)
						}
					}
				}
			}
		}
	}
}
