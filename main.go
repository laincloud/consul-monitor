package main

import (
	"flag"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/marpaia/graphite-golang"
	"log"
	"os"
	"time"
)

var (
	graphiteHost, consulAddr string
	graphitePort             int
)

func init() {
	flag.StringVar(&graphiteHost, "host", "", "graphite host")
	flag.IntVar(&graphitePort, "port", 2003, "graphite port")
	flag.StringVar(&consulAddr, "consul", "127.0.0.1:8500", "consul addr")
	flag.Parse()
}

func main() {
	hostname, err := os.Hostname()
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
		gh, err := graphite.NewGraphite(graphiteHost, graphitePort)
		if err != nil {
			log.Println(err)
			continue
		}
		gh.Prefix = "consul." + hostname
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
