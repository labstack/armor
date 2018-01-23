package cluster

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/hashicorp/logutils"
	"github.com/hashicorp/serf/serf"
	"github.com/labstack/armor"
)

type (
	EventType int
)

const (
	EventAddPlugin EventType = iota
	EventUpdatePlugin
)

func Start(a *armor.Armor) {
	conf := serf.DefaultConfig()
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stderr,
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
				fmt.Printf("%s\n", t.Payload)
			}
		}
	}
}
