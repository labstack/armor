package armor

import (
	"log"
	"net"
	"strconv"

	"github.com/hashicorp/logutils"
	"github.com/hashicorp/serf/serf"
)

// Events
const (
	EventPluginLoad   = "1"
	EventPluginUpdate = "2"
)

func (a *Armor) StartCluster() {
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
				case EventPluginLoad, EventPluginUpdate:
					id := string(t.Payload)
					p, err := a.Store.FindPlugin(id)
					if err != nil {
						a.Logger.Error(err)
					}
					a.LoadPlugin(p, t.Name == EventPluginUpdate)
				}
			}
		}
	}
}
