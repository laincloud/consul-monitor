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
	graphiteHost, consulAddr string
	graphitePort             int
)

func init() {
	flag.StringVar(&graphiteHost, "host", "", "graphite host")
	flag.IntVar(&graphitePort, "port", 2003, "graphite host")
	flag.StringVar(&consulAddr, "consul", "127.0.0.1:8500", "consul addr")
	flag.Parse()
}

func main() {
	gh, err := graphite.NewGraphite(graphiteHost, graphitePort)
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
			gh.SimpleSend("consul_raft_leader", "0")
		} else {
			gh.SimpleSend("consul_raft_leader", "1")
		}
		time.Sleep(time.Second)
	}
}
