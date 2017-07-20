package main

import (
	"flag"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/marpaia/graphite-golang"
	"log"
	"time"
)

var (
	graphiteHost, consulAddr, hostname string
	graphitePort                       int
)

func init() {
	flag.StringVar(&graphiteHost, "ghost", "", "graphite host")
	flag.IntVar(&graphitePort, "gport", 2003, "graphite port")
	flag.StringVar(&consulAddr, "consul", "127.0.0.1:8500", "consul addr")
	flag.StringVar(&hostname, "host", "", "hostname")
	flag.Parse()
}

func main() {
	gh, err := graphite.NewGraphite(graphiteHost, graphitePort)
	gh.Prefix = "consul." + hostname
	if err != nil {
		log.Fatal(err)
	}
	config := &api.Config{
		Address:   consulAddr,
		Scheme:    "http",
		Transport: cleanhttp.DefaultPooledTransport(),
	}
	consulClient, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	for {
		leader, err := consulClient.Status().Leader()
		if err != nil || len(leader) == 0 {
			err = gh.SimpleSend("consul_raft_leader", "0")
			if err != nil {
				log.Println(err)
			}
		} else {
			err = gh.SimpleSend("consul_raft_leader", "1")
			if err != nil {
				log.Println(err)
			}
		}
		time.Sleep(time.Minute)
	}
}
